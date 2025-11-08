package employee

import (
	"database/sql"

	"github.com/champion19/flighthours-api/core/interactor/services/domain"
	"github.com/champion19/flighthours-api/core/ports/output"
)

const (
	QuerySave    = "INSERT INTO employee(id,name,airline,email,identification_number,bp,start_date,end_date,active,role,keycloak_user_id) VALUES(?,?,?,?,?,?,?,?,?,?,?)"
	QueryByEmail = "SELECT id,name,airline,email,identification_number,bp,start_date,end_date,active,role,keycloak_user_id FROM employee WHERE email=? LIMIT 1"
	QueryByID    = "SELECT id,name,airline,email,identification_number,bp,start_date,end_date,active,role,keycloak_user_id FROM employee WHERE id=? LIMIT 1"
	QueryUpdate  = "UPDATE employee SET name=?,airline=?,email=?,identification_number=?,bp=?,start_date=?,end_date=?,active=?,role=?,keycloak_user_id=? WHERE id=?"
	QueryDelete  = "DELETE FROM employee WHERE id=?"
)

type repository struct {
	keycloak       output.AuthClient
	db             *sql.DB
	stmtSave       *sql.Stmt
	stmtGetByEmail *sql.Stmt
	stmtGetByID    *sql.Stmt
	stmtUpdate     *sql.Stmt
	stmtDelete     *sql.Stmt
}

func NewRepository(db *sql.DB, keycloak output.AuthClient) (*repository, error) {
	stmtSave, err := db.Prepare(QuerySave)
	if err != nil {
		return nil, domain.ErrUserCannotSave
	}
	stmtGetByEmail, err := db.Prepare(QueryByEmail)
	if err != nil {
		return nil, domain.ErrUserCannotSave
	}
	stmtGetByID, err := db.Prepare(QueryByID)
	if err != nil {
		return nil, domain.ErrUserCannotSave
	}
	stmtUpdate, err := db.Prepare(QueryUpdate)
	if err != nil {
		return nil, domain.ErrUserCannotSave
	}
	stmtDelete, err := db.Prepare(QueryDelete)
	if err != nil {
		return nil, domain.ErrUserCannotSave
	}
	return &repository{
		keycloak:       keycloak,
		db:             db,
		stmtSave:       stmtSave,
		stmtGetByEmail: stmtGetByEmail,
		stmtGetByID:    stmtGetByID,
		stmtUpdate:     stmtUpdate,
		stmtDelete:     stmtDelete,
	}, nil
}
