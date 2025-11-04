package newemployee

import (
	"database/sql"

	"github.com/champion19/flighthours-api/core/ports"
)



const (
	QuerySave    = "INSERT INTO employee(id,name,airline,email,identification_number,bp,start_date,end_date,active,role,keycloak_user_id) VALUES(?,?,?,?,?,?,?,?,?,?,?)"
	QueryByEmail = "SELECT id,name,airline,email,identification_number,bp,start_date,end_date,active,role,keycloak_user_id FROM employee WHERE email=? LIMIT 1"
	QueryByID    = "SELECT id,name,airline,email,identification_number,bp,start_date,end_date,active,role,keycloak_user_id FROM employee WHERE id=? LIMIT 1"
	QueryUpdate  = "UPDATE employee SET name=?,airline=?,email=?,identification_number=?,bp=?,start_date=?,end_date=?,active=?,role=?,keycloak_user_id=? WHERE id=?"
	QueryDelete  = "DELETE FROM employee WHERE id=?"
)

type repository struct {
	keycloak ports.AuthClient
	db *sql.DB
}

func NewClient(db *sql.DB,keycloak ports.AuthClient)(*repository,error){
	return &repository{
		keycloak:keycloak,

		db:db,
	}, nil
}
