package interactor

import (
	"context"

	"github.com/champion19/flighthours-api/core/interactor/dto"
	"github.com/champion19/flighthours-api/core/interactor/services/domain"
	"github.com/champion19/flighthours-api/core/ports/input"
	"github.com/champion19/flighthours-api/middleware"
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
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := i.logger.WithTraceID(traceID)

	log.Info(logger.LogEmployeeInteractorRegStart, employee.ToLogger())
	//paso 1
	result, err = i.service.RegisterEmployee(ctx, employee)
	if err != nil {
		if err == domain.ErrIncompleteRegistration {
			log.Warn(logger.LogEmployeeInteractorIncompleteDetected, "email", employee.Email)

			if cleanErr := i.service.CheckAndCleanInconsistentState(ctx, employee.Email); cleanErr != nil {
				log.Error(logger.LogEmployeeInteractorCleanup_Error, "email", employee.Email, "error", cleanErr)
				return nil, cleanErr
			}

			log.Success(logger.LogEmployeeInteractorCleanup_OK, "email", employee.Email)
			return nil, err
		}
		log.Error(logger.LogEmployeeInteractorStep1_Error, "error", err)
		return
	}
	log.Success(logger.LogEmployeeInteractorStep1_OK)

	employee.SetID()
	log.Debug(logger.LogEmployeeInteractorIDGenerated, "employee_id", employee.ID)

	//paso 1.5
	if err = i.service.CheckAndCleanInconsistentState(ctx, employee.Email); err != nil {
		log.Error(logger.LogEmployeeInteractorStep15_Error, "email", employee.Email, "error", err)
		return
	}
	log.Success(logger.LogEmployeeInteractorStep15_OK, "email", employee.Email)

	//paso 2
	tx, err := i.service.BeginTx(ctx)
	if err != nil {
		log.Error(logger.LogEmployeeInteractorStep2_Error, "error", err)
		return
	}
	log.Success(logger.LogEmployeeInteractorStep2_OK)

	var keycloakUserID string
	var keycloakCreated bool

	defer func() {
		if err != nil {

			if rbErr := tx.Rollback(); rbErr != nil {
				log.Error(logger.LogEmployeeInteractorRollbackDB_Error,
					"rollback_error", rbErr,
					"original_error", err)
			} else {
				log.Warn(logger.LogEmployeeInteractorRollbackDB_OK)
			}

			if keycloakCreated {
				if kcErr := i.service.RollbackKeycloakUser(ctx, keycloakUserID); kcErr != nil {
					log.Error(logger.LogEmployeeInteractorRollbackKeycloak_Err,
						"keycloak_error", kcErr,
						"keycloak_user_id", keycloakUserID)
				} else {
					log.Warn(logger.LogEmployeeInteractorRollbackKeycloak_OK)
				}
			}
		}
	}()
	//paso 3
	if err = i.service.SaveEmployeeToDB(ctx, tx, employee); err != nil {
		log.Error(logger.LogEmployeeInteractorStep3_Error, "error", err)
		return
	}
	log.Success(logger.LogEmployeeInteractorStep3_OK)

	//paso 4
	keycloakUserID, err = i.service.CreateUserInKeycloak(ctx, &employee)
	if err != nil {
		log.Error(logger.LogEmployeeInteractorStep4_Error, "error", err)
		err = domain.ErrKeycloakUserCreationFailed
		return
	}
	keycloakCreated = true
	log.Success(logger.LogEmployeeInteractorStep4_OK, "keycloak_user_id", keycloakUserID)

	//paso 5
	if err = i.service.SetUserPassword(ctx, keycloakUserID, employee.Password); err != nil {
		log.Error(logger.LogEmployeeInteractorStep5_Error, "error", err)
		return
	}
	log.Success(logger.LogEmployeeInteractorStep5_OK)

	//paso 6
	if err = i.service.AssignUserRole(ctx, keycloakUserID, employee.Role); err != nil {
		log.Error(logger.LogEmployeeInteractorStep6_Error, "error", err)
		return
	}
	log.Success(logger.LogEmployeeInteractorStep6_OK, "role", employee.Role)

	//paso 7
	if err = i.service.UpdateEmployeeKeycloakID(ctx, tx, employee.ID, keycloakUserID); err != nil {
		log.Error(logger.LogEmployeeInteractorStep7_Error, "error", err)
		return
	}
	log.Success(logger.LogEmployeeInteractorStep7_OK)

	//paso 8
	if err = tx.Commit(); err != nil {
		log.Error(logger.LogEmployeeInteractorCommit_Error, "error", err)
		return
	}
	log.Success(logger.LogEmployeeInteractorCommit_OK)

	employee.KeycloakUserID = keycloakUserID
	result.Employee = employee
	result.Message = "user registered successfully"

	log.Success(logger.LogEmployeeInteractorRegComplete, employee.ToLogger(),
		"Keycloak_id", keycloakUserID)
	err = nil

	return
}

func (i *Interactor) Locate(ctx context.Context, id string) (*dto.RegisterEmployee, error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := i.logger.WithTraceID(traceID)

	result, err := i.service.LocateEmployee(ctx, id)
	if err != nil {
		log.Error(logger.LogEmployeeLocateError, "error", err)
		return nil, err
	}
	return result, nil
}
