package route

import (
	"context"
	"database/sql"

	"github.com/champion19/flighthours-api/core/interactor/services/domain"
	"github.com/champion19/flighthours-api/platform/logger"
)

const (
	// Query with JOINs to get denormalized data from airport table
	// estimated_flight_time is returned as TIME format "HH:MM:SS"
	QueryByID = `
		SELECT
			r.id,
			r.origin_airport_id,
			ao.iata_code AS origin_iata_code,
			ao.name AS origin_airport_name,
			r.destination_airport_id,
			ad.iata_code AS destination_iata_code,
			ad.name AS destination_airport_name,
			r.origin_country,
			r.destination_country,
			r.airport_type,
			r.estimated_flight_time,
			CONCAT(ao.iata_code, '-', ad.iata_code) AS route_code
		FROM route r
		JOIN airport ao ON r.origin_airport_id = ao.id
		JOIN airport ad ON r.destination_airport_id = ad.id
		WHERE r.id = ?
		LIMIT 1
	`

	QueryGetAll = `
		SELECT
			r.id,
			r.origin_airport_id,
			ao.iata_code AS origin_iata_code,
			ao.name AS origin_airport_name,
			r.destination_airport_id,
			ad.iata_code AS destination_iata_code,
			ad.name AS destination_airport_name,
			r.origin_country,
			r.destination_country,
			r.airport_type,
			r.estimated_flight_time,
			CONCAT(ao.iata_code, '-', ad.iata_code) AS route_code
		FROM route r
		JOIN airport ao ON r.origin_airport_id = ao.id
		JOIN airport ad ON r.destination_airport_id = ad.id
		ORDER BY ao.iata_code, ad.iata_code
	`

	QueryGetByAirportType = `
		SELECT
			r.id,
			r.origin_airport_id,
			ao.iata_code AS origin_iata_code,
			ao.name AS origin_airport_name,
			r.destination_airport_id,
			ad.iata_code AS destination_iata_code,
			ad.name AS destination_airport_name,
			r.origin_country,
			r.destination_country,
			r.airport_type,
			r.estimated_flight_time,
			CONCAT(ao.iata_code, '-', ad.iata_code) AS route_code
		FROM route r
		JOIN airport ao ON r.origin_airport_id = ao.id
		JOIN airport ad ON r.destination_airport_id = ad.id
		WHERE r.airport_type = ?
		ORDER BY ao.iata_code, ad.iata_code
	`

	QueryGetByOriginCountry = `
		SELECT
			r.id,
			r.origin_airport_id,
			ao.iata_code AS origin_iata_code,
			ao.name AS origin_airport_name,
			r.destination_airport_id,
			ad.iata_code AS destination_iata_code,
			ad.name AS destination_airport_name,
			r.origin_country,
			r.destination_country,
			r.airport_type,
			r.estimated_flight_time,
			CONCAT(ao.iata_code, '-', ad.iata_code) AS route_code
		FROM route r
		JOIN airport ao ON r.origin_airport_id = ao.id
		JOIN airport ad ON r.destination_airport_id = ad.id
		WHERE r.origin_country = ?
		ORDER BY ao.iata_code, ad.iata_code
	`
)

var log logger.Logger = logger.NewSlogLogger()

type repository struct {
	stmtGetByID            *sql.Stmt
	stmtGetAll             *sql.Stmt
	stmtGetByAirportType   *sql.Stmt
	stmtGetByOriginCountry *sql.Stmt
	db                     *sql.DB
}

// NewRouteRepository creates a new route repository with prepared statements
func NewRouteRepository(db *sql.DB) (*repository, error) {
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

	stmtGetByAirportType, err := db.Prepare(QueryGetByAirportType)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, "error preparing statement", err)
		return nil, err
	}

	stmtGetByOriginCountry, err := db.Prepare(QueryGetByOriginCountry)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, "error preparing statement", err)
		return nil, err
	}

	return &repository{
		db:                     db,
		stmtGetByID:            stmtGetByID,
		stmtGetAll:             stmtGetAll,
		stmtGetByAirportType:   stmtGetByAirportType,
		stmtGetByOriginCountry: stmtGetByOriginCountry,
	}, nil
}

// GetRouteByID retrieves a route by ID with denormalized airport data
func (r *repository) GetRouteByID(ctx context.Context, id string) (*domain.Route, error) {
	var route domain.Route
	var originCountry, destinationCountry, estimatedFlightTime sql.NullString

	err := r.stmtGetByID.QueryRowContext(ctx, id).Scan(
		&route.ID,
		&route.OriginAirportID,
		&route.OriginIataCode,
		&route.OriginAirportName,
		&route.DestinationAirportID,
		&route.DestinationIataCode,
		&route.DestinationAirportName,
		&originCountry,
		&destinationCountry,
		&route.AirportType,
		&estimatedFlightTime,
		&route.RouteCode,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrRouteNotFound
		}
		return nil, err
	}

	if originCountry.Valid {
		route.OriginCountry = originCountry.String
	}
	if destinationCountry.Valid {
		route.DestinationCountry = destinationCountry.String
	}
	if estimatedFlightTime.Valid {
		route.EstimatedFlightTime = estimatedFlightTime.String
	}

	return &route, nil
}

// ListRoutes retrieves all routes with optional filters
func (r *repository) ListRoutes(ctx context.Context, filters map[string]interface{}) ([]domain.Route, error) {
	var rows *sql.Rows
	var err error

	// Check if filtering by airport type
	if airportType, ok := filters["airport_type"]; ok {
		rows, err = r.stmtGetByAirportType.QueryContext(ctx, airportType)
	} else if originCountry, ok := filters["origin_country"]; ok {
		rows, err = r.stmtGetByOriginCountry.QueryContext(ctx, originCountry)
	} else {
		rows, err = r.stmtGetAll.QueryContext(ctx)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var routes []domain.Route
	for rows.Next() {
		var route domain.Route
		var originCountry, destinationCountry, estimatedFlightTime sql.NullString

		if err := rows.Scan(
			&route.ID,
			&route.OriginAirportID,
			&route.OriginIataCode,
			&route.OriginAirportName,
			&route.DestinationAirportID,
			&route.DestinationIataCode,
			&route.DestinationAirportName,
			&originCountry,
			&destinationCountry,
			&route.AirportType,
			&estimatedFlightTime,
			&route.RouteCode,
		); err != nil {
			return nil, err
		}

		if originCountry.Valid {
			route.OriginCountry = originCountry.String
		}
		if destinationCountry.Valid {
			route.DestinationCountry = destinationCountry.String
		}
		if estimatedFlightTime.Valid {
			route.EstimatedFlightTime = estimatedFlightTime.String
		}

		routes = append(routes, route)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return routes, nil
}
