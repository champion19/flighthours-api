package output

import (
	"context"

	"github.com/champion19/flighthours-api/core/interactor/services/domain"
)

type Repository interface {
	Save(ctx context.Context, employee domain.Employee) error
	GetEmployeeByEmail(ctx context.Context, email string) (*domain.Employee, error)
	GetEmployeeByID(ctx context.Context, id string) (*domain.Employee, error)
	UpdateEmployee(ctx context.Context, employee domain.Employee) error
	DeleteEmployee(ctx context.Context, id string) error
}
