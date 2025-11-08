package services

import (
	"context"
	"fmt"

	"github.com/champion19/flighthours-api/core/interactor/dto"
	"github.com/champion19/flighthours-api/core/interactor/services/domain"
	"github.com/champion19/flighthours-api/core/ports/input"
	"github.com/champion19/flighthours-api/core/ports/output"
)

type Interactor struct {
	Service    input.Service
	Repository output.Repository
}

func NewInteractor(repo output.Repository, srv input.Service) *Interactor {
	return &Interactor{
		Service:    srv,
		Repository: repo,
	}
}

func (i *Interactor) Execute(employee domain.Employee) (*dto.RegisterEmployee, error) {
	// 1️⃣ Validaciones básicas del dominio o entrada
	if employee.Email == "" || employee.Name == "" {
		return nil, domain.ErrDuplicateUser
	}
	if employee.Role == "" {
		return nil, domain.ErrRoleAssignmentFailed
	}

	// 2️⃣ Lógica del caso de uso → delegar al servicio
	result, err := i.Service.RegisterEmployee(context.Background(),employee)
	if err != nil {
		return nil, domain.ErrUserCannotFound
	}

	// 3️⃣ Si todo fue bien, devolver el DTO con mensaje
	result.Message = fmt.Sprintf("Employee %s registered successfully", employee.Email)
	return result, nil

}
