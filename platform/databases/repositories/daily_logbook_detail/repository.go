package daily_logbook_detail

import (
	"context"
	"database/sql"

	"github.com/champion19/flighthours-api/core/ports/output"
	"github.com/champion19/flighthours-api/platform/databases/common"
	"github.com/champion19/flighthours-api/platform/logger"
)

const (
	// Query for getting a detail by ID with JOINs for denormalized data
	QueryByID = `
		SELECT
			dld.id,
			dld.daily_logbook_id,
			dld.flight_real_date,
			dld.flight_number,
			dld.airline_route_id,
			dld.actual_aircraft_registration_id,
			dld.passengers,
			dld.out_time,
			dld.takeoff_time,
			dld.landing_time,
			dld.in_time,
			dld.pilot_role,
			dld.companion_name,
			dld.air_time,
			dld.block_time,
			dld.duty_time,
			dld.approach_type,
			dld.flight_type,
			dld.employee_logbook_id,
			dl.log_date,
			ar.license_plate,
			am.model_name,
			CONCAT(orig.iata_code, '-', dest.iata_code) as route_code,
			orig.iata_code as origin_iata_code,
			dest.iata_code as destination_iata_code,
			airl.airline_code
		FROM daily_logbook_detail dld
		INNER JOIN daily_logbook dl ON dld.daily_logbook_id = dl.id
		INNER JOIN aircraft_registration ar ON dld.actual_aircraft_registration_id = ar.id
		INNER JOIN aircraft_model am ON ar.aircraft_model_id = am.id
		INNER JOIN airline_route alr ON dld.airline_route_id = alr.id
		INNER JOIN route r ON alr.route_id = r.id
		INNER JOIN airport orig ON r.origin_airport_id = orig.id
		INNER JOIN airport dest ON r.destination_airport_id = dest.id
		INNER JOIN airline airl ON alr.airline_id = airl.id
		WHERE dld.id = ?
		LIMIT 1
	`

	// Query for listing details by logbook ID
	QueryByLogbook = `
		SELECT
			dld.id,
			dld.daily_logbook_id,
			dld.flight_real_date,
			dld.flight_number,
			dld.airline_route_id,
			dld.actual_aircraft_registration_id,
			dld.passengers,
			dld.out_time,
			dld.takeoff_time,
			dld.landing_time,
			dld.in_time,
			dld.pilot_role,
			dld.companion_name,
			dld.air_time,
			dld.block_time,
			dld.duty_time,
			dld.approach_type,
			dld.flight_type,
			dld.employee_logbook_id,
			dl.log_date,
			ar.license_plate,
			am.model_name,
			CONCAT(orig.iata_code, '-', dest.iata_code) as route_code,
			orig.iata_code as origin_iata_code,
			dest.iata_code as destination_iata_code,
			airl.airline_code
		FROM daily_logbook_detail dld
		INNER JOIN daily_logbook dl ON dld.daily_logbook_id = dl.id
		INNER JOIN aircraft_registration ar ON dld.actual_aircraft_registration_id = ar.id
		INNER JOIN aircraft_model am ON ar.aircraft_model_id = am.id
		INNER JOIN airline_route alr ON dld.airline_route_id = alr.id
		INNER JOIN route r ON alr.route_id = r.id
		INNER JOIN airport orig ON r.origin_airport_id = orig.id
		INNER JOIN airport dest ON r.destination_airport_id = dest.id
		INNER JOIN airline airl ON alr.airline_id = airl.id
		WHERE dld.daily_logbook_id = ?
		ORDER BY dld.out_time ASC
	`

	// Insert query
	QueryInsert = `
		INSERT INTO daily_logbook_detail (
			id, daily_logbook_id, flight_real_date, flight_number,
			airline_route_id, actual_aircraft_registration_id, passengers,
			out_time, takeoff_time, landing_time, in_time,
			pilot_role, companion_name,
			air_time, block_time, duty_time,
			approach_type, flight_type, employee_logbook_id
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	// Update query
	QueryUpdate = `
		UPDATE daily_logbook_detail SET
			flight_real_date = ?,
			flight_number = ?,
			airline_route_id = ?,
			actual_aircraft_registration_id = ?,
			passengers = ?,
			out_time = ?,
			takeoff_time = ?,
			landing_time = ?,
			in_time = ?,
			pilot_role = ?,
			companion_name = ?,
			air_time = ?,
			block_time = ?,
			duty_time = ?,
			approach_type = ?,
			flight_type = ?
		WHERE id = ?
	`

	// Delete query
	QueryDelete = `DELETE FROM daily_logbook_detail WHERE id = ?`
)

var log logger.Logger = logger.NewSlogLogger()

type repository struct {
	stmtGetByID      *sql.Stmt
	stmtGetByLogbook *sql.Stmt
	stmtInsert       *sql.Stmt
	stmtUpdate       *sql.Stmt
	stmtDelete       *sql.Stmt
	db               *sql.DB
}

// NewDailyLogbookDetailRepository creates a new daily logbook detail repository with prepared statements
func NewDailyLogbookDetailRepository(db *sql.DB) (*repository, error) {
	if db == nil {
		return nil, sql.ErrConnDone
	}

	stmtGetByID, err := db.Prepare(QueryByID)
	if err != nil {
		log.Error(logger.LogDailyLogbookDetailRepoInitError, "error preparing statement", err)
		return nil, err
	}

	stmtGetByLogbook, err := db.Prepare(QueryByLogbook)
	if err != nil {
		log.Error(logger.LogDailyLogbookDetailRepoInitError, "error preparing statement", err)
		return nil, err
	}

	stmtInsert, err := db.Prepare(QueryInsert)
	if err != nil {
		log.Error(logger.LogDailyLogbookDetailRepoInitError, "error preparing statement", err)
		return nil, err
	}

	stmtUpdate, err := db.Prepare(QueryUpdate)
	if err != nil {
		log.Error(logger.LogDailyLogbookDetailRepoInitError, "error preparing statement", err)
		return nil, err
	}

	stmtDelete, err := db.Prepare(QueryDelete)
	if err != nil {
		log.Error(logger.LogDailyLogbookDetailRepoInitError, "error preparing statement", err)
		return nil, err
	}

	log.Info(logger.LogDailyLogbookDetailRepoInitOK)

	return &repository{
		db:               db,
		stmtGetByID:      stmtGetByID,
		stmtGetByLogbook: stmtGetByLogbook,
		stmtInsert:       stmtInsert,
		stmtUpdate:       stmtUpdate,
		stmtDelete:       stmtDelete,
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
