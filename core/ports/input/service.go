package input

import (
	"context"

	"github.com/champion19/flighthours-api/core/interactor/dto"
	"github.com/champion19/flighthours-api/core/interactor/services/domain"
)

type Service interface {
	RegisterEmployee(ctx context.Context, employee domain.Employee) (*dto.RegisterEmployee, error)
	GetEmployeeByEmail(ctx context.Context, email string) (*domain.Employee, error)
  LocateEmployee(ctx context.Context, id string) (*dto.RegisterEmployee, error)
	SaveEmployeeToDB(ctx context.Context, employee domain.Employee) error
	CreateUserInKeycloak(ctx context.Context, employee *domain.Employee) (string, error)
	SetUserPassword(ctx context.Context, userID string, password string) error
	AssignUserRole(ctx context.Context, userID string, role string) error
	UpdateEmployeeKeycloakID(ctx context.Context, employeeID string, keycloakUserID string) error

	RollbackEmployee(ctx context.Context, employeeID string) error
	RollbackKeycloakUser(ctx context.Context, KeycloakUserID string) error
}
