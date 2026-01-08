package daily_logbook

import (
	"context"

	domain "github.com/champion19/flighthours-api/core/interactor/services/domain"
	"github.com/champion19/flighthours-api/core/ports/output"
	"github.com/champion19/flighthours-api/platform/databases/common"
)

// DeleteDailyLogbook removes a daily logbook entry
func (r *repository) DeleteDailyLogbook(ctx context.Context, tx output.Tx, id string) error {
	sqlTx := tx.(*common.SQLTX)

	result, err := sqlTx.ExecContext(ctx, QueryDelete, id)
	if err != nil {
		return domain.ErrDailyLogbookCannotDelete
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
