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

type MessageRepository interface {
	BeginTx(ctx context.Context) (Tx, error)

	// Message operations - transactional
	SaveMessage(ctx context.Context, tx Tx, message domain.Message) error
	UpdateMessage(ctx context.Context, tx Tx, message domain.Message) error
	DeleteMessage(ctx context.Context, tx Tx, id string) error

	// Message operations - read
	GetAllActive(ctx context.Context) ([]domain.Message, error)
	GetByID(ctx context.Context, id string) (*domain.Message, error)
	GetByCode(ctx context.Context, code string) (*domain.Message, error)
	GetByType(ctx context.Context, msgType string) ([]domain.Message, error)
	GetByModule(ctx context.Context, module string) ([]domain.Message, error)

}
