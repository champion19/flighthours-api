package employee

import (
	"context"
	"database/sql"

	"github.com/champion19/flighthours-api/core/interactor/services/domain"
)


func (r *repository) GetEmployeeByEmail(ctx context.Context,email string) (*domain.Employee, error) {
	tx,err:=r.db.BeginTx(ctx,nil)
	if err!=nil{
		return nil,domain.ErrUserCannotSave
	}

	row := tx.QueryRowContext(ctx,QueryByEmail,email)

	var e Employee
	err = row.Scan(
		&e.ID,
		&e.Name,
		&e.Airline,
		&e.Email,
		&e.IdentificationNumber,
		&e.Bp,
		&e.StartDate,
		&e.EndDate,
		&e.Active,
		&e.Role,
		&e.KeycloakUserID)

   if err != nil {
		tx.Rollback()
		if err== sql.ErrNoRows{
			return nil,domain.ErrPersonNotFound
		}
		return nil,domain.ErrUserCannotSave
	}

	if err:=tx.Commit();err!=nil{
		return nil,domain.ErrUserCannotSave
	}

  domainEmployee := e.ToDomain()
	return &domainEmployee,nil

}
