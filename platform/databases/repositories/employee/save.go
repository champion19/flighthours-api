package employee


import (
	"context"

	"github.com/champion19/flighthours-api/core/interactor/services/domain"
	"github.com/go-sql-driver/mysql"
)

func (r *repository) Save( employee domain.Employee) error {

	employeeToSave := FromDomain(employee)

	//begin transaction
	tx, err := r.db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(context.Background(),QuerySave,
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
		tx.Rollback()
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
			return domain.ErrDuplicateUser
		} else {
			return domain.ErrUserCannotSave
		}
	}

	if err = tx.Commit(); err != nil {
		return domain.ErrUserCannotSave
		}

	return nil
}
