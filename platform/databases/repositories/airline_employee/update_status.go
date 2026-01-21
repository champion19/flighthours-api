package airline_employee

import (
	"context"

	domain "github.com/champion19/flighthours-api/core/interactor/services/domain"
	"github.com/champion19/flighthours-api/core/ports/output"
	"github.com/champion19/flighthours-api/platform/databases/common"
	"github.com/champion19/flighthours-api/platform/logger"
)

// UpdateAirlineEmployeeStatus updates the active status of an airline employee
func (r *repository) UpdateAirlineEmployeeStatus(ctx context.Context, tx output.Tx, id string, status bool) error {
	// Cast the transaction to the concrete type
	dbTx, ok := tx.(*common.SQLTX)
	if !ok {
		log.Error(logger.LogDatabaseUnavailable, "error", "invalid transaction type")
		return domain.ErrInvalidTransaction
	}

	result, err := dbTx.ExecContext(ctx, QueryUpdateStatus, status, id)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, "employee_id", id, "error", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return domain.ErrAirlineEmployeeNotFound
	}

	return nil
}
