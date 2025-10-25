package services

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Nerzal/gocloak/v13"
	"github.com/champion19/Flighthours_backend/config"
	"github.com/champion19/Flighthours_backend/core/domain"
	"github.com/champion19/Flighthours_backend/core/dto"
	"github.com/champion19/Flighthours_backend/core/ports"
)

type service struct {
	repository  ports.Repository
	authService ports.AuthorizationService
	config      *config.Config
}

func NewService(repo ports.Repository, authservice ports.AuthorizationService, cfg *config.Config) ports.Service {
	return &service{
		repository:  repo,
		authService: authservice,
		config:      cfg,
	}
}

func (s service) GetEmployeeByEmail(email string) (*domain.Employee, error) {
	return s.repository.GetEmployeeByEmail(email)
}

func (s service) RegisterEmployee(employee domain.Employee) (*dto.RegisterEmployee, error) {

	existingEmployee, err := s.repository.GetEmployeeByEmail(employee.Email)
	if err == nil && existingEmployee != nil {
		return nil, domain.ErrDuplicateUser
	}

	if employee.Role == "" {
		return nil, fmt.Errorf("role is required")
	}

	employee.SetID()

	plainPassword := employee.Password

	if err := employee.HashPassword(); err != nil {
		return nil, err
	}

	err = s.repository.Save(employee)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	token, err := s.syncUserWithKeycloak(ctx, &employee, plainPassword)
	if err != nil {

		slog.Error("Failed to sync user with Keycloak, rolling back",
			"user_id", employee.ID,
			"email", employee.Email,
			"error", err)

		if deleteErr := s.repository.DeleteEmployee(employee.ID); deleteErr != nil {
			slog.Error("Failed to rollback user creation",
				"user_id", employee.ID,
				"error", deleteErr)
		}

		return nil, fmt.Errorf("registration failed: %w", err)
	}

	slog.Info("User registered successfully",
		"user_id", employee.ID,
		"email", employee.Email,
		"role", employee.Role)

	return &dto.RegisterEmployee{
		Employee: employee,
		Token:    token,
	}, nil
}

func (s service) syncUserWithKeycloak(ctx context.Context, employee *domain.Employee, plainPassword string) (*gocloak.JWT, error) {
	if s.authService == nil {
		return nil, fmt.Errorf("keycloak service not configured")
	}

	keycloakUserID, err := s.authService.SyncUserToKeycloak(ctx, employee)
	if err != nil {
		return nil, fmt.Errorf("failed to sync user: %w", err)
	}

	err = s.authService.SetUserPassword(ctx, keycloakUserID, plainPassword)
	if err != nil {
		_ = s.authService.DeleteUserFromKeycloak(ctx, keycloakUserID)
		return nil, fmt.Errorf("failed to set password: %w", err)
	}

	err = s.authService.AssignRole(ctx, employee.ID, employee.Role)
	if err != nil {
		_ = s.authService.DeleteUserFromKeycloak(ctx, keycloakUserID)
		return nil, fmt.Errorf("failed to assign role: %w", err)
	}

	token, err := s.authService.LoginUser(ctx, employee.Email, plainPassword)
	if err != nil {
		slog.Warn("User created but login failed",
			"user_id", employee.ID,
			"error", err)
		return nil, fmt.Errorf("user created but failed to generate token: %w", err)
	}

	slog.Info("User synced with Keycloak successfully",
		"user_id", employee.ID,
		"email", employee.Email,
		"role", employee.Role,
		"keycloak_user_id", keycloakUserID)

	return token, nil
}

func (s service) LoginEmployee(email, password string) (*gocloak.JWT, error) {

	employee, err := s.repository.GetEmployeeByEmail(email)
	if err != nil {
		return nil, domain.ErrNotFoundUserByEmail
	}

	if !employee.Active {
		return nil, fmt.Errorf("user account is inactive")
	}

	ctx := context.Background()
	token, err := s.authService.LoginUser(ctx, email, password)
	if err != nil {
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	return token, nil
}
