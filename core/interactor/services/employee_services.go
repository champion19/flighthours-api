package services

import (
	"context"

	"github.com/champion19/flighthours-api/core/interactor/dto"
	"github.com/champion19/flighthours-api/core/interactor/services/domain"
	"github.com/champion19/flighthours-api/core/ports/input"
	"github.com/champion19/flighthours-api/core/ports/output"
	"github.com/champion19/flighthours-api/platform/logger"
)

type service struct {
	repository output.Repository
	keycloak   output.AuthClient
	logger     logger.Logger
}

func NewService(repository output.Repository, keycloak output.AuthClient, logger logger.Logger) input.Service {
	return &service{
		repository: repository,
		keycloak:   keycloak,
		logger:     logger,
	}
}

func (s service) GetEmployeeByEmail(ctx context.Context, email string) (*domain.Employee, error) {
	employee, err := s.repository.GetEmployeeByEmail(ctx, nil, email)
	if err != nil {
		s.logger.Error(logger.LogEmployeeGetByEmailError, err)
		return nil, err
	}
	return employee, nil
}

func (s service) GetEmployeeByID(ctx context.Context, id string) (*domain.Employee, error) {
	employee, err := s.repository.GetEmployeeByID(ctx, nil, id)
	if err != nil {
		s.logger.Error(logger.LogEmployeeGetByIDError, err)
		return nil, err
	}
	return employee, nil
}

func (s service) BeginTx(ctx context.Context) (output.Tx, error) {
	return s.repository.BeginTx(ctx)
}

func (s service) RegisterEmployee(ctx context.Context, employee domain.Employee) (*dto.RegisterEmployee, error) {
	existingEmployee, err := s.repository.GetEmployeeByEmail(ctx, nil, employee.Email)
	if err == nil && existingEmployee != nil {
		s.logger.Warn(logger.LogEmployeeExists, err)
		return nil, domain.ErrDuplicateUser
	}

	return &dto.RegisterEmployee{
		Employee: employee,
		Message:  "Employee registered successfully",
	}, nil
}

func (s service) SaveEmployeeToDB(ctx context.Context, tx output.Tx, employee domain.Employee) error {
	err := s.repository.Save(ctx, tx, employee)
	if err != nil {
		s.logger.Error(logger.LogEmployeeSaveError, err)
		return err
	}
	return nil
}

func (s service) CreateUserInKeycloak(ctx context.Context, employee *domain.Employee) (string, error) {
	keycloakUserID, err := s.keycloak.CreateUser(ctx, employee)
	if err != nil {
		s.logger.Error(logger.LogKeycloakUserCreateError, err)
		return "", err
	}
	return keycloakUserID, nil
}

func (s service) SetUserPassword(ctx context.Context, userID string, password string) error {
	err := s.keycloak.SetPassword(ctx, userID, password, true)
	if err != nil {
		s.logger.Error(logger.LogKeycloakPasswordSetError, err)
		return err
	}
	return nil
}

func (s service) AssignUserRole(ctx context.Context, userID string, role string) error {
	err := s.keycloak.AssignRole(ctx, userID, role)
	if err != nil {
		s.logger.Error(logger.LogKeycloakRoleAssignError, err)
		return err
	}
	return nil
}

func (s service) UpdateEmployeeKeycloakID(ctx context.Context, tx output.Tx, employeeID string, keycloakUserID string) error {
	err := s.repository.PatchEmployee(ctx, tx, employeeID, keycloakUserID)
	if err != nil {
		s.logger.Error(logger.LogEmployeeUpdateKeycloakIDError, err)
		return err
	}
	return nil
}

func (s service) RollbackEmployee(ctx context.Context, employeeID string) error {
	err := s.repository.DeleteEmployee(ctx, nil, employeeID)
	if err != nil {
		s.logger.Error(logger.LogEmployeeDeleteError, err)
		return err
	}
	return nil
}

func (s service) RollbackKeycloakUser(ctx context.Context, KeycloakUserID string) error {
	err := s.keycloak.DeleteUser(ctx, KeycloakUserID)
	if err != nil {
		s.logger.Error(logger.LogKeycloakUserDeleteError, err)
		return err
	}
	return nil
}

func (s service) LocateEmployee(ctx context.Context, id string) (*dto.RegisterEmployee, error) {
	employee, err := s.repository.GetEmployeeByID(ctx, nil, id)
	if err != nil {
		s.logger.Error(logger.LogEmployeeGetByIDError, err)
		return nil, err
	}

	if employee == nil {
		s.logger.Error(logger.LogEmployeeNotFound, err)
		return nil, err
	}

	return &dto.RegisterEmployee{
		Employee: *employee,
		Message:  "Employee located successfully",
	}, nil
}
