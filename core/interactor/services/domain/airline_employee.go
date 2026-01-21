package domain

import (
	"time"

	"github.com/google/uuid"
)

// AirlineEmployee represents an employee assigned to a specific airline
// This is a specialized view of Employee where airline IS NOT NULL
// Includes denormalized airline data for convenience
type AirlineEmployee struct {
	ID                   string    `json:"id"`
	Name                 string    `json:"name"`
	AirlineID            string    `json:"airline_id"`
	AirlineName          string    `json:"airline_name"`
	AirlineCode          string    `json:"airline_code"`
	Email                string    `json:"email"`
	IdentificationNumber string    `json:"identification_number"`
	Bp                   string    `json:"bp"`
	StartDate            time.Time `json:"start_date"`
	EndDate              time.Time `json:"end_date"`
	Active               bool      `json:"active"`
	Role                 string    `json:"role"`
	KeycloakUserID       string    `json:"keycloak_user_id,omitempty"`
}

func (e *AirlineEmployee) SetID() {
	e.ID = uuid.New().String()
}

func (e *AirlineEmployee) ToLogger() []string {
	return []string{
		"id:" + e.ID,
		"email:" + e.Email,
		"role:" + e.Role,
		"airline_id:" + e.AirlineID,
		"airline_code:" + e.AirlineCode,
	}
}
