package interactor

import (
	"context"

	"github.com/champion19/flighthours-api/core/interactor/services/domain"
	"github.com/champion19/flighthours-api/core/ports/input"
	"github.com/champion19/flighthours-api/middleware"
	"github.com/champion19/flighthours-api/platform/logger"
)

// AircraftRegistrationInteractor orchestrates aircraft registration operations
type AircraftRegistrationInteractor struct {
	service input.AircraftRegistrationService
	logger  logger.Logger
}

// NewAircraftRegistrationInteractor creates a new aircraft registration interactor
func NewAircraftRegistrationInteractor(service input.AircraftRegistrationService, log logger.Logger) *AircraftRegistrationInteractor {
	return &AircraftRegistrationInteractor{
		service: service,
		logger:  log,
	}
}

// GetAircraftRegistrationByID retrieves an aircraft registration by its ID
func (i *AircraftRegistrationInteractor) GetAircraftRegistrationByID(ctx context.Context, id string) (*domain.AircraftRegistration, error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := i.logger.WithTraceID(traceID)

	log.Info(logger.LogAircraftRegistrationGet, "registration_id", id)

	registration, err := i.service.GetAircraftRegistrationByID(ctx, id)
	if err != nil {
		log.Error(logger.LogAircraftRegistrationGetError, "registration_id", id, "error", err)
		return nil, err
	}

	log.Success(logger.LogAircraftRegistrationGetOK, registration.ToLogger())
	return registration, nil
}

// ListAircraftRegistrations retrieves all aircraft registrations with optional filters
func (i *AircraftRegistrationInteractor) ListAircraftRegistrations(ctx context.Context, filters map[string]interface{}) ([]domain.AircraftRegistration, error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := i.logger.WithTraceID(traceID)

	log.Info(logger.LogAircraftRegistrationList, "filters", filters)

	registrations, err := i.service.ListAircraftRegistrations(ctx, filters)
	if err != nil {
		log.Error(logger.LogAircraftRegistrationListError, "error", err)
		return nil, err
	}

	log.Success(logger.LogAircraftRegistrationListOK, "count", len(registrations))
	return registrations, nil
}

// CreateAircraftRegistration creates a new aircraft registration
func (i *AircraftRegistrationInteractor) CreateAircraftRegistration(ctx context.Context, registration domain.AircraftRegistration) error {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := i.logger.WithTraceID(traceID)

	log.Info(logger.LogAircraftRegistrationCreate, registration.ToLogger())

	if err := i.service.CreateAircraftRegistration(ctx, registration); err != nil {
		log.Error(logger.LogAircraftRegistrationCreateError, "error", err)
		return err
	}

	log.Success(logger.LogAircraftRegistrationCreateOK, registration.ToLogger())
	return nil
}

// UpdateAircraftRegistration updates an existing aircraft registration
func (i *AircraftRegistrationInteractor) UpdateAircraftRegistration(ctx context.Context, registration domain.AircraftRegistration) error {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := i.logger.WithTraceID(traceID)

	log.Info(logger.LogAircraftRegistrationUpdate, registration.ToLogger())

	// Verify registration exists
	_, err := i.service.GetAircraftRegistrationByID(ctx, registration.ID)
	if err != nil {
		log.Error(logger.LogAircraftRegistrationNotFound, "registration_id", registration.ID)
		return err
	}

	if err = i.service.UpdateAircraftRegistration(ctx, registration); err != nil {
		log.Error(logger.LogAircraftRegistrationUpdateError, "error", err)
		return err
	}

	log.Success(logger.LogAircraftRegistrationUpdateOK, registration.ToLogger())
	return nil
}
