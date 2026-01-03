package services

import (
	"context"

	"github.com/champion19/flighthours-api/core/interactor/services/domain"
	"github.com/champion19/flighthours-api/core/ports/output"
	"github.com/champion19/flighthours-api/platform/logger"
)

// AirlineService implements the business logic for airline operations
type AirlineService struct {
	repo   output.AirlineRepository
	logger logger.Logger
}

// NewAirlineService creates a new airline service
func NewAirlineService(repo output.AirlineRepository, log logger.Logger) *AirlineService {
	return &AirlineService{
		repo:   repo,
		logger: log,
	}
}

// BeginTx starts a new database transaction
func (s *AirlineService) BeginTx(ctx context.Context) (output.Tx, error) {
	return s.repo.BeginTx(ctx)
}

// GetAirlineByID retrieves an airline by its ID
func (s *AirlineService) GetAirlineByID(ctx context.Context, id string) (*domain.Airline, error) {
	return s.repo.GetAirlineByID(ctx, id)
}

// UpdateAirlineStatus updates the status of an airline with transaction handling
func (s *AirlineService) UpdateAirlineStatus(ctx context.Context, id string, status bool) error {
	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	if err = s.repo.UpdateAirlineStatus(ctx, tx, id, status); err != nil {
		return err
	}

	return tx.Commit()
}

// ActivateAirline sets the airline status to true (active)
func (s *AirlineService) ActivateAirline(ctx context.Context, id string) error {
	return s.UpdateAirlineStatus(ctx, id, true)
}

// DeactivateAirline sets the airline status to false (inactive)
func (s *AirlineService) DeactivateAirline(ctx context.Context, id string) error {
	return s.UpdateAirlineStatus(ctx, id, false)
}
