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
	log     logger.Logger
}

func NewInteractor(service input.Service, log logger.Logger) *Interactor {
	return &Interactor{
		service: service,
		log:     log,
	}
}
func (i *Interactor) RegisterEmployee(ctx context.Context, employee domain.Employee) (*dto.RegisterEmployee, error) {
	result, err := i.service.RegisterEmployee(ctx, employee)
	if err != nil {
		i.log.Error(logger.LogEmployeeRegisterError, err)
		return nil, err
	}
	employee.SetID()
	i.log.Success(logger.LogEmployeeRegisterSuccess)

	tx, err := i.service.BeginTx(ctx)
	if err != nil {
		i.log.Error(logger.LogDBTransactionBeginErr)
		return nil, err
	}
	if err = i.service.SaveEmployeeToDB(ctx, tx, employee); err != nil {
		i.log.Error(logger.LogEmployeeSaveError)
		_ = tx.Rollback()
		return nil, err
	}

	keycloakUserID, err := i.service.CreateUserInKeycloak(ctx, &employee)
	if err != nil {
		i.log.Error(logger.LogKeycloakUserCreateError)
		_ = tx.Rollback()
		return nil, err
	}

	err = i.service.SetUserPassword(ctx, keycloakUserID, employee.Password)
	if err != nil {
		i.log.Error(logger.LogKeycloakPasswordSetError)
		_ = i.service.RollbackKeycloakUser(ctx, keycloakUserID)
		_ = tx.Rollback()
		return nil, err
	}

	err = i.service.AssignUserRole(ctx, keycloakUserID, employee.Role)
	if err != nil {
		i.log.Error(logger.LogKeycloakRoleAssignError)
		_ = i.service.RollbackKeycloakUser(ctx, keycloakUserID)
		_ = tx.Rollback()
		return nil, err
	}

	err = i.service.SetUserPassword(ctx, keycloakUserID, employee.Password)
	if err != nil {
		i.log.Error(logger.LogKeycloakPasswordSetError)
		_ = i.service.RollbackKeycloakUser(ctx, keycloakUserID)
		_ = tx.Rollback()
		return nil, err
	}

	err = i.service.AssignUserRole(ctx, keycloakUserID, employee.Role)
	if err != nil {
		i.log.Error(logger.LogKeycloakRoleAssignError)
		_ = i.service.RollbackKeycloakUser(ctx, keycloakUserID)
		_ = tx.Rollback()
		return nil, err
	}

	err = i.service.UpdateEmployeeKeycloakID(ctx, tx, employee.ID, keycloakUserID)
	if err != nil {
		i.log.Error(logger.LogEmployeeUpdateKeycloakIDError)
		_ = i.service.RollbackKeycloakUser(ctx, keycloakUserID)
		_ = tx.Rollback()
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		i.log.Error(logger.LogDBTransactionCommitErr)
		_ = tx.Rollback()
		return nil, err
	}

	employee.KeycloakUserID = keycloakUserID
	result.Employee = employee
	result.Message = "user registered successfully"

	i.log.Success(logger.LogRegSuccess,
		"email", employee.Email,
		"id", employee.ID,
		"Keycloak_id", keycloakUserID)

	return result, nil
}

func (i *Interactor) Locate(ctx context.Context, id string) (*dto.RegisterEmployee, error) {
	result, err := i.service.LocateEmployee(ctx, id)
	if err != nil {
		i.log.Error(logger.LogEmployeeLocateError)
		return nil, err
	}
	return result, nil
}
