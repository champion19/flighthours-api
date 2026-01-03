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

	//employee operations - transactional
	Save(ctx context.Context, tx Tx, employee domain.Employee) error
	UpdateEmployee(ctx context.Context, tx Tx, employee domain.Employee) error
	PatchEmployee(ctx context.Context, tx Tx, id string, keycloakUserID string) error
	DeleteEmployee(ctx context.Context, tx Tx, id string) error

	//employee operations - read
	GetEmployeeByEmail(ctx context.Context, email string) (*domain.Employee, error)
	GetEmployeeByID(ctx context.Context, id string) (*domain.Employee, error)
	GetEmployeeByKeycloakID(ctx context.Context, keycloakUserID string) (*domain.Employee, error)
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

// AirlineRepository defines the interface for airline data persistence
type AirlineRepository interface {
	BeginTx(ctx context.Context) (Tx, error)

	// Airline operations - read
	GetAirlineByID(ctx context.Context, id string) (*domain.Airline, error)

	// Airline operations - transactional
	UpdateAirlineStatus(ctx context.Context, tx Tx, id string, status bool) error
}

// AirportRepository defines the interface for airport data persistence
type AirportRepository interface {
	BeginTx(ctx context.Context) (Tx, error)

	// Airport operations - read
	GetAirportByID(ctx context.Context, id string) (*domain.Airport, error)

	// Airport operations - transactional
	UpdateAirportStatus(ctx context.Context, tx Tx, id string, status bool) error
}
