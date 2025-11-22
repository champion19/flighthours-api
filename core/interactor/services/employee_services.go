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
	s.logger.Debug("Getting employee by email", "email", email)
employee, err := s.repository.GetEmployeeByEmail(ctx,email)
if err != nil {
	s.logger.Error("Error getting employee by email","email",email,"error", err)
	return nil, err
}
s.logger.Debug("Employee found successfully","email",email,"employee_id",employee.ID)
return employee, nil
}

func (s service) GetEmployeeByID(ctx context.Context, id string) (*domain.Employee, error) {
	s.logger.Debug("Getting employee by id", "employee_id", id)
	employee,err:= s.repository.GetEmployeeByID(ctx,id)
	if err != nil {
		s.logger.Error("Error getting employee by id","employee_id",id,"error", err)
		return nil, err
	}
	s.logger.Debug("Employee found successfully","employee_id",id,"email",employee.Email)
	return employee, nil
}

func(s service) BeginTx(ctx context.Context) (output.Tx, error) {
	return s.repository.BeginTx(ctx)
}


func (s service) RegisterEmployee(ctx context.Context, employee domain.Employee) (*dto.RegisterEmployee, error) {
	s.logger.Info("Initiating employee registration", "email", employee.Email)
	existingEmployee, err := s.repository.GetEmployeeByEmail(ctx,employee.Email)
	if err == nil && existingEmployee != nil {
		s.logger.Warn("Employee already exists", err)
		return nil, domain.ErrDuplicateUser
	}
	s.logger.Info("Register validation passed")

	return &dto.RegisterEmployee{
		Employee: employee,
		Message:  "Employee registered successfully",
	}, nil
}



func (s service) SaveEmployeeToDB(ctx context.Context, tx output.Tx, employee domain.Employee) error {
	s.logger.Info("Saving employee to database", "email", employee.Email,"employee_id",employee.ID)
	err := s.repository.Save(ctx, tx, employee)
	if err != nil {
		s.logger.Error("Error saving employee to database","email",employee.Email,"error", err)
		return err
	}
	s.logger.Success("Employee saved successfully","email",employee.Email,"employee_id",employee.ID)
	return nil
}

func (s service) CreateUserInKeycloak(ctx context.Context, employee *domain.Employee) (string, error) {
	s.logger.Info("Creating user in keycloak", "email", employee.Email,"employee_id",employee.ID)
  userID, err := s.keycloak.CreateUser(ctx, employee)
	if err != nil {
		s.logger.Error("Error creating user in keycloak", "email", employee.Email,"employee_id","error", err)
		return "", err
	}
	s.logger.Success("User created in keycloak", "email", employee.Email,"employee_id",employee.ID)
	return userID, nil
}

func (s service) SetUserPassword(ctx context.Context, userID string, password string) error {
	s.logger.Info("Setting user password in keycloak", "user_id", userID)
	err := s.keycloak.SetPassword(ctx, userID, password, true)
	if err != nil {
		s.logger.Error("Error setting user password in keycloak", err)
		return err
	}
	s.logger.Success("User password set in keycloak", "user_id", userID)
	return nil
}

func (s service) AssignUserRole(ctx context.Context, userID string, role string) error {
	s.logger.Info("Assigning user role in keycloak", "KeycloakUserID", userID,"role",role)
	err := s.keycloak.AssignRole(ctx, userID, role)
	if err != nil {
		s.logger.Error("Error assigning user role in keycloak", "KeycloakUserID", userID,"role",role,"error", err)
		return err
	}
	s.logger.Success("User role assigned in keycloak", "KeycloakUserID", userID,"role",role)
	return nil
}

func (s service) UpdateEmployeeKeycloakID(ctx context.Context, tx output.Tx, employeeID string, keycloakUserID string) error {
	s.logger.Info("Updating employee keycloak id in database", "employee_id", employeeID,"keycloak_user_id",keycloakUserID)
	err := s.repository.PatchEmployee(ctx, tx, employeeID, keycloakUserID)
	if err != nil {
		s.logger.Error("Error updating employee keycloak id in database", "employee_id", employeeID,"error", err)
		return err
	}
	s.logger.Success("Employee keycloak id updated in database", "employee_id", employeeID,"keycloak_user_id",keycloakUserID)
	return nil
}

func (s service) RollbackEmployee(ctx context.Context,  employeeID string) error {
  s.logger.Warn("Executing rollback:deleting employee from database", "employee_id", employeeID)
	err := s.repository.DeleteEmployee(ctx,nil, employeeID)
	if err != nil {
		s.logger.Error("Error in rollback of employee", "employee_id", employeeID,"error", err)
		return err
	}
	s.logger.Info("Rollback of employee completed successfully", "employee_id", employeeID)
	return nil
}

func (s service) RollbackKeycloakUser(ctx context.Context, KeycloakUserID string) error {
	s.logger.Warn("Executing rollback:deleting user from keycloak", "KeycloakUserID", KeycloakUserID)
	err := s.keycloak.DeleteUser(ctx,KeycloakUserID)
	if err != nil {
		s.logger.Error("Error in rollback of user from keycloak", "KeycloakUserID", KeycloakUserID,"error", err)
		return err
	}
	s.logger.Info("Rollback of user completed successfully", "KeycloakUserID", KeycloakUserID)
	return nil
}

func (s service) LocateEmployee(ctx context.Context, id string) (*dto.RegisterEmployee, error) {
	employee, err := s.repository.GetEmployeeByID(ctx, id)
	if err != nil {
		s.logger.Error("Error getting employee by id", err)
		return nil, err
	}

	if employee == nil {
		s.logger.Error("Employee not found", err)
		return nil, err
	}

	return &dto.RegisterEmployee{
		Employee: *employee,
		Message:  "Employee located successfully",
	}, nil
}
