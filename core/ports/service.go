package ports

import (
	"github.com/Nerzal/gocloak/v13"
	"github.com/champion19/Flighthours_backend/core/domain"
	"github.com/champion19/Flighthours_backend/core/dto"
)
type Service interface {
	RegisterEmployee(employee domain.Employee) (*dto.RegisterEmployee, error)
	GetEmployeeByEmail(email string) (*domain.Employee, error)
	LoginEmployee(email, password string) (*gocloak.JWT, error)
}
