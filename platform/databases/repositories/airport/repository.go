package airport

import (
	"context"
	"database/sql"

	"github.com/champion19/flighthours-api/core/ports/output"
	"github.com/champion19/flighthours-api/platform/databases/common"
	"github.com/champion19/flighthours-api/platform/logger"
)

const (
	QueryByID         = "SELECT id, name, city, country, iata_code, status, airport_type FROM airport WHERE id = ? LIMIT 1"
	QueryUpdateStatus = "UPDATE airport SET status = ? WHERE id = ?"
	QueryGetAll       = "SELECT id, name, city, country, iata_code, status, airport_type FROM airport ORDER BY name"
	QueryGetByStatus  = "SELECT id, name, city, country, iata_code, status, airport_type FROM airport WHERE status = ? ORDER BY name"
)

var log logger.Logger = logger.NewSlogLogger()

type repository struct {
	stmtGetByID      *sql.Stmt
	stmtUpdateStatus *sql.Stmt
	stmtGetAll       *sql.Stmt
	stmtGetByStatus  *sql.Stmt
	db               *sql.DB
}

// NewAirportRepository creates a new airport repository with prepared statements
func NewAirportRepository(db *sql.DB) (*repository, error) {
	if db == nil {
		return nil, sql.ErrConnDone
	}

	stmtGetByID, err := db.Prepare(QueryByID)
	if err != nil {
		log.Error(logger.LogAirportRepoInitError, "error preparing statement", err)
		return nil, err
	}

	stmtUpdateStatus, err := db.Prepare(QueryUpdateStatus)
	if err != nil {
		log.Error(logger.LogAirportRepoInitError, "error preparing statement", err)
		return nil, err
	}

	stmtGetAll, err := db.Prepare(QueryGetAll)
	if err != nil {
		log.Error(logger.LogAirportRepoInitError, "error preparing statement", err)
		return nil, err
	}

	stmtGetByStatus, err := db.Prepare(QueryGetByStatus)
	if err != nil {
		log.Error(logger.LogAirportRepoInitError, "error preparing statement", err)
		return nil, err
	}

	return &repository{
		db:               db,
		stmtGetByID:      stmtGetByID,
		stmtUpdateStatus: stmtUpdateStatus,
		stmtGetAll:       stmtGetAll,
		stmtGetByStatus:  stmtGetByStatus,
	}, nil
}

// BeginTx starts a new database transaction
func (r *repository) BeginTx(ctx context.Context) (output.Tx, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return common.NewSQLTx(tx), nil
}
