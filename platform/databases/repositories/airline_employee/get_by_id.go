package airline_employee

import (
	"context"
	"database/sql"

	domain "github.com/champion19/flighthours-api/core/interactor/services/domain"
)

// GetAirlineEmployeeByID retrieves an airline employee by ID with denormalized airline data
func (r *repository) GetAirlineEmployeeByID(ctx context.Context, id string) (*domain.AirlineEmployee, error) {
	var e AirlineEmployee

	err := r.stmtGetByID.QueryRowContext(ctx, id).Scan(
		&e.ID,
		&e.Name,
		&e.AirlineID,
		&e.AirlineName,
		&e.AirlineCode,
		&e.Email,
		&e.IdentificationNumber,
		&e.Bp,
		&e.StartDate,
		&e.EndDate,
		&e.Active,
		&e.Role,
		&e.KeycloakUserID,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrAirlineEmployeeNotFound
		}
		return nil, err
	}

	return e.ToDomain(), nil
}
