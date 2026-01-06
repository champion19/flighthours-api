package handlers

import domain "github.com/champion19/flighthours-api/core/interactor/services/domain"

// AirlineResponse - Response DTO for airline data
type AirlineResponse struct {
	ID          string `json:"id"`
	UUID        string `json:"uuid"`
	AirlineName string `json:"airline_name"`
	AirlineCode string `json:"airline_code"`
	Status      string `json:"status"`
}

// FromDomainAirline converts domain.Airline to AirlineResponse with encoded ID
func FromDomainAirline(airline *domain.Airline, encodedID string) AirlineResponse {
	return AirlineResponse{
		ID:          encodedID,
		UUID:        airline.ID,
		AirlineName: airline.AirlineName,
		AirlineCode: airline.AirlineCode,
		Status:      airline.Status,
	}
}

// UpdateAirlineStatusRequest - Request DTO for updating airline status
type UpdateAirlineStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=active inactive"`
}

// AirlineStatusResponse - Response DTO for status update
type AirlineStatusResponse struct {
	ID      string `json:"id"`
	Status  string `json:"status"`
	Updated bool   `json:"updated"`
}

// AirlineListResponse - Response DTO for listing airlines
type AirlineListResponse struct {
	Airlines []AirlineResponse `json:"airlines"`
	Total    int               `json:"total"`
}

// ToAirlineListResponse converts a slice of domain.Airline to AirlineListResponse
func ToAirlineListResponse(airlines []domain.Airline, encodeFunc func(string) (string, error)) AirlineListResponse {
	response := AirlineListResponse{
		Airlines: make([]AirlineResponse, 0, len(airlines)),
		Total:    len(airlines),
	}

	for _, airline := range airlines {
		encodedID, err := encodeFunc(airline.ID)
		if err != nil {
			// If encoding fails, use the original UUID
			encodedID = airline.ID
		}
		response.Airlines = append(response.Airlines, AirlineResponse{
			ID:          encodedID,
			UUID:        airline.ID,
			AirlineName: airline.AirlineName,
			AirlineCode: airline.AirlineCode,
			Status:      airline.Status,
		})
	}

	return response
}
