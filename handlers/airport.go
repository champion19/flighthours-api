package handlers

import domain "github.com/champion19/flighthours-api/core/interactor/services/domain"

// AirportResponse - Response DTO for airport data
type AirportResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	City        string `json:"city,omitempty"`
	Country     string `json:"country,omitempty"`
	IATACode    string `json:"iata_code,omitempty"`
	Status      string `json:"status"`
	AirportType string `json:"airport_type,omitempty"`
}

// FromDomainAirport converts domain.Airport to AirportResponse with encoded ID
func FromDomainAirport(airport *domain.Airport, encodedID string) AirportResponse {
	status := "inactive"
	if airport.Status {
		status = "active"
	}
	return AirportResponse{
		ID:          encodedID,
		Name:        airport.Name,
		City:        airport.City,
		Country:     airport.Country,
		IATACode:    airport.IATACode,
		Status:      status,
		AirportType: airport.AirportType,
	}
}

// AirportStatusResponse - Response DTO for status update
type AirportStatusResponse struct {
	ID      string `json:"id"`
	Status  string `json:"status"`
	Updated bool   `json:"updated"`
}

// AirportListResponse - Response DTO for listing airports
type AirportListResponse struct {
	Airports []AirportResponse `json:"airports"`
	Total    int               `json:"total"`
}

// ToAirportListResponse converts a slice of domain airports to AirportListResponse
func ToAirportListResponse(airports []domain.Airport, encodeFunc func(string) (string, error)) AirportListResponse {
	response := AirportListResponse{
		Airports: make([]AirportResponse, 0, len(airports)),
		Total:    len(airports),
	}

	for _, airport := range airports {
		encodedID, err := encodeFunc(airport.ID)
		if err != nil {
			// If encoding fails, use the original UUID
			encodedID = airport.ID
		}
		status := "inactive"
		if airport.Status {
			status = "active"
		}
		response.Airports = append(response.Airports, AirportResponse{
			ID:          encodedID,
			Name:        airport.Name,
			City:        airport.City,
			Country:     airport.Country,
			IATACode:    airport.IATACode,
			Status:      status,
			AirportType: airport.AirportType,
		})
	}

	return response
}
