package aircraft_registration

import (
	"context"
	"database/sql"

	"github.com/champion19/flighthours-api/core/ports/output"
	"github.com/champion19/flighthours-api/platform/databases/common"
	"github.com/champion19/flighthours-api/platform/logger"
)

const (
	// Query with JOINs to get denormalized data (Numero, Modelo, Aerolinea)
	QueryByID = `
		SELECT
			ar.id,
			ar.license_plate,
			ar.aircraft_model_id,
			ar.airline_id,
			am.model_name,
			a.airline_name
		FROM aircraft_registration ar
		INNER JOIN aircraft_model am ON ar.aircraft_model_id = am.id
		INNER JOIN airline a ON ar.airline_id = a.id
		WHERE ar.id = ?
		LIMIT 1`

	QueryGetAll = `
		SELECT
			ar.id,
			ar.license_plate,
			ar.aircraft_model_id,
			ar.airline_id,
			am.model_name,
			a.airline_name
		FROM aircraft_registration ar
		INNER JOIN aircraft_model am ON ar.aircraft_model_id = am.id
		INNER JOIN airline a ON ar.airline_id = a.id
		ORDER BY ar.license_plate`

	QueryGetByAirline = `
		SELECT
			ar.id,
			ar.license_plate,
			ar.aircraft_model_id,
			ar.airline_id,
			am.model_name,
			a.airline_name
		FROM aircraft_registration ar
		INNER JOIN aircraft_model am ON ar.aircraft_model_id = am.id
		INNER JOIN airline a ON ar.airline_id = a.id
		WHERE ar.airline_id = ?
		ORDER BY ar.license_plate`

	QueryInsert = `INSERT INTO aircraft_registration (id, license_plate, aircraft_model_id, airline_id) VALUES (?, ?, ?, ?)`
	QueryUpdate = `UPDATE aircraft_registration SET license_plate = ?, aircraft_model_id = ?, airline_id = ? WHERE id = ?`
)

var log logger.Logger = logger.NewSlogLogger()

type repository struct {
	stmtGetByID      *sql.Stmt
	stmtGetAll       *sql.Stmt
	stmtGetByAirline *sql.Stmt
	stmtInsert       *sql.Stmt
	stmtUpdate       *sql.Stmt
	db               *sql.DB
}

// NewAircraftRegistrationRepository creates a new aircraft registration repository with prepared statements
func NewAircraftRegistrationRepository(db *sql.DB) (*repository, error) {
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

	stmtGetByAirline, err := db.Prepare(QueryGetByAirline)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, "error preparing statement", err)
		return nil, err
	}

	stmtInsert, err := db.Prepare(QueryInsert)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, "error preparing statement", err)
		return nil, err
	}

	stmtUpdate, err := db.Prepare(QueryUpdate)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, "error preparing statement", err)
		return nil, err
	}

	return &repository{
		db:               db,
		stmtGetByID:      stmtGetByID,
		stmtGetAll:       stmtGetAll,
		stmtGetByAirline: stmtGetByAirline,
		stmtInsert:       stmtInsert,
		stmtUpdate:       stmtUpdate,
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
