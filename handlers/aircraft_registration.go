package handlers

import (
	domain "github.com/champion19/flighthours-api/core/interactor/services/domain"
	"github.com/google/uuid"
)

// AircraftRegistrationResponse - Response DTO for aircraft registration data
type AircraftRegistrationResponse struct {
	ID              string `json:"id"`
	UUID            string `json:"uuid"`
	LicensePlate    string `json:"license_plate"` // Numero de Matrícula
	ModelName       string `json:"model_name"`    // Modelo (denormalized)
	AirlineName     string `json:"airline_name"`  // Aerolínea (denormalized)
	AircraftModelID string `json:"aircraft_model_id"`
	AirlineID       string `json:"airline_id"`
	Links           []Link `json:"_links,omitempty"`
}

// FromDomainAircraftRegistration converts domain.AircraftRegistration to AircraftRegistrationResponse with encoded ID
func FromDomainAircraftRegistration(reg *domain.AircraftRegistration, encodedID string) AircraftRegistrationResponse {
	return AircraftRegistrationResponse{
		ID:              encodedID,
		UUID:            reg.ID,
		LicensePlate:    reg.LicensePlate,
		ModelName:       reg.ModelName,
		AirlineName:     reg.AirlineName,
		AircraftModelID: reg.AircraftModelID,
		AirlineID:       reg.AirlineID,
	}
}

// CreateAircraftRegistrationRequest - Request DTO for creating aircraft registration (HU34)
type CreateAircraftRegistrationRequest struct {
	LicensePlate    string `json:"license_plate" binding:"required"`
	AircraftModelID string `json:"aircraft_model_id" binding:"required"`
	AirlineID       string `json:"airline_id" binding:"required"`
}

// Sanitize trims whitespace from CreateAircraftRegistrationRequest fields
func (r *CreateAircraftRegistrationRequest) Sanitize() {
	r.LicensePlate = TrimString(r.LicensePlate)
	r.AircraftModelID = TrimString(r.AircraftModelID)
	r.AirlineID = TrimString(r.AirlineID)
}

// ToDomain converts CreateAircraftRegistrationRequest to domain.AircraftRegistration
func (r *CreateAircraftRegistrationRequest) ToDomain() domain.AircraftRegistration {
	return domain.AircraftRegistration{
		ID:              uuid.New().String(),
		LicensePlate:    r.LicensePlate,
		AircraftModelID: r.AircraftModelID,
		AirlineID:       r.AirlineID,
	}
}

// UpdateAircraftRegistrationRequest - Request DTO for updating aircraft registration (HU35)
type UpdateAircraftRegistrationRequest struct {
	LicensePlate    string `json:"license_plate" binding:"required"`
	AircraftModelID string `json:"aircraft_model_id" binding:"required"`
	AirlineID       string `json:"airline_id" binding:"required"`
}

// Sanitize trims whitespace from UpdateAircraftRegistrationRequest fields
func (r *UpdateAircraftRegistrationRequest) Sanitize() {
	r.LicensePlate = TrimString(r.LicensePlate)
	r.AircraftModelID = TrimString(r.AircraftModelID)
	r.AirlineID = TrimString(r.AirlineID)
}

// ToDomain converts UpdateAircraftRegistrationRequest to domain.AircraftRegistration
func (r *UpdateAircraftRegistrationRequest) ToDomain(id string) domain.AircraftRegistration {
	return domain.AircraftRegistration{
		ID:              id,
		LicensePlate:    r.LicensePlate,
		AircraftModelID: r.AircraftModelID,
		AirlineID:       r.AirlineID,
	}
}

// AircraftRegistrationListResponse - Response DTO for listing aircraft registrations
type AircraftRegistrationListResponse struct {
	Registrations []AircraftRegistrationResponse `json:"registrations"`
	Total         int                            `json:"total"`
	Links         []Link                         `json:"_links,omitempty"`
}

// ToAircraftRegistrationListResponse converts a slice of domain.AircraftRegistration to AircraftRegistrationListResponse
// baseURL is used to build HATEOAS links for each registration
func ToAircraftRegistrationListResponse(registrations []domain.AircraftRegistration, encodeFunc func(string) (string, error), baseURL string) AircraftRegistrationListResponse {
	response := AircraftRegistrationListResponse{
		Registrations: make([]AircraftRegistrationResponse, 0, len(registrations)),
		Total:         len(registrations),
	}

	for _, reg := range registrations {
		encodedID, err := encodeFunc(reg.ID)
		if err != nil {
			// If encoding fails, use the original UUID
			encodedID = reg.ID
		}
		regResp := AircraftRegistrationResponse{
			ID:              encodedID,
			UUID:            reg.ID,
			LicensePlate:    reg.LicensePlate,
			ModelName:       reg.ModelName,
			AirlineName:     reg.AirlineName,
			AircraftModelID: reg.AircraftModelID,
			AirlineID:       reg.AirlineID,
		}
		// Add HATEOAS links to each registration
		if baseURL != "" {
			regResp.Links = BuildAircraftRegistrationLinks(baseURL, encodedID)
		}
		response.Registrations = append(response.Registrations, regResp)
	}

	// Add collection-level links
	if baseURL != "" {
		response.Links = BuildAircraftRegistrationListLinks(baseURL)
	}

	return response
}
