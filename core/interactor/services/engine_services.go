package services

import (
	"context"

	"github.com/champion19/flighthours-api/core/interactor/services/domain"
	"github.com/champion19/flighthours-api/core/ports/output"
	"github.com/champion19/flighthours-api/platform/logger"
)

// EngineService implements the business logic for engine operations
type EngineService struct {
	repo   output.EngineRepository
	logger logger.Logger
}

// NewEngineService creates a new engine service
func NewEngineService(repo output.EngineRepository, log logger.Logger) *EngineService {
	return &EngineService{
		repo:   repo,
		logger: log,
	}
}

// GetEngineByID retrieves an engine by its ID
func (s *EngineService) GetEngineByID(ctx context.Context, id string) (*domain.Engine, error) {
	return s.repo.GetEngineByID(ctx, id)
}

// ListEngines retrieves all engines
func (s *EngineService) ListEngines(ctx context.Context) ([]domain.Engine, error) {
	return s.repo.ListEngines(ctx)
}
