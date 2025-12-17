package handlers

import (
	"html/template"
	"time"

	domain "github.com/champion19/flighthours-api/core/interactor/services/domain"
)

type EmployeeRequest struct {
	Name                 string `json:"name"`
	Airline              string `json:"airline"`
	Email                string `json:"email"`
	Password             string `json:"password"`
	IdentificationNumber string `json:"identificationNumber"`
	Bp                   string `json:"bp"`
	StartDate            string `json:"start_date"`
	EndDate              string `json:"end_date"`
	Active               bool   `json:"active"`
	Role                 string `json:"role"`
}

type EmployeeResponse struct {
	ID                   string    `json:"id"`
	Name                 string    `json:"name"`
	Airline              string    `json:"airline,omitempty"`
	Email                string    `json:"email"`
	IdentificationNumber string    `json:"identification_number"`
	Bp                   string    `json:"bp,omitempty"`
	StartDate            time.Time `json:"start_date"`
	EndDate              time.Time `json:"end_date"`
	Active               bool      `json:"active"`
	Role                 string    `json:"role,omitempty"`
}

type RegisterEmployeeResponse struct {
	Message string `json:"message"`
	Links   []Link `json:"_links"`
}

type ResponseEmail struct {
	Title   string
	Content template.HTML
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

type ResendVerificationEmailRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ResendVerificationEmailResponse struct {
	Message string          `json:"message"`
	Links   map[string]Link `json:"_links"`
}

type PasswordResetRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type PasswordResetResponse struct {
	Message string          `json:"message"`
	Links   map[string]Link `json:"_links"`
}

func (e EmployeeRequest) ToDomain() domain.Employee {
	layout := "2006-01-02"

	startDate, err := time.Parse(layout, e.StartDate)
	if err != nil {
		return domain.Employee{}
	}

	var endDate time.Time
	if e.EndDate != "" {
		endDate, err = time.Parse(layout, e.EndDate)
		if err != nil {
			return domain.Employee{}
		}
	}

	return domain.Employee{
		Name:                 e.Name,
		Airline:              e.Airline,
		Email:                e.Email,
		Password:             e.Password,
		IdentificationNumber: e.IdentificationNumber,
		Bp:                   e.Bp,
		StartDate:            startDate,
		EndDate:              endDate,
		Active:               e.Active,
		Role:                 e.Role,
	}
}
