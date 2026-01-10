package interactor

import (
	"context"

	"github.com/champion19/flighthours-api/core/interactor/services/domain"
	"github.com/champion19/flighthours-api/core/ports/input"
	"github.com/champion19/flighthours-api/middleware"
	"github.com/champion19/flighthours-api/platform/logger"
)

// AircraftModelInteractor orchestrates aircraft model operations
type AircraftModelInteractor struct {
	service input.AircraftModelService
	logger  logger.Logger
}

// NewAircraftModelInteractor creates a new aircraft model interactor
func NewAircraftModelInteractor(service input.AircraftModelService, log logger.Logger) *AircraftModelInteractor {
	return &AircraftModelInteractor{
		service: service,
		logger:  log,
	}
}

// GetAircraftModelByID retrieves an aircraft model by its ID
func (i *AircraftModelInteractor) GetAircraftModelByID(ctx context.Context, id string) (*domain.AircraftModel, error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := i.logger.WithTraceID(traceID)

	log.Info(logger.LogAircraftModelGet, "aircraft_model_id", id)

	model, err := i.service.GetAircraftModelByID(ctx, id)
	if err != nil {
		log.Error(logger.LogAircraftModelGetError, "aircraft_model_id", id, "error", err)
		return nil, err
	}

	log.Success(logger.LogAircraftModelGetOK, model.ToLogger())
	return model, nil
}

// ListAircraftModels retrieves all aircraft models with optional filters
func (i *AircraftModelInteractor) ListAircraftModels(ctx context.Context, filters map[string]interface{}) ([]domain.AircraftModel, error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := i.logger.WithTraceID(traceID)

	log.Info(logger.LogAircraftModelList, "filters", filters)

	models, err := i.service.ListAircraftModels(ctx, filters)
	if err != nil {
		log.Error(logger.LogAircraftModelListError, "error", err)
		return nil, err
	}

	log.Success(logger.LogAircraftModelListOK, "count", len(models))
	return models, nil
}
