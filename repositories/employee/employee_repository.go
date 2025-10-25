package employee

import (
	"database/sql"
	"fmt"

	domain "github.com/champion19/Flighthours_backend/core/domain"
	"github.com/champion19/Flighthours_backend/core/ports"
	mysql "github.com/go-sql-driver/mysql"
)

const (
	QuerySave    = "INSERT INTO employee(id,name,airline,email,password,email_confirmed,identification_number,bp,start_date,end_date,active,role,keycloak_user_id) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?)"
	QueryByEmail = "SELECT id,name,airline,email,password,email_confirmed,identification_number,bp,start_date,end_date,active,role,keycloak_user_id FROM employee WHERE email=? LIMIT 1"
	QueryByID    = "SELECT id,name,airline,email,password,email_confirmed,identification_number,bp,start_date,end_date,active,role,keycloak_user_id FROM employee WHERE id=? LIMIT 1"
	QueryUpdate  = "UPDATE employee SET name=?,airline=?,email=?,password=?,email_confirmed=?,identification_number=?,bp=?,start_date=?,end_date=?,active=?,role=?,keycloak_user_id=? WHERE id=?"
	QueryDelete  = "DELETE FROM employee WHERE id=?"
)

type repository struct {
	db             *sql.DB
	stmtSave       *sql.Stmt
	stmtGetByEmail *sql.Stmt
	stmtGetByID    *sql.Stmt
	stmtUpdate     *sql.Stmt
	stmtDelete     *sql.Stmt
}

func NewRepository(db *sql.DB) (ports.Repository, error) {
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
		db:             db,
		stmtSave:       stmtSave,
		stmtGetByEmail: stmtGetByEmail,
		stmtGetByID:    stmtGetByID,
		stmtUpdate:     stmtUpdate,
		stmtDelete:     stmtDelete,
	}, nil
}

func (r *repository) GetEmployeeByEmail(email string) (*domain.Employee, error) {
	var e Employee
	err := r.stmtGetByEmail.QueryRow(email).Scan(
		&e.ID,
		&e.Name,
		&e.Airline,
		&e.Email,
		&e.Password,
		&e.Emailconfirmed,
		&e.IdentificationNumber,
		&e.Bp,
		&e.StartDate,
		&e.EndDate,
		&e.Active,
		&e.Role,
		&e.KeycloakUserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrPersonNotFound
		}
		return nil, err
	}
	d := e.ToDomain()
	return &d, nil
}

func (r *repository) GetEmployeeByID(id string) (*domain.Employee, error) {
	var e Employee
	err := r.stmtGetByID.QueryRow(id).Scan(
		&e.ID,
		&e.Name,
		&e.Airline,
		&e.Email,
		&e.Password,
		&e.Emailconfirmed,
		&e.IdentificationNumber,
		&e.Bp, &e.StartDate,
		&e.EndDate,
		&e.Active,
		&e.Role,
		&e.KeycloakUserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrPersonNotFound
		}
		return nil, err
	}
	d := e.ToDomain()
	return &d, nil
}

func (r *repository) Save(employee domain.Employee) error {

employeeToSave := FromDomain(employee)
	_, err := r.stmtSave.Exec(
		employeeToSave.ID,
		employeeToSave.Name,
		employeeToSave.Airline,
		employeeToSave.Email,
		employeeToSave.Password,
		employeeToSave.Emailconfirmed,
		employeeToSave.IdentificationNumber,
		employeeToSave.Bp,
		employeeToSave.StartDate,
		employeeToSave.EndDate,
		employeeToSave.Active,
		employeeToSave.Role,
		employeeToSave.KeycloakUserID,
	)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
			return domain.ErrDuplicateUser
		} else {
			return domain.ErrUserCannotSave
		}
	}

	return nil
}

func (r *repository) UpdateEmployee(employee domain.Employee) error {
	employeeToUpdate := FromDomain(employee)

	_, err := r.stmtUpdate.Exec(
		employeeToUpdate.Name,
		employeeToUpdate.Airline,
		employeeToUpdate.Email,
		employeeToUpdate.Password,
		employeeToUpdate.Emailconfirmed,
		employeeToUpdate.IdentificationNumber,
		employeeToUpdate.Bp,
		employeeToUpdate.StartDate,
		employeeToUpdate.EndDate,
		employeeToUpdate.Active,
		employeeToUpdate.Role,
	  employeeToUpdate.KeycloakUserID,
	  employeeToUpdate.ID,
	)
	if err != nil {
		return domain.ErrUserCannotSave
	}

	return nil
}

func (r *repository) DeleteEmployee(id string) error {
	result,err:=r.db.Exec(QueryDelete,id)
		if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected,err:=result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return domain.ErrPersonNotFound
	}
	return nil
}
