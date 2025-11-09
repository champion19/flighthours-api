package employee

import (
	"context"


	"github.com/champion19/flighthours-api/core/interactor/services/domain"
)

const (
	QueryPatch = "UPDATE employee SET keycloak_user_id=? WHERE id=?"
)

func (r *repository) PatchEmployee( id string, keycloakUserID string) error {
	tx, err := r.db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}
	result, err := tx.ExecContext(context.Background(), QueryPatch,keycloakUserID,id)
	if err != nil {
		tx.Rollback()
		return domain.ErrUserCannotSave
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		tx.Rollback()
		return 	err
	}
	if rowsAffected == 0 {
		tx.Rollback()
		return 	err
	}

	err=tx.Commit()
	if err != nil {
		return err
	}
	return nil
}
