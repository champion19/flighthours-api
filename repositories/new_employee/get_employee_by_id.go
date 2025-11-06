package employee

import(

	"context"
	"database/sql"

	"github.com/champion19/flighthours-api/core/domain"
)


func (r *repository) GetEmployeeByID(ctx context.Context,id string) (*domain.Employee, error) {

	tx,err:=r.db.BeginTx(ctx, nil)
	if err!=nil{
		return nil,domain.ErrUserCannotSave
	}
	row := tx.QueryRowContext(context.Background(),QueryByID,id)

	var e Employee
	err = row.Scan(
		&e.ID,
		&e.Name,
		&e.Airline,
		&e.Email,
		&e.IdentificationNumber,
		&e.Bp, &e.StartDate,
		&e.EndDate,
		&e.Active,
		&e.Role,
		&e.KeycloakUserID)

	if err != nil {
		tx.Rollback()
		if err == sql.ErrNoRows {
			return nil, domain.ErrPersonNotFound
		}
		return nil, domain.ErrUserCannotSave
	}



if e.KeycloakUserID == "" {
		keycloakUser, keycloakErr := r.keycloak.GetUserByEmail(ctx, e.Email)
		if keycloakErr != nil {
			tx.Rollback()
			return nil, domain.ErrUserCannotGet
		}
		if keycloakUser != nil && keycloakUser.ID != nil {
			e.KeycloakUserID = *keycloakUser.ID

			// Podrías opcionalmente actualizar el registro con el nuevo ID
			updateErr := r.UpdateEmployee(ctx, e.ToDomain())
			if updateErr != nil {
				tx.Rollback()
				return nil, domain.ErrUserCannotSave
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	domainEmployee := e.ToDomain()
	return &domainEmployee,nil

}

