package output

import (
	"context"

	"github.com/champion19/flighthours-api/core/interactor/services/domain"
)

type Tx interface {
	Commit() error
	Rollback() error
}

type Repository interface {
	BeginTx(ctx context.Context) (Tx, error)
	Save(ctx context.Context, tx Tx, employee domain.Employee) error
	GetEmployeeByEmail(ctx context.Context, tx Tx, email string) (*domain.Employee, error)
	GetEmployeeByID(ctx context.Context, tx Tx, id string) (*domain.Employee, error)
	UpdateEmployee(ctx context.Context, tx Tx, employee domain.Employee) error
	PatchEmployee(ctx context.Context, tx Tx, id string, keycloakUserID string) error
	DeleteEmployee(ctx context.Context, tx Tx, id string) error
}
