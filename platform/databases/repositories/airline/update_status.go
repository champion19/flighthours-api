package airline

import (
	"context"

	domain "github.com/champion19/flighthours-api/core/interactor/services/domain"
	"github.com/champion19/flighthours-api/core/ports/output"
	"github.com/champion19/flighthours-api/platform/databases/common"
)

// UpdateAirlineStatus updates the status of an airline (active/inactive)
func (r *repository) UpdateAirlineStatus(ctx context.Context, tx output.Tx, id string, status string) error {
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
		return domain.ErrAirlineNotFound
	}

	return nil
}
