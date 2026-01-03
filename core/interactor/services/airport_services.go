package services

import (
	"context"

	"github.com/champion19/flighthours-api/core/interactor/services/domain"
	"github.com/champion19/flighthours-api/core/ports/output"
	"github.com/champion19/flighthours-api/platform/logger"
)

// AirportService implements the business logic for airport operations
type AirportService struct {
	repo   output.AirportRepository
	logger logger.Logger
}

// NewAirportService creates a new airport service
func NewAirportService(repo output.AirportRepository, log logger.Logger) *AirportService {
	return &AirportService{
		repo:   repo,
		logger: log,
	}
}

// BeginTx starts a new database transaction
func (s *AirportService) BeginTx(ctx context.Context) (output.Tx, error) {
	return s.repo.BeginTx(ctx)
}

// GetAirportByID retrieves an airport by its ID
func (s *AirportService) GetAirportByID(ctx context.Context, id string) (*domain.Airport, error) {
	return s.repo.GetAirportByID(ctx, id)
}

// UpdateAirportStatus updates the status of an airport with transaction handling
func (s *AirportService) UpdateAirportStatus(ctx context.Context, id string, status bool) error {
	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	if err = s.repo.UpdateAirportStatus(ctx, tx, id, status); err != nil {
		return err
	}

	return tx.Commit()
}

// ActivateAirport sets the airport status to true (active)
func (s *AirportService) ActivateAirport(ctx context.Context, id string) error {
	return s.UpdateAirportStatus(ctx, id, true)
}

// DeactivateAirport sets the airport status to false (inactive)
func (s *AirportService) DeactivateAirport(ctx context.Context, id string) error {
	return s.UpdateAirportStatus(ctx, id, false)
}
