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

// FromDomain convierte un domain.Employee a EmployeeResponse
// Esta función excluye explícitamente el campo Password para no exponerlo en la API
func FromDomain(employee *domain.Employee, encodedID string) EmployeeResponse {
	return EmployeeResponse{
		ID:                   encodedID,
		Name:                 employee.Name,
		Airline:              employee.Airline,
		Email:                employee.Email,
		IdentificationNumber: employee.IdentificationNumber,
		Bp:                   employee.Bp,
		StartDate:            employee.StartDate,
		EndDate:              employee.EndDate,
		Active:               employee.Active,
		Role:                 employee.Role,
	}
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

// ResendVerificationEmailResponse - Respuesta de reenvío de email de verificación
type ResendVerificationEmailResponse struct {
	Sent  bool   `json:"sent"`
	Email string `json:"email,omitempty"`
}

// PasswordResetRequest - DTO para solicitar recuperación de contraseña
type PasswordResetRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// PasswordResetResponse - Respuesta de solicitud de recuperación de contraseña
type PasswordResetResponse struct {
	Sent bool `json:"sent"`
}

// VerifyEmailRequest - DTO para verificar email mediante token proxy
// Este token es un JWT que contiene el email del usuario
type VerifyEmailRequest struct {
	Token string `json:"token" binding:"required"`
}

// VerifyEmailResponse - Respuesta de verificación de email
type VerifyEmailResponse struct {
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

// UpdateEmployeeRequest - DTO para actualizar información general del empleado
// Excluye email y password ya que se manejan en endpoints separados
type UpdateEmployeeRequest struct {
	Name                 string `json:"name"`
	Airline              string `json:"airline"`
	IdentificationNumber string `json:"identificationNumber"`
	Bp                   string `json:"bp"`
	StartDate            string `json:"start_date"`
	EndDate              string `json:"end_date"`
	Active               bool   `json:"active"`
	Role                 string `json:"role"`
}

// ToUpdateData convierte el request a un mapa de campos actualizables
// Se usa con un empleado existente para preservar email, password y keycloak_user_id
func (u UpdateEmployeeRequest) ToUpdateData(existing *domain.Employee) domain.Employee {
	layout := "2006-01-02"

	startDate := existing.StartDate
	if u.StartDate != "" {
		if parsed, err := time.Parse(layout, u.StartDate); err == nil {
			startDate = parsed
		}
	}

	endDate := existing.EndDate
	if u.EndDate != "" {
		if parsed, err := time.Parse(layout, u.EndDate); err == nil {
			endDate = parsed
		}
	}

	return domain.Employee{
		ID:                   existing.ID,
		Name:                 u.Name,
		Airline:              u.Airline,
		Email:                existing.Email,    // Preservar email
		Password:             existing.Password, // Preservar password (aunque no está en BD)
		IdentificationNumber: u.IdentificationNumber,
		Bp:                   u.Bp,
		StartDate:            startDate,
		EndDate:              endDate,
		Active:               u.Active,
		Role:                 u.Role,
		KeycloakUserID:       existing.KeycloakUserID, // Preservar keycloak_user_id
	}
}

// UpdateEmployeeResponse - Respuesta de actualización de empleado
type UpdateEmployeeResponse struct {
	ID      string `json:"id"`
	Updated bool   `json:"updated"`
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
