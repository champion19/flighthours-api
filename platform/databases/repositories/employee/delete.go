package employee

import (
	"context"
	"database/sql"

	"github.com/champion19/flighthours-api/core/interactor/services/domain"
	"github.com/champion19/flighthours-api/core/ports/output"
)

func (r *repository) DeleteEmployee(ctx context.Context, tx output.Tx, id string) error {
	dbtx, ok := tx.(*sql.Tx)
	if !ok {
		return domain.ErrInvalidRequest
	}

	_, err := dbtx.ExecContext(context.Background(), QueryDelete, id)
	if err != nil {
		return domain.ErrUserCannotSave
	}

	return nil
}
