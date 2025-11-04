package newemployee

import (
	"context"

	"github.com/champion19/flighthours-api/core/domain"
	"github.com/go-sql-driver/mysql"
)

func (r *repository) SaveEmployee(ctx context.Context, employee domain.Employee) error {

	employeeToSave := FromDomain(employee)

	//begin transaction
	tx, err := r.db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, QuerySave,
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

	// TODO logica de keycloak

keycloakUserID, err := r.keycloak.CreateUser(ctx, &employee)
if err != nil {

		existingUser, getErr := r.keycloak.GetUserByEmail(ctx, employee.Email)
		if getErr != nil {
			return domain.ErrUserCannotSave
		}
		keycloakUserID = *existingUser.ID
	}
		employee.KeycloakUserID = keycloakUserID
	err = r.UpdateEmployee(ctx,employee)
	if err != nil {
		tx.Rollback()
		_ = r.keycloak.DeleteUser(ctx, keycloakUserID)
		return domain.ErrUserCannotSave
	}


	// commit transaction
	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}
