package airline_employee

import (
	"time"

	domain "github.com/champion19/flighthours-api/core/interactor/services/domain"
)

// AirlineEmployee is the database entity for airline employee (employee with airline NOT NULL)
type AirlineEmployee struct {
	ID                   string    `db:"id"`
	Name                 string    `db:"name"`
	AirlineID            string    `db:"airline"`
	AirlineName          string    `db:"airline_name"`
	AirlineCode          string    `db:"airline_code"`
	Email                string    `db:"email"`
	IdentificationNumber string    `db:"identification_number"`
	Bp                   *string   `db:"bp"`
	StartDate            time.Time `db:"start_date"`
	EndDate              time.Time `db:"end_date"`
	Active               bool      `db:"active"`
	Role                 string    `db:"role"`
	KeycloakUserID       *string   `db:"keycloak_user_id"`
}

// ToDomain converts the database entity to domain model
func (e *AirlineEmployee) ToDomain() *domain.AirlineEmployee {
	bp := ""
	if e.Bp != nil {
		bp = *e.Bp
	}
	keycloakUserID := ""
	if e.KeycloakUserID != nil {
		keycloakUserID = *e.KeycloakUserID
	}

	return &domain.AirlineEmployee{
		ID:                   e.ID,
		Name:                 e.Name,
		AirlineID:            e.AirlineID,
		AirlineName:          e.AirlineName,
		AirlineCode:          e.AirlineCode,
		Email:                e.Email,
		IdentificationNumber: e.IdentificationNumber,
		Bp:                   bp,
		StartDate:            e.StartDate,
		EndDate:              e.EndDate,
		Active:               e.Active,
		Role:                 e.Role,
		KeycloakUserID:       keycloakUserID,
	}
}

// stringPtrOrNil returns nil for empty strings, otherwise a pointer to the string
func stringPtrOrNil(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// FromDomain converts a domain model to database entity
func FromDomain(domainEmployee *domain.AirlineEmployee) *AirlineEmployee {
	return &AirlineEmployee{
		ID:                   domainEmployee.ID,
		Name:                 domainEmployee.Name,
		AirlineID:            domainEmployee.AirlineID,
		AirlineName:          domainEmployee.AirlineName,
		AirlineCode:          domainEmployee.AirlineCode,
		Email:                domainEmployee.Email,
		IdentificationNumber: domainEmployee.IdentificationNumber,
		Bp:                   stringPtrOrNil(domainEmployee.Bp),
		StartDate:            domainEmployee.StartDate,
		EndDate:              domainEmployee.EndDate,
		Active:               domainEmployee.Active,
		Role:                 domainEmployee.Role,
		KeycloakUserID:       stringPtrOrNil(domainEmployee.KeycloakUserID),
	}
}
