package services

import (
	"context"

	"github.com/champion19/flighthours-api/core/interactor/services/domain"
	"github.com/champion19/flighthours-api/core/ports/output"
	"github.com/champion19/flighthours-api/platform/logger"
)

// AircraftRegistrationService implements the business logic for aircraft registration operations
type AircraftRegistrationService struct {
	repo   output.AircraftRegistrationRepository
	logger logger.Logger
}

// NewAircraftRegistrationService creates a new aircraft registration service
func NewAircraftRegistrationService(repo output.AircraftRegistrationRepository, log logger.Logger) *AircraftRegistrationService {
	return &AircraftRegistrationService{
		repo:   repo,
		logger: log,
	}
}

// BeginTx starts a new database transaction
func (s *AircraftRegistrationService) BeginTx(ctx context.Context) (output.Tx, error) {
	return s.repo.BeginTx(ctx)
}

// GetAircraftRegistrationByID retrieves an aircraft registration by its ID
func (s *AircraftRegistrationService) GetAircraftRegistrationByID(ctx context.Context, id string) (*domain.AircraftRegistration, error) {
	return s.repo.GetAircraftRegistrationByID(ctx, id)
}

// ListAircraftRegistrations retrieves all aircraft registrations with optional filters
func (s *AircraftRegistrationService) ListAircraftRegistrations(ctx context.Context, filters map[string]interface{}) ([]domain.AircraftRegistration, error) {
	return s.repo.ListAircraftRegistrations(ctx, filters)
}

// CreateAircraftRegistration creates a new aircraft registration with transaction handling
func (s *AircraftRegistrationService) CreateAircraftRegistration(ctx context.Context, registration domain.AircraftRegistration) error {
	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	if err = s.repo.SaveAircraftRegistration(ctx, tx, registration); err != nil {
		return err
	}

	return tx.Commit()
}

// UpdateAircraftRegistration updates an existing aircraft registration with transaction handling
func (s *AircraftRegistrationService) UpdateAircraftRegistration(ctx context.Context, registration domain.AircraftRegistration) error {
	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	if err = s.repo.UpdateAircraftRegistration(ctx, tx, registration); err != nil {
		return err
	}

	return tx.Commit()
}
