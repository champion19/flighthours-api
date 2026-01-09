package aircraft_registration

import (
	"context"
	"database/sql"

	domain "github.com/champion19/flighthours-api/core/interactor/services/domain"
)

// GetAircraftRegistrationByID retrieves an aircraft registration by its UUID
func (r *repository) GetAircraftRegistrationByID(ctx context.Context, id string) (*domain.AircraftRegistration, error) {
	var ar AircraftRegistration
	err := r.stmtGetByID.QueryRowContext(ctx, id).Scan(
		&ar.ID,
		&ar.LicensePlate,
		&ar.AircraftModelID,
		&ar.AirlineID,
		&ar.ModelName,
		&ar.AirlineName,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrAircraftRegistrationNotFound
		}
		return nil, err
	}
	return ar.ToDomain(), nil
}
