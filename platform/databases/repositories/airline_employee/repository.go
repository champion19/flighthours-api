package airline_employee

import (
	"context"
	"database/sql"

	"github.com/champion19/flighthours-api/core/ports/output"
	"github.com/champion19/flighthours-api/platform/databases/common"
	"github.com/champion19/flighthours-api/platform/logger"
)

const (
	// Query with JOIN to get employee with airline data (airline NOT NULL)
	QueryByID = `
		SELECT
			e.id,
			e.name,
			e.airline,
			a.airline_name,
			a.airline_code,
			e.email,
			e.identification_number,
			e.bp,
			e.start_date,
			e.end_date,
			e.active,
			e.role,
			e.keycloak_user_id
		FROM employee e
		JOIN airline a ON e.airline = a.id
		WHERE e.id = ? AND e.airline IS NOT NULL
		LIMIT 1
	`

	QueryGetAll = `
		SELECT
			e.id,
			e.name,
			e.airline,
			a.airline_name,
			a.airline_code,
			e.email,
			e.identification_number,
			e.bp,
			e.start_date,
			e.end_date,
			e.active,
			e.role,
			e.keycloak_user_id
		FROM employee e
		JOIN airline a ON e.airline = a.id
		WHERE e.airline IS NOT NULL
		ORDER BY a.airline_name, e.name
	`

	QueryGetByAirlineID = `
		SELECT
			e.id,
			e.name,
			e.airline,
			a.airline_name,
			a.airline_code,
			e.email,
			e.identification_number,
			e.bp,
			e.start_date,
			e.end_date,
			e.active,
			e.role,
			e.keycloak_user_id
		FROM employee e
		JOIN airline a ON e.airline = a.id
		WHERE e.airline = ?
		ORDER BY e.name
	`

	QueryGetByStatus = `
		SELECT
			e.id,
			e.name,
			e.airline,
			a.airline_name,
			a.airline_code,
			e.email,
			e.identification_number,
			e.bp,
			e.start_date,
			e.end_date,
			e.active,
			e.role,
			e.keycloak_user_id
		FROM employee e
		JOIN airline a ON e.airline = a.id
		WHERE e.airline IS NOT NULL AND e.active = ?
		ORDER BY a.airline_name, e.name
	`

	QueryInsert = `
		INSERT INTO employee (id, name, airline, email, identification_number, bp, start_date, end_date, active, role, keycloak_user_id)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	QueryUpdate = `
		UPDATE employee SET
			name = ?,
			airline = ?,
			email = ?,
			identification_number = ?,
			bp = ?,
			start_date = ?,
			end_date = ?,
			active = ?,
			role = ?,
			keycloak_user_id = ?
		WHERE id = ? AND airline IS NOT NULL
	`

	QueryUpdateStatus = `
		UPDATE employee SET active = ?
		WHERE id = ? AND airline IS NOT NULL
	`
)

var log logger.Logger = logger.NewSlogLogger()

type repository struct {
	stmtGetByID        *sql.Stmt
	stmtGetAll         *sql.Stmt
	stmtGetByAirlineID *sql.Stmt
	stmtGetByStatus    *sql.Stmt
	stmtInsert         *sql.Stmt
	stmtUpdate         *sql.Stmt
	stmtUpdateStatus   *sql.Stmt
	db                 *sql.DB
}

// NewAirlineEmployeeRepository creates a new airline employee repository with prepared statements
func NewAirlineEmployeeRepository(db *sql.DB) (*repository, error) {
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

	stmtGetByAirlineID, err := db.Prepare(QueryGetByAirlineID)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, "error preparing statement", err)
		return nil, err
	}

	stmtGetByStatus, err := db.Prepare(QueryGetByStatus)
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

	stmtUpdateStatus, err := db.Prepare(QueryUpdateStatus)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, "error preparing statement", err)
		return nil, err
	}

	log.Info(logger.LogDatabaseAvailable, "repository", "airline_employee")

	return &repository{
		db:                 db,
		stmtGetByID:        stmtGetByID,
		stmtGetAll:         stmtGetAll,
		stmtGetByAirlineID: stmtGetByAirlineID,
		stmtGetByStatus:    stmtGetByStatus,
		stmtInsert:         stmtInsert,
		stmtUpdate:         stmtUpdate,
		stmtUpdateStatus:   stmtUpdateStatus,
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
