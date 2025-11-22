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
	tx *sql.Tx
	closed bool
}

func (t *sqlTx) Commit() error {
	if t.closed {
		panic("sqlTx: commit on closed")}
	t.closed = true
	return t.tx.Commit()
}

func (t *sqlTx) Rollback() error {
	if t.closed {
		panic("sqlTx: rollback on closed")
	}
	t.closed = true
	return t.tx.Rollback()
}

func (t *sqlTx) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return t.tx.ExecContext(ctx, query, args...)
}

func (t *sqlTx) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return t.tx.QueryRowContext(ctx, query, args...)
}
type repository struct {
	db             *sql.DB
}

func NewClientRepository(db *sql.DB) (*repository, error) {
	if db == nil {
		return nil,sql.ErrConnDone
	}

	return &repository{
		db:             db,
	}, nil
}

func (r *repository) BeginTx(ctx context.Context) (output.Tx, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &sqlTx{tx: tx}, nil
}
