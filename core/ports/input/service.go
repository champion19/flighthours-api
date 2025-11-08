package input

import (
	"context"

	"github.com/Nerzal/gocloak/v13"
	"github.com/champion19/flighthours-api/core/interactor/dto"
	"github.com/champion19/flighthours-api/core/interactor/services/domain"
)
type Service interface {
	RegisterEmployee(ctx context.Context,employee domain.Employee) (*dto.RegisterEmployee, error)
	GetEmployeeByEmail(ctx context.Context,email string) (*domain.Employee, error)
	LoginEmployee(ctx context.Context,email, password string) (*gocloak.JWT, error)
}
