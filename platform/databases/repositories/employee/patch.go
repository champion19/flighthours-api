package employee

import (
	"context"


	"github.com/champion19/flighthours-api/core/interactor/services/domain"
	"github.com/champion19/flighthours-api/core/ports/output"
)

const (
	QueryPatch = "UPDATE employee SET keycloak_user_id=? WHERE id=?"
)

func (r *repository) PatchEmployee( ctx context.Context,tx output.Tx,id string, keycloakUserID string) error {
  dbTx, ok := tx.(*sqlTx)

	if !ok {
		return domain.ErrInvalidRequest
	}


	result, err := dbTx.ExecContext(ctx, QueryPatch,keycloakUserID,id)
	if err != nil {
		return domain.ErrUserCannotSave
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 	err
	}
	if rowsAffected == 0 {
		return 	domain.ErrUserCannotSave
	}
	return nil
}
