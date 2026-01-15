package daily_logbook_detail

import (
	"context"

	"github.com/champion19/flighthours-api/core/interactor/services/domain"
	"github.com/champion19/flighthours-api/core/ports/output"
	"github.com/champion19/flighthours-api/platform/databases/common"
	"github.com/champion19/flighthours-api/platform/logger"
)

// DeleteDailyLogbookDetail deletes a daily logbook detail by its ID
func (r *repository) DeleteDailyLogbookDetail(ctx context.Context, tx output.Tx, id string) error {
	log.Info(logger.LogDailyLogbookDetailDelete, "id", id)

	sqlTx, ok := tx.(*common.SQLTX)
	if !ok {
		log.Error(logger.LogDailyLogbookDetailDeleteError, "error", "invalid transaction type")
		return domain.ErrInvalidTransaction
	}

	stmt := sqlTx.Tx.StmtContext(ctx, r.stmtDelete)

	result, err := stmt.ExecContext(ctx, id)
	if err != nil {
		log.Error(logger.LogDailyLogbookDetailDeleteError, "id", id, "error", err)
		return domain.ErrFlightCannotDelete
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Error(logger.LogDailyLogbookDetailDeleteError, "id", id, "error", err)
		return domain.ErrFlightCannotDelete
	}

	if rowsAffected == 0 {
		log.Warn(logger.LogDailyLogbookDetailNotFound, "id", id)
		return domain.ErrFlightNotFound
	}

	log.Info(logger.LogDailyLogbookDetailDeleteOK, "id", id)
	return nil
}
