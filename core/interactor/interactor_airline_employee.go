package interactor

import (
	"context"

	"github.com/champion19/flighthours-api/core/interactor/services/domain"
	"github.com/champion19/flighthours-api/core/ports/input"
	"github.com/champion19/flighthours-api/middleware"
	"github.com/champion19/flighthours-api/platform/logger"
)

// AirlineEmployeeInteractor orchestrates airline employee operations (Release 15)
type AirlineEmployeeInteractor struct {
	service input.AirlineEmployeeService
	logger  logger.Logger
}

// NewAirlineEmployeeInteractor creates a new airline employee interactor
func NewAirlineEmployeeInteractor(service input.AirlineEmployeeService, log logger.Logger) *AirlineEmployeeInteractor {
	return &AirlineEmployeeInteractor{
		service: service,
		logger:  log,
	}
}

// GetAirlineEmployeeByID retrieves an airline employee by its ID (HU26)
func (i *AirlineEmployeeInteractor) GetAirlineEmployeeByID(ctx context.Context, id string) (*domain.AirlineEmployee, error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := i.logger.WithTraceID(traceID)

	log.Info(logger.LogDatabaseAvailable, "operation", "get_airline_employee", "employee_id", id)

	employee, err := i.service.GetAirlineEmployeeByID(ctx, id)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, "operation", "get_airline_employee", "employee_id", id, "error", err)
		return nil, err
	}

	log.Success(logger.LogDatabaseAvailable, employee.ToLogger())
	return employee, nil
}

// ListAirlineEmployees retrieves all airline employees with optional filters
func (i *AirlineEmployeeInteractor) ListAirlineEmployees(ctx context.Context, filters map[string]interface{}) ([]domain.AirlineEmployee, error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := i.logger.WithTraceID(traceID)

	log.Info(logger.LogDatabaseAvailable, "operation", "list_airline_employees", "filters", filters)

	employees, err := i.service.ListAirlineEmployees(ctx, filters)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, "operation", "list_airline_employees", "error", err)
		return nil, err
	}

	log.Success(logger.LogDatabaseAvailable, "operation", "list_airline_employees", "count", len(employees))
	return employees, nil
}

// CreateAirlineEmployee creates a new airline employee (HU28)
func (i *AirlineEmployeeInteractor) CreateAirlineEmployee(ctx context.Context, employee domain.AirlineEmployee) (*domain.AirlineEmployee, error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := i.logger.WithTraceID(traceID)

	log.Info(logger.LogDatabaseAvailable, "operation", "create_airline_employee", "email", employee.Email)

	// Generate ID
	employee.SetID()

	if err := i.service.CreateAirlineEmployee(ctx, employee); err != nil {
		log.Error(logger.LogDatabaseUnavailable, "operation", "create_airline_employee", "error", err)
		return nil, err
	}

	log.Success(logger.LogDatabaseAvailable, employee.ToLogger())
	return &employee, nil
}

// UpdateAirlineEmployee updates an existing airline employee (HU27)
func (i *AirlineEmployeeInteractor) UpdateAirlineEmployee(ctx context.Context, id string, employee domain.AirlineEmployee) error {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := i.logger.WithTraceID(traceID)

	log.Info(logger.LogDatabaseAvailable, "operation", "update_airline_employee", "employee_id", id)

	// Verify employee exists first
	existing, err := i.service.GetAirlineEmployeeByID(ctx, id)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, "operation", "get_airline_employee", "employee_id", id, "error", err)
		return err
	}

	// Preserve protected fields
	employee.ID = existing.ID
	employee.Email = existing.Email
	employee.KeycloakUserID = existing.KeycloakUserID

	if err := i.service.UpdateAirlineEmployee(ctx, employee); err != nil {
		log.Error(logger.LogDatabaseUnavailable, "operation", "update_airline_employee", "employee_id", id, "error", err)
		return err
	}

	log.Success(logger.LogDatabaseAvailable, employee.ToLogger())
	return nil
}

// ActivateAirlineEmployee activates an airline employee (HU29)
func (i *AirlineEmployeeInteractor) ActivateAirlineEmployee(ctx context.Context, id string) error {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := i.logger.WithTraceID(traceID)

	log.Info(logger.LogDatabaseAvailable, "operation", "activate_airline_employee", "employee_id", id)

	if err := i.service.ActivateAirlineEmployee(ctx, id); err != nil {
		log.Error(logger.LogDatabaseUnavailable, "operation", "activate_airline_employee", "employee_id", id, "error", err)
		return err
	}

	log.Success(logger.LogDatabaseAvailable, "operation", "activate_airline_employee", "employee_id", id)
	return nil
}

// DeactivateAirlineEmployee deactivates an airline employee (HU30)
func (i *AirlineEmployeeInteractor) DeactivateAirlineEmployee(ctx context.Context, id string) error {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := i.logger.WithTraceID(traceID)

	log.Info(logger.LogDatabaseAvailable, "operation", "deactivate_airline_employee", "employee_id", id)

	if err := i.service.DeactivateAirlineEmployee(ctx, id); err != nil {
		log.Error(logger.LogDatabaseUnavailable, "operation", "deactivate_airline_employee", "employee_id", id, "error", err)
		return err
	}

	log.Success(logger.LogDatabaseAvailable, "operation", "deactivate_airline_employee", "employee_id", id)
	return nil
}
