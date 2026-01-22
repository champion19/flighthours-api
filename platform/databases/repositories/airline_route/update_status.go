package airline_route

import (
	"context"
	"database/sql"

	domain "github.com/champion19/flighthours-api/core/interactor/services/domain"
	"github.com/champion19/flighthours-api/core/ports/output"
)

// UpdateAirlineRouteStatus updates the status of an airline route
func (r *repository) UpdateAirlineRouteStatus(ctx context.Context, tx output.Tx, id string, status bool) error {
	sqlTx, ok := tx.(*sql.Tx)
	if !ok {
		return domain.ErrInvalidTransaction
	}

	stmt := sqlTx.StmtContext(ctx, r.stmtUpdateStatus)
	result, err := stmt.ExecContext(ctx, status, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return domain.ErrAirlineRouteNotFound
	}

	return nil
}
