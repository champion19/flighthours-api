package handlers

import (
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
	Role                 string    `json:"role"`
}

type RegisterEmployeeResponse struct {
	Links []Link `json:"_links"`
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

// ResendVerificationEmailRequest - DTO para reenviar email de verificación
type ResendVerificationEmailRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// PasswordResetRequest - DTO para solicitar recuperación de contraseña
type PasswordResetRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// VerifyEmailRequest - DTO para verificar email mediante token proxy
// Este token es un JWT que contiene el email del usuario
type verifyEmailRequest struct {
	Token string `json:"token" binding:"required"`
}

// VerifyEmailResponse - Respuesta de verificación de email
type verifyEmailResponse struct {
	Verified bool   `json:"verified"`
	Email    string `json:"email,omitempty"`
}

// UpdatePasswordRequest - DTO para actualizar contraseña con token de Keycloak
type UpdatePasswordRequest struct {
	Token           string `json:"token" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" binding:"required,min=8"`
}

// UpdatePasswordResponse - Respuesta de actualización de contraseña
type UpdatePasswordResponse struct {
	Updated bool   `json:"updated"`
	Email   string `json:"email,omitempty"`
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
