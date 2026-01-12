package airline_route

import (
	"context"
	"database/sql"

	"github.com/champion19/flighthours-api/core/interactor/services/domain"
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

// GetAirlineRouteByID retrieves an airline route by ID with denormalized data
func (r *repository) GetAirlineRouteByID(ctx context.Context, id string) (*domain.AirlineRoute, error) {
	var airlineRoute domain.AirlineRoute
	var estimatedFlightTime sql.NullString

	err := r.stmtGetByID.QueryRowContext(ctx, id).Scan(
		&airlineRoute.ID,
		&airlineRoute.RouteID,
		&airlineRoute.AirlineID,
		&airlineRoute.Status,
		&airlineRoute.AirlineCode,
		&airlineRoute.AirlineName,
		&airlineRoute.OriginIataCode,
		&airlineRoute.DestinationIataCode,
		&airlineRoute.RouteCode,
		&airlineRoute.OriginAirportName,
		&airlineRoute.DestinationAirportName,
		&airlineRoute.AirportType,
		&estimatedFlightTime,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrAirlineRouteNotFound
		}
		return nil, err
	}

	if estimatedFlightTime.Valid {
		airlineRoute.EstimatedFlightTime = estimatedFlightTime.String
	}

	return &airlineRoute, nil
}

// ListAirlineRoutes retrieves all airline routes with optional filters
func (r *repository) ListAirlineRoutes(ctx context.Context, filters map[string]interface{}) ([]domain.AirlineRoute, error) {
	var rows *sql.Rows
	var err error

	// Check filters
	if airlineID, ok := filters["airline_id"]; ok {
		rows, err = r.stmtGetByAirlineID.QueryContext(ctx, airlineID)
	} else if airlineCode, ok := filters["airline_code"]; ok {
		rows, err = r.stmtGetByAirlineCode.QueryContext(ctx, airlineCode)
	} else if status, ok := filters["status"]; ok {
		rows, err = r.stmtGetByStatus.QueryContext(ctx, status)
	} else {
		rows, err = r.stmtGetAll.QueryContext(ctx)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var airlineRoutes []domain.AirlineRoute
	for rows.Next() {
		var ar domain.AirlineRoute
		var estimatedFlightTime sql.NullString

		if err := rows.Scan(
			&ar.ID,
			&ar.RouteID,
			&ar.AirlineID,
			&ar.Status,
			&ar.AirlineCode,
			&ar.AirlineName,
			&ar.OriginIataCode,
			&ar.DestinationIataCode,
			&ar.RouteCode,
			&ar.OriginAirportName,
			&ar.DestinationAirportName,
			&ar.AirportType,
			&estimatedFlightTime,
		); err != nil {
			return nil, err
		}

		if estimatedFlightTime.Valid {
			ar.EstimatedFlightTime = estimatedFlightTime.String
		}

		airlineRoutes = append(airlineRoutes, ar)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return airlineRoutes, nil
}

// UpdateAirlineRouteStatus updates the status of an airline route
func (r *repository) UpdateAirlineRouteStatus(ctx context.Context, tx output.Tx, id string, status bool) error {
	sqlTx, ok := tx.(*sql.Tx)
	if !ok {
		return domain.ErrInvalidTransaction
	}

	stmt := sqlTx.StmtContext(ctx, r.stmtUpdateStatus)
	result, err := stmt.ExecContext(ctx, status, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return domain.ErrAirlineRouteNotFound
	}

	return nil
}
