package employee

import (
	"context"

	"github.com/champion19/flighthours-api/core/interactor/services/domain"
)



func (r *repository) UpdateEmployee(employee domain.Employee) error {
	employeeToUpdate := FromDomain(employee)

	tx, err := r.db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(context.Background(), QueryUpdate,

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

 err = tx.Commit()
 if err != nil {
	return err
}
return nil
}
