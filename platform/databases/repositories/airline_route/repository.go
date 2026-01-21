package airline_route

import (
	"context"
	"database/sql"

	"github.com/champion19/flighthours-api/core/ports/output"
	"github.com/champion19/flighthours-api/platform/logger"
)

const (
	// Query with JOINs to get denormalized data from route and airline tables
	QueryByID = `
		SELECT
			ar.id,
			ar.route_id,
			ar.airline_id,
			ar.status,
			a.airline_code,
			a.airline_name,
			ao.iata_code AS origin_iata_code,
			ad.iata_code AS destination_iata_code,
			CONCAT(ao.iata_code, '-', ad.iata_code) AS route_code,
			ao.name AS origin_airport_name,
			ad.name AS destination_airport_name,
			r.airport_type,
			r.estimated_flight_time
		FROM airline_route ar
		JOIN airline a ON ar.airline_id = a.id
		JOIN route r ON ar.route_id = r.id
		JOIN airport ao ON r.origin_airport_id = ao.id
		JOIN airport ad ON r.destination_airport_id = ad.id
		WHERE ar.id = ?
		LIMIT 1
	`

	QueryGetAll = `
		SELECT
			ar.id,
			ar.route_id,
			ar.airline_id,
			ar.status,
			a.airline_code,
			a.airline_name,
			ao.iata_code AS origin_iata_code,
			ad.iata_code AS destination_iata_code,
			CONCAT(ao.iata_code, '-', ad.iata_code) AS route_code,
			ao.name AS origin_airport_name,
			ad.name AS destination_airport_name,
			r.airport_type,
			r.estimated_flight_time
		FROM airline_route ar
		JOIN airline a ON ar.airline_id = a.id
		JOIN route r ON ar.route_id = r.id
		JOIN airport ao ON r.origin_airport_id = ao.id
		JOIN airport ad ON r.destination_airport_id = ad.id
		ORDER BY a.airline_code, ao.iata_code, ad.iata_code
	`

	QueryGetByAirlineID = `
		SELECT
			ar.id,
			ar.route_id,
			ar.airline_id,
			ar.status,
			a.airline_code,
			a.airline_name,
			ao.iata_code AS origin_iata_code,
			ad.iata_code AS destination_iata_code,
			CONCAT(ao.iata_code, '-', ad.iata_code) AS route_code,
			ao.name AS origin_airport_name,
			ad.name AS destination_airport_name,
			r.airport_type,
			r.estimated_flight_time
		FROM airline_route ar
		JOIN airline a ON ar.airline_id = a.id
		JOIN route r ON ar.route_id = r.id
		JOIN airport ao ON r.origin_airport_id = ao.id
		JOIN airport ad ON r.destination_airport_id = ad.id
		WHERE ar.airline_id = ?
		ORDER BY ao.iata_code, ad.iata_code
	`

	QueryGetByAirlineCode = `
		SELECT
			ar.id,
			ar.route_id,
			ar.airline_id,
			ar.status,
			a.airline_code,
			a.airline_name,
			ao.iata_code AS origin_iata_code,
			ad.iata_code AS destination_iata_code,
			CONCAT(ao.iata_code, '-', ad.iata_code) AS route_code,
			ao.name AS origin_airport_name,
			ad.name AS destination_airport_name,
			r.airport_type,
			r.estimated_flight_time
		FROM airline_route ar
		JOIN airline a ON ar.airline_id = a.id
		JOIN route r ON ar.route_id = r.id
		JOIN airport ao ON r.origin_airport_id = ao.id
		JOIN airport ad ON r.destination_airport_id = ad.id
		WHERE a.airline_code = ?
		ORDER BY ao.iata_code, ad.iata_code
	`

	QueryGetByStatus = `
		SELECT
			ar.id,
			ar.route_id,
			ar.airline_id,
			ar.status,
			a.airline_code,
			a.airline_name,
			ao.iata_code AS origin_iata_code,
			ad.iata_code AS destination_iata_code,
			CONCAT(ao.iata_code, '-', ad.iata_code) AS route_code,
			ao.name AS origin_airport_name,
			ad.name AS destination_airport_name,
			r.airport_type,
			r.estimated_flight_time
		FROM airline_route ar
		JOIN airline a ON ar.airline_id = a.id
		JOIN route r ON ar.route_id = r.id
		JOIN airport ao ON r.origin_airport_id = ao.id
		JOIN airport ad ON r.destination_airport_id = ad.id
		WHERE ar.status = ?
		ORDER BY a.airline_code, ao.iata_code, ad.iata_code
	`

	QueryUpdateStatus = `
		UPDATE airline_route
		SET status = ?
		WHERE id = ?
	`
)

var log logger.Logger = logger.NewSlogLogger()

type repository struct {
	stmtGetByID          *sql.Stmt
	stmtGetAll           *sql.Stmt
	stmtGetByAirlineID   *sql.Stmt
	stmtGetByAirlineCode *sql.Stmt
	stmtGetByStatus      *sql.Stmt
	stmtUpdateStatus     *sql.Stmt
	db                   *sql.DB
}

// NewAirlineRouteRepository creates a new airline route repository with prepared statements
func NewAirlineRouteRepository(db *sql.DB) (*repository, error) {
	if db == nil {
		return nil, sql.ErrConnDone
	}

	stmtGetByID, err := db.Prepare(QueryByID)
	if err != nil {
		log.Error(logger.LogAirlineRouteRepoInitError, "error preparing statement", err)
		return nil, err
	}

	stmtGetAll, err := db.Prepare(QueryGetAll)
	if err != nil {
		log.Error(logger.LogAirlineRouteRepoInitError, "error preparing statement", err)
		return nil, err
	}

	stmtGetByAirlineID, err := db.Prepare(QueryGetByAirlineID)
	if err != nil {
		log.Error(logger.LogAirlineRouteRepoInitError, "error preparing statement", err)
		return nil, err
	}

	stmtGetByAirlineCode, err := db.Prepare(QueryGetByAirlineCode)
	if err != nil {
		log.Error(logger.LogAirlineRouteRepoInitError, "error preparing statement", err)
		return nil, err
	}

	stmtGetByStatus, err := db.Prepare(QueryGetByStatus)
	if err != nil {
		log.Error(logger.LogAirlineRouteRepoInitError, "error preparing statement", err)
		return nil, err
	}

	stmtUpdateStatus, err := db.Prepare(QueryUpdateStatus)
	if err != nil {
		log.Error(logger.LogAirlineRouteRepoInitError, "error preparing statement", err)
		return nil, err
	}

	log.Info(logger.LogAirlineRouteRepoInitOK)

	return &repository{
		db:                   db,
		stmtGetByID:          stmtGetByID,
		stmtGetAll:           stmtGetAll,
		stmtGetByAirlineID:   stmtGetByAirlineID,
		stmtGetByAirlineCode: stmtGetByAirlineCode,
		stmtGetByStatus:      stmtGetByStatus,
		stmtUpdateStatus:     stmtUpdateStatus,
	}, nil
}

// BeginTx starts a new database transaction
func (r *repository) BeginTx(ctx context.Context) (output.Tx, error) {
	return r.db.BeginTx(ctx, nil)
}
