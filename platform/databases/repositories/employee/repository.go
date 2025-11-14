package employee

import (
	"context"
	"database/sql"

	"github.com/champion19/flighthours-api/core/ports/output"
)

const (
	QuerySave    = "INSERT INTO employee(id,name,airline,email,identification_number,bp,start_date,end_date,active,role,keycloak_user_id) VALUES(?,?,?,?,?,?,?,?,?,?,?)"
	QueryByEmail = "SELECT id,name,airline,email,identification_number,bp,start_date,end_date,active,role,keycloak_user_id FROM employee WHERE email=? LIMIT 1"
	QueryByID    = "SELECT id,name,airline,email,identification_number,bp,start_date,end_date,active,role,keycloak_user_id FROM employee WHERE id=? LIMIT 1"
	QueryUpdate  = "UPDATE employee SET name=?,airline=?,email=?,identification_number=?,bp=?,start_date=?,end_date=?,active=?,role=?,keycloak_user_id=? WHERE id=?"
	QueryDelete  = "DELETE FROM employee WHERE id=?"
)

type sqlTx struct {
	*sql.Tx
}

func (t *sqlTx) Commit() error {
	return t.Tx.Commit()
}

func (t *sqlTx) Rollback() error {
	return t.Tx.Rollback()
}
type repository struct {
	keycloak       output.AuthClient
	db             *sql.DB
}

func NewClientRepository(db *sql.DB, keycloak output.AuthClient) (*repository, error) {
	return &repository{
		keycloak:       keycloak,
		db:             db,

	}, nil
}

func (r *repository) BeginTx(ctx context.Context) (output.Tx, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &sqlTx{Tx: tx}, nil
}
