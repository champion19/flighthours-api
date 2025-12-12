package interactor

import (
	"context"

	"github.com/champion19/flighthours-api/core/interactor/dto"
	"github.com/champion19/flighthours-api/core/interactor/services/domain"
	"github.com/champion19/flighthours-api/core/ports/input"
	"github.com/champion19/flighthours-api/platform/logger"
)

type Interactor struct {
	service input.Service
	logger  logger.Logger
}

func NewInteractor(service input.Service, log logger.Logger) *Interactor {
	return &Interactor{
		service: service,
		logger:  log,
	}
}
func (i *Interactor) RegisterEmployee(ctx context.Context, employee domain.Employee) (result *dto.RegisterEmployee, err error) {
	i.logger.Info(logger.LogEmployeeInteractorRegStart, employee.ToLogger())
	//paso 1
	result, err = i.service.RegisterEmployee(ctx, employee)
	if err != nil {
		if err == domain.ErrIncompleteRegistration {
			i.logger.Warn(logger.LogEmployeeInteractorIncompleteDetected, "email", employee.Email)

			if cleanErr := i.service.CheckAndCleanInconsistentState(ctx, employee.Email); cleanErr != nil {
				i.logger.Error(logger.LogEmployeeInteractorCleanup_Error, "email", employee.Email, "error", cleanErr)
				return nil, cleanErr
			}

			i.logger.Success(logger.LogEmployeeInteractorCleanup_OK, "email", employee.Email)
			return nil, err
		}
		i.logger.Error(logger.LogEmployeeInteractorStep1_Error, "error", err)
		return
	}
	i.logger.Success(logger.LogEmployeeInteractorStep1_OK)

	employee.SetID()
	i.logger.Debug(logger.LogEmployeeInteractorIDGenerated, "employee_id", employee.ID)

	//paso 1.5
	if err = i.service.CheckAndCleanInconsistentState(ctx, employee.Email); err != nil {
		i.logger.Error(logger.LogEmployeeInteractorStep15_Error, "email", employee.Email, "error", err)
		return
	}
	i.logger.Success(logger.LogEmployeeInteractorStep15_OK, "email", employee.Email)

	//paso 2
	tx, err := i.service.BeginTx(ctx)
	if err != nil {
		i.logger.Error(logger.LogEmployeeInteractorStep2_Error, "error", err)
		return
	}
	i.logger.Success(logger.LogEmployeeInteractorStep2_OK)

	var keycloakUserID string
	var keycloakCreated bool

	defer func() {
		if err != nil {

			if rbErr := tx.Rollback(); rbErr != nil {
				i.logger.Error(logger.LogEmployeeInteractorRollbackDB_Error,
					"rollback_error", rbErr,
					"original_error", err)
			} else {
				i.logger.Warn(logger.LogEmployeeInteractorRollbackDB_OK)
			}

			if keycloakCreated {
				if kcErr := i.service.RollbackKeycloakUser(ctx, keycloakUserID); kcErr != nil {
					i.logger.Error(logger.LogEmployeeInteractorRollbackKeycloak_Err,
						"keycloak_error", kcErr,
						"keycloak_user_id", keycloakUserID)
				} else {
					i.logger.Warn(logger.LogEmployeeInteractorRollbackKeycloak_OK)
				}
			}
		}
	}()
	//paso 3
	if err = i.service.SaveEmployeeToDB(ctx, tx, employee); err != nil {
		i.logger.Error(logger.LogEmployeeInteractorStep3_Error, "error", err)
		return
	}
	i.logger.Success(logger.LogEmployeeInteractorStep3_OK)

	//paso 4
	keycloakUserID, err = i.service.CreateUserInKeycloak(ctx, &employee)
	if err != nil {
		i.logger.Error(logger.LogEmployeeInteractorStep4_Error, "error", err)
		err = domain.ErrKeycloakUserCreationFailed
		return
	}
	keycloakCreated = true
	i.logger.Success(logger.LogEmployeeInteractorStep4_OK, "keycloak_user_id", keycloakUserID)

	//paso 5
	if err = i.service.SetUserPassword(ctx, keycloakUserID, employee.Password); err != nil {
		i.logger.Error(logger.LogEmployeeInteractorStep5_Error, "error", err)
		return
	}
	i.logger.Success(logger.LogEmployeeInteractorStep5_OK)

	//paso 6
	if err = i.service.AssignUserRole(ctx, keycloakUserID, employee.Role); err != nil {
		i.logger.Error(logger.LogEmployeeInteractorStep6_Error, "error", err)
		return
	}
	i.logger.Success(logger.LogEmployeeInteractorStep6_OK, "role", employee.Role)

	//paso 7
	if err = i.service.UpdateEmployeeKeycloakID(ctx, tx, employee.ID, keycloakUserID); err != nil {
		i.logger.Error(logger.LogEmployeeInteractorStep7_Error, "error", err)
		return
	}
	i.logger.Success(logger.LogEmployeeInteractorStep7_OK)

	//paso 8
	if err = tx.Commit(); err != nil {
		i.logger.Error(logger.LogEmployeeInteractorCommit_Error, "error", err)
		return
	}
	i.logger.Success(logger.LogEmployeeInteractorCommit_OK)

	employee.KeycloakUserID = keycloakUserID
	result.Employee = employee
	result.Message = "user registered successfully"

	i.logger.Success(logger.LogEmployeeInteractorRegComplete, employee.ToLogger(),
		"Keycloak_id", keycloakUserID)
	err = nil

	return 
}

func (i *Interactor) Locate(ctx context.Context, id string) (*dto.RegisterEmployee, error) {
	result, err := i.service.LocateEmployee(ctx, id)
	if err != nil {
		i.logger.Error(logger.LogEmployeeLocateError, "error", err)
		return nil, err
	}
	return result, nil
}
