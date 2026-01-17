package services

import (
	"context"

	"github.com/champion19/flighthours-api/core/interactor/services/domain"
	"github.com/champion19/flighthours-api/core/ports/output"
	"github.com/champion19/flighthours-api/platform/logger"
)

// AircraftModelService implements the business logic for aircraft model operations
type AircraftModelService struct {
	repo   output.AircraftModelRepository
	logger logger.Logger
}

// NewAircraftModelService creates a new aircraft model service
func NewAircraftModelService(repo output.AircraftModelRepository, log logger.Logger) *AircraftModelService {
	return &AircraftModelService{
		repo:   repo,
		logger: log,
	}
}

// GetAircraftModelByID retrieves an aircraft model by its ID
func (s *AircraftModelService) GetAircraftModelByID(ctx context.Context, id string) (*domain.AircraftModel, error) {
	return s.repo.GetAircraftModelByID(ctx, id)
}

// ListAircraftModels retrieves all aircraft models with optional filters
func (s *AircraftModelService) ListAircraftModels(ctx context.Context, filters map[string]interface{}) ([]domain.AircraftModel, error) {
	return s.repo.ListAircraftModels(ctx, filters)
}

// GetAircraftModelsByFamily retrieves all aircraft models for a specific family (HU32)
func (s *AircraftModelService) GetAircraftModelsByFamily(ctx context.Context, family string) ([]domain.AircraftModel, error) {
	return s.repo.GetAircraftModelsByFamily(ctx, family)
}
