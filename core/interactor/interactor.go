package interactor

import (
	"context"


	"github.com/champion19/flighthours-api/core/interactor/dto"
	"github.com/champion19/flighthours-api/core/interactor/services/domain"
	"github.com/champion19/flighthours-api/core/ports/input"

)

type Interactor struct {
	service    input.Service
}

func NewInteractor(service input.Service) *Interactor {
	return &Interactor{
		service:    service,
	}
}
func (i *Interactor) RegisterEmployee(ctx context.Context, employee domain.Employee) (*dto.RegisterEmployee, error) {
  result,err:=i.service.RegisterEmployee(ctx, employee)
	if err != nil {
		return nil, err
	}
	employee.SetID()

  tx,err:=i.service.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	if err = i.service.SaveEmployeeToDB(ctx, tx, employee); err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	keycloakUserID,err := i.service.CreateUserInKeycloak(ctx, &employee)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	err = i.service.SetUserPassword(ctx, keycloakUserID, employee.Password)
	if err != nil {
		_ = i.service.RollbackKeycloakUser(ctx, keycloakUserID)
		_ = tx.Rollback()
		return nil, err
	}

	err = i.service.AssignUserRole(ctx, keycloakUserID, employee.Role)
	if err != nil {
		_ = i.service.RollbackKeycloakUser(ctx, keycloakUserID)
		_ = tx.Rollback()
		return nil, err
	}


	err = i.service.SetUserPassword(ctx, keycloakUserID, employee.Password)
	if err != nil {
		_ = i.service.RollbackKeycloakUser(ctx, keycloakUserID)
		_ = tx.Rollback()
		return nil, err
	}

	err = i.service.AssignUserRole(ctx, keycloakUserID, employee.Role)
	if err != nil {
		_ = i.service.RollbackKeycloakUser(ctx, keycloakUserID)
		_ = tx.Rollback()
		return nil, err
	}

  err = i.service.UpdateEmployeeKeycloakID(ctx, tx, employee.ID, keycloakUserID)
	if err != nil {
		_ = i.service.RollbackKeycloakUser(ctx, keycloakUserID)
		_ = tx.Rollback()
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		_ = tx.Rollback()
		return nil, err
	}



	employee.KeycloakUserID = keycloakUserID
	result.Employee=employee
	result.Message="user registered successfully"

	return result, nil
}


func (i *Interactor) Locate(ctx context.Context, id string) (*dto.RegisterEmployee, error) {
	result, err := i.service.LocateEmployee(ctx, id)
	if err != nil {
		return nil, err
	}
	return result, nil
}
