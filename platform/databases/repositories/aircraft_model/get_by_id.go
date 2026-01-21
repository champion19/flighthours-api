package aircraft_model

import (
	"context"
	"database/sql"

	domain "github.com/champion19/flighthours-api/core/interactor/services/domain"
)

// GetAircraftModelByID retrieves an aircraft model by ID
func (r *repository) GetAircraftModelByID(ctx context.Context, id string) (*domain.AircraftModel, error) {
	var model AircraftModel
	var engineTypeName sql.NullString
	var manufacturer sql.NullString

	err := r.stmtGetByID.QueryRowContext(ctx, id).Scan(
		&model.ID,
		&model.ModelName,
		&model.AircraftTypeName,
		&engineTypeName,
		&model.Family,
		&manufacturer,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrAircraftModelNotFound
		}
		return nil, err
	}

	if engineTypeName.Valid {
		model.EngineTypeName = engineTypeName.String
	}
	if manufacturer.Valid {
		model.Manufacturer = manufacturer.String
	}

	return model.ToDomain(), nil
}
