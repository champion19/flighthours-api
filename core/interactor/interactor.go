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
		i.log.Error("failed to register employee", err)
		return nil, err
	}
	employee.SetID()
	i.log.Success("employee registered successfully")

	tx, err := i.service.BeginTx(ctx)
	if err != nil {
		i.log.Error("Error beginning transaction")
	}
   var keycloakUserID string
	var keycloakCreated bool

	defer func() {
		if err != nil {

			if rbErr := tx.Rollback(); rbErr != nil {
				i.log.Error("ROLLBACK BD FALLÓ - ALERTA CRÍTICA",
					"rollback_error", rbErr,
					"original_error", err)
			} else {
				i.log.Warn("Rollback BD ejecutado correctamente")
			}

			if keycloakCreated {
				if kcErr := i.service.RollbackKeycloakUser(ctx, keycloakUserID); kcErr != nil {
					i.log.Error("ROLLBACK KEYCLOAK FALLÓ - ALERTA CRÍTICA",
						"keycloak_error", kcErr,
						"keycloak_user_id", keycloakUserID)
						} else {
					i.log.Warn("Rollback Keycloak ejecutado correctamente")
				}
			}
		}
	}()

	if err = i.service.SaveEmployeeToDB(ctx, tx, employee); err != nil {
		i.log.Error("failed to save employee in database")
		return nil, err
	}

	keycloakUserID, err = i.service.CreateUserInKeycloak(ctx, &employee)
	if err != nil {
		i.log.Error("failed to create user in keycloak")
		return nil, err
	}

	keycloakCreated = true
	if err = i.service.SetUserPassword(ctx, keycloakUserID, employee.Password); err != nil {
		i.log.Error("failed to set user password in keycloak")
		return nil, err
	}

	if err = i.service.AssignUserRole(ctx, keycloakUserID, employee.Role); err != nil {
		i.log.Error("failed to assign user role in keycloak")
		return nil, err
	}

	if err = i.service.UpdateEmployeeKeycloakID(ctx, tx, employee.ID, keycloakUserID); err != nil {
		i.log.Error("failed to update employee keycloak id in database")
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		i.log.Error("failed to commit transaction")
		return nil, err
	}

	employee.KeycloakUserID = keycloakUserID
	result.Employee = employee
	result.Message = "user registered successfully"

	i.log.Success("user registered successfully",
		"email", employee.Email,
		"id", employee.ID,
		"Keycloak_id", keycloakUserID)


		err = nil
	return result, nil

}

func (i *Interactor) Locate(ctx context.Context, id string) (*dto.RegisterEmployee, error) {
	result, err := i.service.LocateEmployee(ctx, id)
	if err != nil {
		i.log.Error("failed to locate employee")
		return nil, err
	}
	return result, nil
}
