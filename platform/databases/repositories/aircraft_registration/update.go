package aircraft_registration

import (
	"context"

	domain "github.com/champion19/flighthours-api/core/interactor/services/domain"
	"github.com/champion19/flighthours-api/core/ports/output"
	"github.com/champion19/flighthours-api/platform/databases/common"
)

// UpdateAircraftRegistration updates an existing aircraft registration in the database
func (r *repository) UpdateAircraftRegistration(ctx context.Context, tx output.Tx, registration domain.AircraftRegistration) error {
	sqlTx := tx.(*common.SQLTX)

	result, err := sqlTx.ExecContext(ctx, QueryUpdate,
		registration.LicensePlate,
		registration.AircraftModelID,
		registration.AirlineID,
		registration.ID,
	)
	if err != nil {
		return domain.ErrAircraftRegistrationCannotUpdate
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return domain.ErrAircraftRegistrationCannotUpdate
	}

	if rowsAffected == 0 {
		return domain.ErrAircraftRegistrationNotFound
	}

	return nil
}
