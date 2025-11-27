package input

import (
	"context"

	"github.com/champion19/flighthours-api/core/interactor/dto"
	"github.com/champion19/flighthours-api/core/interactor/services/domain"
	"github.com/champion19/flighthours-api/core/ports/output"
)

type Service interface {
	BeginTx(ctx context.Context) (output.Tx, error)
	RegisterEmployee(ctx context.Context, employee domain.Employee) (*dto.RegisterEmployee, error)
	GetEmployeeByEmail(ctx context.Context, email string) (*domain.Employee, error)
  GetEmployeeByID(ctx context.Context,id string)(*domain.Employee,error)
	LocateEmployee(ctx context.Context, id string) (*dto.RegisterEmployee, error)
	SaveEmployeeToDB(ctx context.Context, tx output.Tx, employee domain.Employee) error
	CreateUserInKeycloak(ctx context.Context, employee *domain.Employee) (string, error)
	SetUserPassword(ctx context.Context, userID string, password string) error
	AssignUserRole(ctx context.Context, userID string, role string) error
	UpdateEmployeeKeycloakID(ctx context.Context, tx output.Tx, employeeID string, keycloakUserID string) error

	RollbackEmployee(ctx context.Context, employeeID string) error
	RollbackKeycloakUser(ctx context.Context, KeycloakUserID string) error
}
