package manufacturer

import (
	"context"
	"database/sql"

	"github.com/champion19/flighthours-api/core/interactor/services/domain"
	"github.com/champion19/flighthours-api/platform/logger"
)

const (
	QueryByID   = "SELECT id, name FROM manufacturer WHERE id = ? LIMIT 1"
	QueryGetAll = "SELECT id, name FROM manufacturer ORDER BY name"
)

var log logger.Logger = logger.NewSlogLogger()

type repository struct {
	stmtGetByID *sql.Stmt
	stmtGetAll  *sql.Stmt
	db          *sql.DB
}

// NewManufacturerRepository creates a new manufacturer repository with prepared statements
func NewManufacturerRepository(db *sql.DB) (*repository, error) {
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

	return &repository{
		db:          db,
		stmtGetByID: stmtGetByID,
		stmtGetAll:  stmtGetAll,
	}, nil
}

// GetManufacturerByID retrieves a manufacturer by ID
func (r *repository) GetManufacturerByID(ctx context.Context, id string) (*domain.Manufacturer, error) {
	var manufacturer domain.Manufacturer

	err := r.stmtGetByID.QueryRowContext(ctx, id).Scan(
		&manufacturer.ID,
		&manufacturer.Name,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrManufacturerNotFound
		}
		return nil, err
	}

	return &manufacturer, nil
}

// ListManufacturers retrieves all manufacturers
func (r *repository) ListManufacturers(ctx context.Context) ([]domain.Manufacturer, error) {
	rows, err := r.stmtGetAll.QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var manufacturers []domain.Manufacturer
	for rows.Next() {
		var manufacturer domain.Manufacturer

		if err := rows.Scan(
			&manufacturer.ID,
			&manufacturer.Name,
		); err != nil {
			return nil, err
		}

		manufacturers = append(manufacturers, manufacturer)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return manufacturers, nil
}
