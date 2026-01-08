package interactor

import (
	"context"

	"github.com/champion19/flighthours-api/core/interactor/services/domain"
	"github.com/champion19/flighthours-api/core/ports/input"
	"github.com/champion19/flighthours-api/middleware"
	"github.com/champion19/flighthours-api/platform/logger"
)

// DailyLogbookInteractor orchestrates daily logbook operations
type DailyLogbookInteractor struct {
	service input.DailyLogbookService
	logger  logger.Logger
}

// NewDailyLogbookInteractor creates a new daily logbook interactor
func NewDailyLogbookInteractor(service input.DailyLogbookService, log logger.Logger) *DailyLogbookInteractor {
	return &DailyLogbookInteractor{
		service: service,
		logger:  log,
	}
}

// GetDailyLogbookByID retrieves a daily logbook by its ID
func (i *DailyLogbookInteractor) GetDailyLogbookByID(ctx context.Context, id string) (*domain.DailyLogbook, error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := i.logger.WithTraceID(traceID)

	log.Info(logger.LogDailyLogbookGet, "logbook_id", id)

	logbook, err := i.service.GetDailyLogbookByID(ctx, id)
	if err != nil {
		log.Error(logger.LogDailyLogbookGetError, "logbook_id", id, "error", err)
		return nil, err
	}

	log.Success(logger.LogDailyLogbookGetOK, logbook.ToLogger())
	return logbook, nil
}

// ListDailyLogbooksByEmployee retrieves all daily logbooks for the authenticated employee
func (i *DailyLogbookInteractor) ListDailyLogbooksByEmployee(ctx context.Context, employeeID string, filters map[string]interface{}) ([]domain.DailyLogbook, error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := i.logger.WithTraceID(traceID)

	log.Info(logger.LogDailyLogbookList, "employee_id", employeeID, "filters", filters)

	logbooks, err := i.service.ListDailyLogbooksByEmployee(ctx, employeeID, filters)
	if err != nil {
		log.Error(logger.LogDailyLogbookListError, "error", err)
		return nil, err
	}

	log.Success(logger.LogDailyLogbookListOK, "count", len(logbooks))
	return logbooks, nil
}

// CreateDailyLogbook creates a new daily logbook for the authenticated employee
func (i *DailyLogbookInteractor) CreateDailyLogbook(ctx context.Context, logbook domain.DailyLogbook) error {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := i.logger.WithTraceID(traceID)

	log.Info(logger.LogDailyLogbookCreate, "employee_id", logbook.EmployeeID)

	if err := i.service.CreateDailyLogbook(ctx, logbook); err != nil {
		log.Error(logger.LogDailyLogbookCreateError, "error", err)
		return err
	}

	log.Success(logger.LogDailyLogbookCreateOK, logbook.ToLogger())
	return nil
}

// UpdateDailyLogbook updates an existing daily logbook
func (i *DailyLogbookInteractor) UpdateDailyLogbook(ctx context.Context, logbook domain.DailyLogbook) error {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := i.logger.WithTraceID(traceID)

	log.Info(logger.LogDailyLogbookUpdate, "logbook_id", logbook.ID)

	// Verify logbook exists
	_, err := i.service.GetDailyLogbookByID(ctx, logbook.ID)
	if err != nil {
		log.Error(logger.LogDailyLogbookNotFound, "logbook_id", logbook.ID)
		return err
	}

	if err = i.service.UpdateDailyLogbook(ctx, logbook); err != nil {
		log.Error(logger.LogDailyLogbookUpdateError, "logbook_id", logbook.ID, "error", err)
		return err
	}

	log.Success(logger.LogDailyLogbookUpdateOK, logbook.ToLogger())
	return nil
}

// DeleteDailyLogbook deletes a daily logbook
func (i *DailyLogbookInteractor) DeleteDailyLogbook(ctx context.Context, id string) error {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := i.logger.WithTraceID(traceID)

	log.Info(logger.LogDailyLogbookDelete, "logbook_id", id)

	// Verify logbook exists
	_, err := i.service.GetDailyLogbookByID(ctx, id)
	if err != nil {
		log.Error(logger.LogDailyLogbookNotFound, "logbook_id", id)
		return err
	}

	if err = i.service.DeleteDailyLogbook(ctx, id); err != nil {
		log.Error(logger.LogDailyLogbookDeleteError, "logbook_id", id, "error", err)
		return err
	}

	log.Success(logger.LogDailyLogbookDeleteOK, "logbook_id", id)
	return nil
}

// ActivateDailyLogbook sets a daily logbook's status to active
func (i *DailyLogbookInteractor) ActivateDailyLogbook(ctx context.Context, id string) error {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := i.logger.WithTraceID(traceID)

	log.Info(logger.LogDailyLogbookActivate, "logbook_id", id)

	// Verify logbook exists
	_, err := i.service.GetDailyLogbookByID(ctx, id)
	if err != nil {
		log.Error(logger.LogDailyLogbookNotFound, "logbook_id", id)
		return err
	}

	if err = i.service.ActivateDailyLogbook(ctx, id); err != nil {
		log.Error(logger.LogDailyLogbookActivateError, "logbook_id", id, "error", err)
		return err
	}

	log.Success(logger.LogDailyLogbookActivateOK, "logbook_id", id)
	return nil
}

// DeactivateDailyLogbook sets a daily logbook's status to inactive
func (i *DailyLogbookInteractor) DeactivateDailyLogbook(ctx context.Context, id string) error {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := i.logger.WithTraceID(traceID)

	log.Info(logger.LogDailyLogbookDeactivate, "logbook_id", id)

	// Verify logbook exists
	_, err := i.service.GetDailyLogbookByID(ctx, id)
	if err != nil {
		log.Error(logger.LogDailyLogbookNotFound, "logbook_id", id)
		return err
	}

	if err = i.service.DeactivateDailyLogbook(ctx, id); err != nil {
		log.Error(logger.LogDailyLogbookDeactivateError, "logbook_id", id, "error", err)
		return err
	}

	log.Success(logger.LogDailyLogbookDeactivateOK, "logbook_id", id)
	return nil
}
