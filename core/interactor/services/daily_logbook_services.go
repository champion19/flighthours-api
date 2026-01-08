package services

import (
	"context"

	"github.com/champion19/flighthours-api/core/interactor/services/domain"
	"github.com/champion19/flighthours-api/core/ports/output"
	"github.com/champion19/flighthours-api/platform/logger"
)

// DailyLogbookService implements the business logic for daily logbook operations
type DailyLogbookService struct {
	repo   output.DailyLogbookRepository
	logger logger.Logger
}

// NewDailyLogbookService creates a new daily logbook service
func NewDailyLogbookService(repo output.DailyLogbookRepository, log logger.Logger) *DailyLogbookService {
	return &DailyLogbookService{
		repo:   repo,
		logger: log,
	}
}

// BeginTx starts a new database transaction
func (s *DailyLogbookService) BeginTx(ctx context.Context) (output.Tx, error) {
	return s.repo.BeginTx(ctx)
}

// GetDailyLogbookByID retrieves a daily logbook by its ID
func (s *DailyLogbookService) GetDailyLogbookByID(ctx context.Context, id string) (*domain.DailyLogbook, error) {
	return s.repo.GetDailyLogbookByID(ctx, id)
}

// ListDailyLogbooksByEmployee retrieves all daily logbooks for a specific employee
func (s *DailyLogbookService) ListDailyLogbooksByEmployee(ctx context.Context, employeeID string, filters map[string]interface{}) ([]domain.DailyLogbook, error) {
	return s.repo.ListDailyLogbooksByEmployee(ctx, employeeID, filters)
}

// CreateDailyLogbook creates a new daily logbook entry with transaction handling
func (s *DailyLogbookService) CreateDailyLogbook(ctx context.Context, logbook domain.DailyLogbook) error {
	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Generate ID if not set
	if logbook.ID == "" {
		logbook.SetID()
	}

	if err = s.repo.SaveDailyLogbook(ctx, tx, logbook); err != nil {
		return err
	}

	return tx.Commit()
}

// UpdateDailyLogbook updates an existing daily logbook with transaction handling
func (s *DailyLogbookService) UpdateDailyLogbook(ctx context.Context, logbook domain.DailyLogbook) error {
	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	if err = s.repo.UpdateDailyLogbook(ctx, tx, logbook); err != nil {
		return err
	}

	return tx.Commit()
}

// DeleteDailyLogbook deletes a daily logbook with transaction handling
func (s *DailyLogbookService) DeleteDailyLogbook(ctx context.Context, id string) error {
	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	if err = s.repo.DeleteDailyLogbook(ctx, tx, id); err != nil {
		return err
	}

	return tx.Commit()
}

// UpdateDailyLogbookStatus updates the status of a daily logbook with transaction handling
func (s *DailyLogbookService) updateDailyLogbookStatus(ctx context.Context, id string, status bool) error {
	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	if err = s.repo.UpdateDailyLogbookStatus(ctx, tx, id, status); err != nil {
		return err
	}

	return tx.Commit()
}

// ActivateDailyLogbook sets the daily logbook status to true (active)
func (s *DailyLogbookService) ActivateDailyLogbook(ctx context.Context, id string) error {
	return s.updateDailyLogbookStatus(ctx, id, true)
}

// DeactivateDailyLogbook sets the daily logbook status to false (inactive)
func (s *DailyLogbookService) DeactivateDailyLogbook(ctx context.Context, id string) error {
	return s.updateDailyLogbookStatus(ctx, id, false)
}
