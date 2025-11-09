package output

import (
	"github.com/champion19/flighthours-api/core/interactor/services/domain"
)

type Repository interface {
	Save(employee domain.Employee) error
	GetEmployeeByEmail(email string) (*domain.Employee, error)
	GetEmployeeByID(id string) (*domain.Employee, error)
	UpdateEmployee(employee domain.Employee) error
	PatchEmployee(id string, keycloakUserID string) error
	DeleteEmployee(id string) error
}
