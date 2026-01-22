package aircraft_model

import (
	"context"
	"database/sql"

	domain "github.com/champion19/flighthours-api/core/interactor/services/domain"
)

// GetAircraftModelsByFamily retrieves all aircraft models for a specific family (HU32)
func (r *repository) GetAircraftModelsByFamily(ctx context.Context, family string) ([]domain.AircraftModel, error) {
	rows, err := r.stmtGetByFamily.QueryContext(ctx, family)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var models []domain.AircraftModel
	for rows.Next() {
		var model AircraftModel
		var engineTypeName sql.NullString
		var manufacturer sql.NullString

		if err := rows.Scan(
			&model.ID,
			&model.ModelName,
			&model.AircraftTypeName,
			&engineTypeName,
			&model.Family,
			&manufacturer,
		); err != nil {
			return nil, err
		}

		if engineTypeName.Valid {
			model.EngineTypeName = engineTypeName.String
		}
		if manufacturer.Valid {
			model.Manufacturer = manufacturer.String
		}

		models = append(models, *model.ToDomain())
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return models, nil
}
