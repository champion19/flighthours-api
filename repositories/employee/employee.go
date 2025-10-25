package employee

import (
	"time"

	"github.com/champion19/Flighthours_backend/core/domain"
)

type Employee struct {
	ID                   string    `db:"id"`
	Name                 string    `db:"name"`
	Airline              string    `db:"airline"`
	Email                string    `db:"email"`
	Password             string    `db:"password"`
	Emailconfirmed       bool      `db:"emailconfirmed"`
	IdentificationNumber string    `db:"identification_number"`
	Bp                   string    `db:"bp"`
	StartDate            time.Time `db:"start_date"`
	EndDate              time.Time `db:"end_date"`
	Active               bool      `db:"active"`
	Role                 string    `db:"role"`
	KeycloakUserID       string    `db:"keycloak_user_id"`
}

func (e Employee) ToDomain() domain.Employee {
	return domain.Employee{
		ID:                   e.ID,
		Name:                 e.Name,
		Email:                e.Email,
		Airline:              e.Airline,
		Password:             e.Password,
		Emailconfirmed:       e.Emailconfirmed,
		IdentificationNumber: e.IdentificationNumber,
		Bp:                   e.Bp,
		StartDate:            e.StartDate,
		EndDate:              e.EndDate,
		Active:               e.Active,
		Role:                 e.Role,
		KeycloakUserID:       e.KeycloakUserID,
	}
}

func FromDomain(domainEmployee domain.Employee) Employee {
	return Employee{
		ID:                   domainEmployee.ID,
		Name:                 domainEmployee.Name,
		Email:                domainEmployee.Email,
		Airline:              domainEmployee.Airline,
		Password:             domainEmployee.Password,
		Emailconfirmed:       domainEmployee.Emailconfirmed,
		IdentificationNumber: domainEmployee.IdentificationNumber,
		Bp:                   domainEmployee.Bp,
		StartDate:            domainEmployee.StartDate,
		EndDate:              domainEmployee.EndDate,
		Active:               domainEmployee.Active,
		Role:                 domainEmployee.Role,
		KeycloakUserID:       domainEmployee.KeycloakUserID,
	}
}


