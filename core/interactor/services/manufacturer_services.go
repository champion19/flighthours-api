package services

import (
	"context"

	"github.com/champion19/flighthours-api/core/interactor/services/domain"
	"github.com/champion19/flighthours-api/core/ports/output"
	"github.com/champion19/flighthours-api/platform/logger"
)

// ManufacturerService implements the business logic for manufacturer operations
type ManufacturerService struct {
	repo   output.ManufacturerRepository
	logger logger.Logger
}

// NewManufacturerService creates a new manufacturer service
func NewManufacturerService(repo output.ManufacturerRepository, log logger.Logger) *ManufacturerService {
	return &ManufacturerService{
		repo:   repo,
		logger: log,
	}
}

// GetManufacturerByID retrieves a manufacturer by its ID
func (s *ManufacturerService) GetManufacturerByID(ctx context.Context, id string) (*domain.Manufacturer, error) {
	return s.repo.GetManufacturerByID(ctx, id)
}

// ListManufacturers retrieves all manufacturers
func (s *ManufacturerService) ListManufacturers(ctx context.Context) ([]domain.Manufacturer, error) {
	return s.repo.ListManufacturers(ctx)
}
