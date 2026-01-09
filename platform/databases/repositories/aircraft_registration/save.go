package aircraft_registration

import (
	"context"
	"strings"

	domain "github.com/champion19/flighthours-api/core/interactor/services/domain"
	"github.com/champion19/flighthours-api/core/ports/output"
	"github.com/champion19/flighthours-api/platform/databases/common"
)

// SaveAircraftRegistration saves a new aircraft registration to the database
func (r *repository) SaveAircraftRegistration(ctx context.Context, tx output.Tx, registration domain.AircraftRegistration) error {
	sqlTx := tx.(*common.SQLTX)

	_, err := sqlTx.ExecContext(ctx, QueryInsert,
		registration.ID,
		registration.LicensePlate,
		registration.AircraftModelID,
		registration.AirlineID,
	)
	if err != nil {
		log.Error("SaveAircraftRegistration failed",
			"id", registration.ID,
			"license_plate", registration.LicensePlate,
			"aircraft_model_id", registration.AircraftModelID,
			"airline_id", registration.AirlineID,
			"error", err.Error())

		// Check for specific MySQL errors
		if strings.Contains(err.Error(), "Duplicate entry") {
			return domain.ErrAircraftRegistrationDuplicatePlate
		}
		if strings.Contains(err.Error(), "foreign key constraint") || strings.Contains(err.Error(), "FOREIGN KEY") {
			// Check which foreign key failed
			if strings.Contains(err.Error(), "aircraft_model") {
				return domain.ErrAircraftRegistrationInvalidModel
			}
			if strings.Contains(err.Error(), "airline") {
				return domain.ErrAircraftRegistrationInvalidAirline
			}
		}
		return domain.ErrAircraftRegistrationCannotSave
	}
	return nil
}
