package daily_logbook

import (
	"context"

	domain "github.com/champion19/flighthours-api/core/interactor/services/domain"
	"github.com/champion19/flighthours-api/core/ports/output"
	"github.com/champion19/flighthours-api/platform/databases/common"
)

// UpdateDailyLogbook updates an existing daily logbook entry
func (r *repository) UpdateDailyLogbook(ctx context.Context, tx output.Tx, logbook domain.DailyLogbook) error {
	sqlTx := tx.(*common.SQLTX)

	result, err := sqlTx.ExecContext(ctx, QueryUpdate,
		logbook.LogDate,
		logbook.BookPage,
		logbook.Status,
		logbook.ID,
	)
	if err != nil {
		return domain.ErrDailyLogbookCannotUpdate
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return domain.ErrDailyLogbookNotFound
	}

	return nil
}

// UpdateDailyLogbookStatus updates only the status of a daily logbook
func (r *repository) UpdateDailyLogbookStatus(ctx context.Context, tx output.Tx, id string, status bool) error {
	sqlTx := tx.(*common.SQLTX)

	result, err := sqlTx.ExecContext(ctx, QueryUpdateStatus, status, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return domain.ErrDailyLogbookNotFound
	}

	return nil
}
