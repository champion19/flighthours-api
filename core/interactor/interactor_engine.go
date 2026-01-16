package interactor

import (
	"context"

	"github.com/champion19/flighthours-api/core/interactor/services/domain"
	"github.com/champion19/flighthours-api/core/ports/input"
	"github.com/champion19/flighthours-api/middleware"
	"github.com/champion19/flighthours-api/platform/logger"
)

// EngineInteractor orchestrates engine operations
type EngineInteractor struct {
	service input.EngineService
	logger  logger.Logger
}

// NewEngineInteractor creates a new engine interactor
func NewEngineInteractor(service input.EngineService, log logger.Logger) *EngineInteractor {
	return &EngineInteractor{
		service: service,
		logger:  log,
	}
}

// GetEngineByID retrieves an engine by its ID
func (i *EngineInteractor) GetEngineByID(ctx context.Context, id string) (*domain.Engine, error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := i.logger.WithTraceID(traceID)

	log.Info(logger.LogEngineGet, "engine_id", id)

	engine, err := i.service.GetEngineByID(ctx, id)
	if err != nil {
		log.Error(logger.LogEngineGetError, "engine_id", id, "error", err)
		return nil, err
	}

	log.Success(logger.LogEngineGetOK, engine.ToLogger())
	return engine, nil
}

// ListEngines retrieves all engines
func (i *EngineInteractor) ListEngines(ctx context.Context) ([]domain.Engine, error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := i.logger.WithTraceID(traceID)

	log.Info(logger.LogEngineList)

	engines, err := i.service.ListEngines(ctx)
	if err != nil {
		log.Error(logger.LogEngineListError, "error", err)
		return nil, err
	}

	log.Success(logger.LogEngineListOK, "count", len(engines))
	return engines, nil
}
