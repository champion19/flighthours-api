package aircraft_model

import (
	"context"
	"database/sql"

	"github.com/champion19/flighthours-api/core/interactor/services/domain"
	"github.com/champion19/flighthours-api/platform/logger"
)

const (
	QueryByID            = "SELECT id, model_name, aircraft_type_name, engine_type_name, family, manufacturer FROM aircraft_model WHERE id = ? LIMIT 1"
	QueryGetAll          = "SELECT id, model_name, aircraft_type_name, engine_type_name, family, manufacturer FROM aircraft_model ORDER BY model_name"
	QueryGetByEngineType = "SELECT id, model_name, aircraft_type_name, engine_type_name, family, manufacturer FROM aircraft_model WHERE engine_type_name = ? ORDER BY model_name"
)

var log logger.Logger = logger.NewSlogLogger()

type repository struct {
	stmtGetByID         *sql.Stmt
	stmtGetAll          *sql.Stmt
	stmtGetByEngineType *sql.Stmt
	db                  *sql.DB
}

// NewAircraftModelRepository creates a new aircraft model repository with prepared statements
func NewAircraftModelRepository(db *sql.DB) (*repository, error) {
	if db == nil {
		return nil, sql.ErrConnDone
	}

	stmtGetByID, err := db.Prepare(QueryByID)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, "error preparing statement", err)
		return nil, err
	}

	stmtGetAll, err := db.Prepare(QueryGetAll)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, "error preparing statement", err)
		return nil, err
	}

	stmtGetByEngineType, err := db.Prepare(QueryGetByEngineType)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, "error preparing statement", err)
		return nil, err
	}

	return &repository{
		db:                  db,
		stmtGetByID:         stmtGetByID,
		stmtGetAll:          stmtGetAll,
		stmtGetByEngineType: stmtGetByEngineType,
	}, nil
}

// GetAircraftModelByID retrieves an aircraft model by ID
func (r *repository) GetAircraftModelByID(ctx context.Context, id string) (*domain.AircraftModel, error) {
	var model domain.AircraftModel
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

	return &model, nil
}

// ListAircraftModels retrieves all aircraft models with optional filters
func (r *repository) ListAircraftModels(ctx context.Context, filters map[string]interface{}) ([]domain.AircraftModel, error) {
	var rows *sql.Rows
	var err error

	// Check if filtering by engine type
	if engineType, ok := filters["engine_type"]; ok {
		rows, err = r.stmtGetByEngineType.QueryContext(ctx, engineType)
	} else {
		rows, err = r.stmtGetAll.QueryContext(ctx)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var models []domain.AircraftModel
	for rows.Next() {
		var model domain.AircraftModel
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

		models = append(models, model)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return models, nil
}
