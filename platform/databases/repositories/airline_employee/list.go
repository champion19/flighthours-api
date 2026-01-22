package airline_employee

import (
	"context"
	"database/sql"

	domain "github.com/champion19/flighthours-api/core/interactor/services/domain"
)

// ListAirlineEmployees retrieves all airline employees with optional filters
func (r *repository) ListAirlineEmployees(ctx context.Context, filters map[string]interface{}) ([]domain.AirlineEmployee, error) {
	var rows *sql.Rows
	var err error

	// Check filters
	if airlineID, ok := filters["airline_id"]; ok {
		rows, err = r.stmtGetByAirlineID.QueryContext(ctx, airlineID)
	} else if status, ok := filters["active"]; ok {
		rows, err = r.stmtGetByStatus.QueryContext(ctx, status)
	} else {
		rows, err = r.stmtGetAll.QueryContext(ctx)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var employees []domain.AirlineEmployee
	for rows.Next() {
		var e AirlineEmployee

		if err := rows.Scan(
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
		); err != nil {
			return nil, err
		}

		employees = append(employees, *e.ToDomain())
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return employees, nil
}
