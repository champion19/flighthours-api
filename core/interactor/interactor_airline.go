package interactor

import (
	"context"

	"github.com/champion19/flighthours-api/core/interactor/services/domain"
	"github.com/champion19/flighthours-api/core/ports/input"
	"github.com/champion19/flighthours-api/middleware"
	"github.com/champion19/flighthours-api/platform/logger"
)

// AirlineInteractor orchestrates airline operations
type AirlineInteractor struct {
	service input.AirlineService
	logger  logger.Logger
}

// NewAirlineInteractor creates a new airline interactor
func NewAirlineInteractor(service input.AirlineService, log logger.Logger) *AirlineInteractor {
	return &AirlineInteractor{
		service: service,
		logger:  log,
	}
}

// GetAirlineByID retrieves an airline by its ID
func (i *AirlineInteractor) GetAirlineByID(ctx context.Context, id string) (*domain.Airline, error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := i.logger.WithTraceID(traceID)

	log.Info(logger.LogAirlineGet, "airline_id", id)

	airline, err := i.service.GetAirlineByID(ctx, id)
	if err != nil {
		log.Error(logger.LogAirlineGetError, "airline_id", id, "error", err)
		return nil, err
	}

	log.Success(logger.LogAirlineGetOK, airline.ToLogger())
	return airline, nil
}

// ActivateAirline sets an airline's status to active
func (i *AirlineInteractor) ActivateAirline(ctx context.Context, id string) error {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := i.logger.WithTraceID(traceID)

	log.Info(logger.LogAirlineActivate, "airline_id", id)

	// Verify airline exists
	_, err := i.service.GetAirlineByID(ctx, id)
	if err != nil {
		log.Error(logger.LogAirlineNotFound, "airline_id", id)
		return err
	}

	if err = i.service.ActivateAirline(ctx, id); err != nil {
		log.Error(logger.LogAirlineActivateError, "airline_id", id, "error", err)
		return err
	}

	log.Success(logger.LogAirlineActivateOK, "airline_id", id)
	return nil
}

// DeactivateAirline sets an airline's status to inactive
func (i *AirlineInteractor) DeactivateAirline(ctx context.Context, id string) error {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := i.logger.WithTraceID(traceID)

	log.Info(logger.LogAirlineDeactivate, "airline_id", id)

	// Verify airline exists
	_, err := i.service.GetAirlineByID(ctx, id)
	if err != nil {
		log.Error(logger.LogAirlineNotFound, "airline_id", id)
		return err
	}

	if err = i.service.DeactivateAirline(ctx, id); err != nil {
		log.Error(logger.LogAirlineDeactivateError, "airline_id", id, "error", err)
		return err
	}

	log.Success(logger.LogAirlineDeactivateOK, "airline_id", id)
	return nil
}

// ListAirlines retrieves all airlines with optional filters
func (i *AirlineInteractor) ListAirlines(ctx context.Context, filters map[string]interface{}) ([]domain.Airline, error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := i.logger.WithTraceID(traceID)

	log.Info(logger.LogAirlineList, "filters", filters)

	airlines, err := i.service.ListAirlines(ctx, filters)
	if err != nil {
		log.Error(logger.LogAirlineListError, "error", err)
		return nil, err
	}

	log.Success(logger.LogAirlineListOK, "count", len(airlines))
	return airlines, nil
}
