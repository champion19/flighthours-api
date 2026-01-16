package engine

import (
	"context"
	"database/sql"

	"github.com/champion19/flighthours-api/core/interactor/services/domain"
	"github.com/champion19/flighthours-api/platform/logger"
)

const (
	QueryByID   = "SELECT id, name FROM engine WHERE id = ? LIMIT 1"
	QueryGetAll = "SELECT id, name FROM engine ORDER BY name"
)

var log logger.Logger = logger.NewSlogLogger()

type repository struct {
	stmtGetByID *sql.Stmt
	stmtGetAll  *sql.Stmt
	db          *sql.DB
}

// NewEngineRepository creates a new engine repository with prepared statements
func NewEngineRepository(db *sql.DB) (*repository, error) {
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

// GetEngineByID retrieves an engine by ID
func (r *repository) GetEngineByID(ctx context.Context, id string) (*domain.Engine, error) {
	var engine domain.Engine

	err := r.stmtGetByID.QueryRowContext(ctx, id).Scan(
		&engine.ID,
		&engine.Name,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrEngineNotFound
		}
		return nil, err
	}

	return &engine, nil
}

// ListEngines retrieves all engines
func (r *repository) ListEngines(ctx context.Context) ([]domain.Engine, error) {
	rows, err := r.stmtGetAll.QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var engines []domain.Engine
	for rows.Next() {
		var engine domain.Engine

		if err := rows.Scan(
			&engine.ID,
			&engine.Name,
		); err != nil {
			return nil, err
		}

		engines = append(engines, engine)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return engines, nil
}
