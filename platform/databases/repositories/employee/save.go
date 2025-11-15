package employee


import (
	"context"

	"github.com/champion19/flighthours-api/core/interactor/services/domain"
	"github.com/champion19/flighthours-api/core/ports/output"
	"github.com/go-sql-driver/mysql"
)

func (r *repository) Save(ctx context.Context, tx output.Tx, 	employee domain.Employee) error {

	employeeToSave := FromDomain(employee)

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

	_, err := dbTx.ExecContext(ctx,QuerySave,
		employeeToSave.ID,
		employeeToSave.Name,
		employeeToSave.Airline,
		employeeToSave.Email,
		employeeToSave.IdentificationNumber,
		employeeToSave.Bp,
		employeeToSave.StartDate,
		employeeToSave.EndDate,
		employeeToSave.Active,
		employeeToSave.Role,
		employeeToSave.KeycloakUserID)

	if err != nil {
		if shouldCommit {
			dbTx.Rollback()
		}
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
			return domain.ErrDuplicateUser
		} else {
			return domain.ErrUserCannotSave
		}
	}

	if shouldCommit {
		if err = dbTx.Commit(); err != nil {
			return domain.ErrUserCannotSave
		}
	}

	return nil
}
