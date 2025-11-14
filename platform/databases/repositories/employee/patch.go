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
var dbTx *sqlTx
	var shouldCommit bool

	if tx != nil {
		// Usar la transacción existente
		dbTx = tx.(*sqlTx)
		shouldCommit = false
	} else {
		// Crear nueva transacción
		newTx, err := r.db.BeginTx(ctx, nil)
		if err != nil {
			return err
		}
		dbTx = &sqlTx{Tx: newTx}
		shouldCommit = true
	}


	result, err := dbTx.ExecContext(ctx, QueryPatch,keycloakUserID,id)
	if err != nil {
		dbTx.Rollback()
		return domain.ErrUserCannotSave
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		dbTx.Rollback()
		return 	err
	}
	if rowsAffected == 0 {
		dbTx.Rollback()
		return 	err
	}

	if shouldCommit {
		err=dbTx.Commit()
		if err != nil {
			return err
		}
	}
	return nil
}
