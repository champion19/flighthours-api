package aircraft_registration

import (
	"context"
	"database/sql"

	domain "github.com/champion19/flighthours-api/core/interactor/services/domain"
)

// ListAircraftRegistrations retrieves all aircraft registrations, optionally filtered by airline_id
func (r *repository) ListAircraftRegistrations(ctx context.Context, filters map[string]interface{}) ([]domain.AircraftRegistration, error) {
	var rows *sql.Rows
	var err error

	// Check if airline_id filter is provided
	if airlineID, ok := filters["airline_id"]; ok {
		airlineIDStr, isString := airlineID.(string)
		if isString && airlineIDStr != "" {
			rows, err = r.stmtGetByAirline.QueryContext(ctx, airlineIDStr)
		} else {
			rows, err = r.stmtGetAll.QueryContext(ctx)
		}
	} else {
		rows, err = r.stmtGetAll.QueryContext(ctx)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var registrations []domain.AircraftRegistration
	for rows.Next() {
		var ar AircraftRegistration
		if err := rows.Scan(
			&ar.ID,
			&ar.LicensePlate,
			&ar.AircraftModelID,
			&ar.AirlineID,
			&ar.ModelName,
			&ar.AirlineName,
		); err != nil {
			return nil, err
		}
		registrations = append(registrations, *ar.ToDomain())
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return registrations, nil
}
