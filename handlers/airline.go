package handlers

import domain "github.com/champion19/flighthours-api/core/interactor/services/domain"

// AirlineResponse - Response DTO for airline data
type AirlineResponse struct {
	ID          string `json:"id"`
	AirlineName string `json:"airline_name"`
	AirlineCode string `json:"airline_code"`
	Status      string `json:"status"`
}

// FromDomainAirline converts domain.Airline to AirlineResponse with encoded ID
func FromDomainAirline(airline *domain.Airline, encodedID string) AirlineResponse {
	return AirlineResponse{
		ID:          encodedID,
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
