package employee

import (
	"context"

	"github.com/champion19/flighthours-api/core/interactor/services/domain"
	"github.com/champion19/flighthours-api/core/ports/output"
)



func (r *repository) UpdateEmployee(ctx context.Context,tx output.Tx,employee domain.Employee) error {
	employeeToUpdate := FromDomain(employee)

	dbTx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	_, err = dbTx.ExecContext(ctx, QueryUpdate,

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
			dbTx.Rollback()
		return domain.ErrUserCannotSave
	}

 err = dbTx.Commit()
 if err != nil {
	return err
}
return nil
}
