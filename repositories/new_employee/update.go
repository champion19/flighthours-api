package newemployee

import (
	"context"

	"github.com/champion19/flighthours-api/core/domain"
)



func (r *repository) UpdateEmployee(ctx context.Context,employee domain.Employee) error {
	employeeToUpdate := FromDomain(employee)

	tx, err := r.db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx,QueryUpdate,

		employeeToUpdate.Name,
		employeeToUpdate.Airline,
		employeeToUpdate.Email,
		employeeToUpdate.IdentificationNumber,
		employeeToUpdate.Bp,
		employeeToUpdate.StartDate,
		employeeToUpdate.EndDate,
		employeeToUpdate.Active,
		employeeToUpdate.Role,
	  employeeToUpdate.KeycloakUserID,
	  employeeToUpdate.ID,
	)

	if err != nil {
		tx.Rollback()
		return domain.ErrUserCannotSave
	}

	return nil
}
