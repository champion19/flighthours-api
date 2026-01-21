package services

import (
	"context"

	"github.com/champion19/flighthours-api/core/interactor/services/domain"
	"github.com/champion19/flighthours-api/core/ports/output"
	"github.com/champion19/flighthours-api/platform/logger"
)

// airlineEmployeeService implements the AirlineEmployeeService interface
type airlineEmployeeService struct {
	repository output.AirlineEmployeeRepository
	logger     logger.Logger
}

// NewAirlineEmployeeService creates a new airline employee service
func NewAirlineEmployeeService(repository output.AirlineEmployeeRepository, log logger.Logger) *airlineEmployeeService {
	return &airlineEmployeeService{
		repository: repository,
		logger:     log,
	}
}

// BeginTx starts a new database transaction
func (s *airlineEmployeeService) BeginTx(ctx context.Context) (output.Tx, error) {
	return s.repository.BeginTx(ctx)
}

// GetAirlineEmployeeByID retrieves an airline employee by ID (HU26)
func (s *airlineEmployeeService) GetAirlineEmployeeByID(ctx context.Context, id string) (*domain.AirlineEmployee, error) {
	employee, err := s.repository.GetAirlineEmployeeByID(ctx, id)
	if err != nil {
		s.logger.Debug("GetAirlineEmployeeByID: error", "id", id, "error", err)
		return nil, err
	}
	return employee, nil
}

// ListAirlineEmployees retrieves all airline employees with optional filters
func (s *airlineEmployeeService) ListAirlineEmployees(ctx context.Context, filters map[string]interface{}) ([]domain.AirlineEmployee, error) {
	employees, err := s.repository.ListAirlineEmployees(ctx, filters)
	if err != nil {
		s.logger.Debug("ListAirlineEmployees: error", "filters", filters, "error", err)
		return nil, err
	}
	return employees, nil
}

// CreateAirlineEmployee creates a new airline employee (HU28)
func (s *airlineEmployeeService) CreateAirlineEmployee(ctx context.Context, employee domain.AirlineEmployee) error {
	tx, err := s.BeginTx(ctx)
	if err != nil {
		s.logger.Error(logger.LogDatabaseUnavailable, "error", err)
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Generate ID if not provided
	if employee.ID == "" {
		employee.SetID()
	}

	// Save employee to database
	if err = s.repository.SaveAirlineEmployee(ctx, tx, employee); err != nil {
		s.logger.Error(logger.LogDatabaseUnavailable, "operation", "save", "employee_id", employee.ID, "error", err)
		return err
	}

	if err = tx.Commit(); err != nil {
		s.logger.Error(logger.LogDatabaseUnavailable, "operation", "commit", "employee_id", employee.ID, "error", err)
		return err
	}

	s.logger.Info(logger.LogDatabaseAvailable, "operation", "create_airline_employee", "employee_id", employee.ID)
	return nil
}

// UpdateAirlineEmployee updates an existing airline employee (HU27)
func (s *airlineEmployeeService) UpdateAirlineEmployee(ctx context.Context, employee domain.AirlineEmployee) error {
	tx, err := s.BeginTx(ctx)
	if err != nil {
		s.logger.Error(logger.LogDatabaseUnavailable, "error", err)
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Update employee in database
	if err = s.repository.UpdateAirlineEmployee(ctx, tx, employee); err != nil {
		s.logger.Error(logger.LogDatabaseUnavailable, "operation", "update", "employee_id", employee.ID, "error", err)
		return err
	}

	if err = tx.Commit(); err != nil {
		s.logger.Error(logger.LogDatabaseUnavailable, "operation", "commit", "employee_id", employee.ID, "error", err)
		return err
	}

	s.logger.Info(logger.LogDatabaseAvailable, "operation", "update_airline_employee", "employee_id", employee.ID)
	return nil
}

// ActivateAirlineEmployee activates an airline employee (HU29)
func (s *airlineEmployeeService) ActivateAirlineEmployee(ctx context.Context, id string) error {
	// First verify the employee exists
	_, err := s.repository.GetAirlineEmployeeByID(ctx, id)
	if err != nil {
		return err
	}

	tx, err := s.BeginTx(ctx)
	if err != nil {
		s.logger.Error(logger.LogDatabaseUnavailable, "error", err)
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	if err = s.repository.UpdateAirlineEmployeeStatus(ctx, tx, id, true); err != nil {
		s.logger.Error(logger.LogDatabaseUnavailable, "operation", "activate", "employee_id", id, "error", err)
		return err
	}

	if err = tx.Commit(); err != nil {
		s.logger.Error(logger.LogDatabaseUnavailable, "operation", "commit", "employee_id", id, "error", err)
		return err
	}

	s.logger.Info(logger.LogDatabaseAvailable, "operation", "activate_airline_employee", "employee_id", id)
	return nil
}

// DeactivateAirlineEmployee deactivates an airline employee (HU30)
func (s *airlineEmployeeService) DeactivateAirlineEmployee(ctx context.Context, id string) error {
	// First verify the employee exists
	_, err := s.repository.GetAirlineEmployeeByID(ctx, id)
	if err != nil {
		return err
	}

	tx, err := s.BeginTx(ctx)
	if err != nil {
		s.logger.Error(logger.LogDatabaseUnavailable, "error", err)
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	if err = s.repository.UpdateAirlineEmployeeStatus(ctx, tx, id, false); err != nil {
		s.logger.Error(logger.LogDatabaseUnavailable, "operation", "deactivate", "employee_id", id, "error", err)
		return err
	}

	if err = tx.Commit(); err != nil {
		s.logger.Error(logger.LogDatabaseUnavailable, "operation", "commit", "employee_id", id, "error", err)
		return err
	}

	s.logger.Info(logger.LogDatabaseAvailable, "operation", "deactivate_airline_employee", "employee_id", id)
	return nil
}
