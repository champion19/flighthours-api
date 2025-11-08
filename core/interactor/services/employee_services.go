package services

import (
	"context"
	"fmt"

	"github.com/Nerzal/gocloak/v13"
	"github.com/champion19/flighthours-api/config"
	"github.com/champion19/flighthours-api/core/interactor/dto"
	"github.com/champion19/flighthours-api/core/interactor/services/domain"
	"github.com/champion19/flighthours-api/core/ports/input"
	"github.com/champion19/flighthours-api/core/ports/output"
)

type service struct {
	repository output.Repository
	authClient output.AuthClient
	config     *config.Config
}

func NewService(repo output.Repository, authClient output.AuthClient, cfg *config.Config) input.Service {
	return &service{
		repository: repo,
		authClient: authClient,
		config:     cfg,
	}
}

func (s service) GetEmployeeByEmail(ctx context.Context,email string) (*domain.Employee, error) {
	return s.repository.GetEmployeeByEmail(ctx,email)
}

func (s service) RegisterEmployee(ctx context.Context, employee domain.Employee) (*dto.RegisterEmployee, error) {

	existingEmployee, err := s.repository.GetEmployeeByEmail(ctx, employee.Email)
	if err == nil && existingEmployee != nil {
		return nil, domain.ErrDuplicateUser
	}

	if employee.Role == "" {
		return nil, fmt.Errorf("role is required")
	}

	userID, err := s.authClient.CreateUser(ctx, &employee)
	if err != nil {
		return nil, fmt.Errorf("failed to create user in Keycloak: %w", err)
	}

	err = s.authClient.AssignRole(ctx, userID, employee.Role)
	if err != nil {
		_ = s.authClient.DeleteUser(ctx, userID)
		return nil, fmt.Errorf("failed to assign role in Keycloak: %w", err)
	}

	employee.KeycloakUserID = userID
	err = s.repository.Save(ctx, employee)
	if err != nil {
		_ = s.authClient.DeleteUser(ctx, userID)
		return nil, err
	}
	//dto de respuesta
	return &dto.RegisterEmployee{
		Employee: employee,
		Message:  "Employee registered successfully",
	}, nil
}

func (s service) LoginEmployee(ctx context.Context,email, password string) (*gocloak.JWT, error) {

	token, err := s.authClient.LoginUser(ctx, email, password)
	if err != nil {
		return nil, fmt.Errorf("login failed: %w", err)
	}
	return token, nil
}
