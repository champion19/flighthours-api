package domain

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Employee struct {
	ID                   string    `json:"id"`
	Name                 string    `json:"name"`
	Airline              string    `json:"airline"`
	Email                string    `json:"email"`
	Password             string    `json:"password"`
	Emailconfirmed       bool      `json:"emailconfirmed"`
	IdentificationNumber string    `json:"identification_number"`
	Bp                   string    `json:"bp"`
	StartDate            time.Time `json:"start_date"`
	EndDate              time.Time `json:"end_date"`
	Active               bool      `json:"active"`
	Role                 string    `json:"role,omitempty"`
	KeycloakUserID       string    `json:"keycloak_user_id,omitempty"`
}

func (e *Employee) SetID() {
	e.ID = uuid.New().String()
}

func (e *Employee) HashPassword() error {
	hash, err := bcrypt.GenerateFromPassword([]byte(e.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	e.Password = string(hash)
	return nil
}

type Airline struct {
	ID     string
	Name   string
	Code   string
	Status string
}
