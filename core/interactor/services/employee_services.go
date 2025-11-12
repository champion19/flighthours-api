package services

import (
	"context"

	"github.com/champion19/flighthours-api/core/interactor/dto"
	"github.com/champion19/flighthours-api/core/interactor/services/domain"
	"github.com/champion19/flighthours-api/core/ports/input"
	"github.com/champion19/flighthours-api/core/ports/output"
)

type service struct {
	repository output.Repository
	keycloak   output.AuthClient
}

func NewService(repository output.Repository, keycloak output.AuthClient) input.Service {
	return &service{
		repository: repository,
		keycloak:   keycloak,
	}
}

func (s service) GetEmployeeByEmail(ctx context.Context, email string) (*domain.Employee, error) {
	return s.repository.GetEmployeeByEmail(email)
}



func (s service) RegisterEmployee(ctx context.Context, employee domain.Employee) (*dto.RegisterEmployee, error) {
	existingEmployee, err := s.repository.GetEmployeeByEmail(employee.Email)
	if err == nil && existingEmployee != nil {
		return nil, domain.ErrDuplicateUser
	}
	if employee.Role == "" {
		return nil, domain.ErrRoleRequired
	}

	//dto de respuesta
	return &dto.RegisterEmployee{
		Employee: employee,
		Message:  "Employee registered successfully",
	}, nil
}
func (s service) SaveEmployeeToDB(ctx context.Context, employee domain.Employee) error {
	return s.repository.Save(employee)
}

func (s service) CreateUserInKeycloak(ctx context.Context, employee *domain.Employee) (string, error) {
	return s.keycloak.CreateUser(ctx, employee)
}

func (s service) SetUserPassword(ctx context.Context, userID string, password string) error {
	return s.keycloak.SetPassword(ctx, userID, password, true)
}

func (s service) AssignUserRole(ctx context.Context, userID string, role string) error {
	return s.keycloak.AssignRole(ctx, userID, role)
}

func (s service) UpdateEmployeeKeycloakID(ctx context.Context, employeeID string, keycloakUserID string) error {
	return s.repository.PatchEmployee(employeeID, keycloakUserID)
}

func (s service) RollbackEmployee(ctx context.Context, employeeID string) error {
	return s.repository.DeleteEmployee(employeeID)
}

func (s service) RollbackKeycloakUser(ctx context.Context, KeycloakUserID string) error {
	return s.keycloak.DeleteUser(ctx, KeycloakUserID)
}

func (s service) LocateEmployee(ctx context.Context, id string) (*dto.RegisterEmployee, error) {
	employee, err := s.repository.GetEmployeeByID(id)
	if err != nil {
		return nil, err
	}

	if employee == nil {
		return nil, err
	}

	return &dto.RegisterEmployee{
		Employee: *employee,
		Message:  "Employee located successfully",
	}, nil
}
