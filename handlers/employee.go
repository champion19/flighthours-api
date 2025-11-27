package handlers

import (
	"html/template"
	"time"

	domain "github.com/champion19/flighthours-api/core/interactor/services/domain"
	"github.com/champion19/flighthours-api/platform/logger"
	"github.com/gin-gonic/gin"
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

// GetBaseURL extrae la URL base de la petición (scheme + host)
func GetBaseURL(c *gin.Context) string {
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	return scheme + "://" + c.Request.Host
}

// EncodeID ofusca un UUID usando el encoder del handler
// Retorna el ID ofuscado o un error si falla
func (h *handler) EncodeID(uuid string) (string, error) {
	encodedID, err := h.IDEncoder.Encode(uuid)
	if err != nil {
		h.Logger.Error(logger.Messages.IDErrorEncoding,
			"uuid", uuid,
			"error", err)
		return "", err
	}
	return encodedID, nil
}

// DecodeID desofusca un ID ofuscado a UUID usando el encoder del handler
// Retorna el UUID o un error si falla
func (h *handler) DecodeID(encodedID string) (string, error) {
	uuid, err := h.IDEncoder.Decode(encodedID)
	if err != nil {
		h.Logger.Error(logger.Messages.IDErrorDecoding,
			"encoded_id", encodedID,
			"error", err)
		return "", err
	}
	return uuid, nil
}

// HandleIDEncodingError maneja errores de ofuscamiento y envía respuesta apropiada
func (h *handler) HandleIDEncodingError(c *gin.Context, uuid string, err error) {
	h.Logger.Error(logger.Messages.IDErrorEncoding,
		"uuid", uuid,
		"error", err,
		"client_ip", c.ClientIP())
	c.Error(domain.ErrInvalidID)
}

// HandleIDDecodingError maneja errores de desofuscamiento y envía respuesta apropiada
func (h *handler) HandleIDDecodingError(c *gin.Context, encodedID string, err error) {
	h.Logger.Error(logger.Messages.IDErrorDecoding,
		"encoded_id", encodedID,
		"error", err,
		"client_ip", c.ClientIP())
	c.Error(domain.ErrInvalidID)
}

// SetLocationHeader establece el Location header con la URL del recurso
func SetLocationHeader(c *gin.Context, baseURL, resource, encodedID string) {
	locationURL := BuildResourceURL(baseURL, resource, encodedID)
	c.Header("Location", locationURL)
}
