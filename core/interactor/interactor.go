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
	var (
		employeeSaved      bool
		keycloakUserID   string
	)

  result,err:=i.service.RegisterEmployee(ctx, employee)
	if err != nil {
		return nil, err
	}
	employee.SetID()

  err =i.service.SaveEmployeeToDB(ctx, employee)
	if err != nil {
		return nil, err
	}
	employeeSaved	= true

	keycloakUserID,err = i.service.CreateUserInKeycloak(ctx, &employee)
	if err != nil {
		if employeeSaved {
			_= i.service.RollbackEmployee(ctx, employee.ID)
		}
		return nil, err
	}

	err = i.service.SetUserPassword(ctx, keycloakUserID, employee.Password)
	if err != nil {
		_ = i.service.RollbackKeycloakUser(ctx, keycloakUserID)
		_ = i.service.RollbackEmployee(ctx, employee.ID)
		return nil, err
	}

	err = i.service.AssignUserRole(ctx, keycloakUserID, employee.Role)
	if err != nil {
		_ = i.service.RollbackKeycloakUser(ctx, keycloakUserID)
		_ = i.service.RollbackEmployee(ctx, employee.ID)
		return nil, err
	}

  err = i.service.UpdateEmployeeKeycloakID(ctx, employee.ID, keycloakUserID)
	if err != nil {
		_ = i.service.RollbackKeycloakUser(ctx, keycloakUserID)
		_ = i.service.RollbackEmployee(ctx, employee.ID)
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
