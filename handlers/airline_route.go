package handlers

import domain "github.com/champion19/flighthours-api/core/interactor/services/domain"

// AirlineRouteResponse - Response DTO for airline route data
type AirlineRouteResponse struct {
	ID                     string `json:"id"`
	RouteID                string `json:"route_id"`
	AirlineID              string `json:"airline_id"`
	Status                 bool   `json:"status"`
	AirlineCode            string `json:"airline_code"` // IATA code (e.g., "AV", "LA")
	AirlineName            string `json:"airline_name"`
	OriginIataCode         string `json:"origin_iata_code"`
	DestinationIataCode    string `json:"destination_iata_code"`
	RouteCode              string `json:"route_code"` // e.g., "BOG-CLO"
	OriginAirportName      string `json:"origin_airport_name,omitempty"`
	DestinationAirportName string `json:"destination_airport_name,omitempty"`
	AirportType            string `json:"airport_type,omitempty"`
	EstimatedFlightTime    string `json:"estimated_flight_time,omitempty"`
	Links                  []Link `json:"_links,omitempty"`
}

// FromDomainAirlineRoute converts domain.AirlineRoute to AirlineRouteResponse with encoded IDs
func FromDomainAirlineRoute(ar *domain.AirlineRoute, encodedID, encodedRouteID, encodedAirlineID string) AirlineRouteResponse {
	return AirlineRouteResponse{
		ID:                     encodedID,
		RouteID:                encodedRouteID,
		AirlineID:              encodedAirlineID,
		Status:                 ar.Status,
		AirlineCode:            ar.AirlineCode,
		AirlineName:            ar.AirlineName,
		OriginIataCode:         ar.OriginIataCode,
		DestinationIataCode:    ar.DestinationIataCode,
		RouteCode:              ar.RouteCode,
		OriginAirportName:      ar.OriginAirportName,
		DestinationAirportName: ar.DestinationAirportName,
		AirportType:            ar.AirportType,
		EstimatedFlightTime:    ar.EstimatedFlightTime,
	}
}

// AirlineRouteListResponse - Response DTO for listing airline routes
type AirlineRouteListResponse struct {
	AirlineRoutes []AirlineRouteResponse `json:"airline_routes"`
	Total         int                    `json:"total"`
	Links         []Link                 `json:"_links,omitempty"`
}

// ToAirlineRouteListResponse converts a slice of domain.AirlineRoute to AirlineRouteListResponse
// baseURL is used to build HATEOAS links for each airline route
func ToAirlineRouteListResponse(airlineRoutes []domain.AirlineRoute, encodeFunc func(string) (string, error), baseURL string) AirlineRouteListResponse {
	response := AirlineRouteListResponse{
		AirlineRoutes: make([]AirlineRouteResponse, 0, len(airlineRoutes)),
		Total:         len(airlineRoutes),
	}

	for _, ar := range airlineRoutes {
		encodedID, err := encodeFunc(ar.ID)
		if err != nil {
			encodedID = ar.ID
		}
		encodedRouteID, err := encodeFunc(ar.RouteID)
		if err != nil {
			encodedRouteID = ar.RouteID
		}
		encodedAirlineID, err := encodeFunc(ar.AirlineID)
		if err != nil {
			encodedAirlineID = ar.AirlineID
		}
		arResp := AirlineRouteResponse{
			ID:                     encodedID,
			RouteID:                encodedRouteID,
			AirlineID:              encodedAirlineID,
			Status:                 ar.Status,
			AirlineCode:            ar.AirlineCode,
			AirlineName:            ar.AirlineName,
			OriginIataCode:         ar.OriginIataCode,
			DestinationIataCode:    ar.DestinationIataCode,
			RouteCode:              ar.RouteCode,
			OriginAirportName:      ar.OriginAirportName,
			DestinationAirportName: ar.DestinationAirportName,
			AirportType:            ar.AirportType,
			EstimatedFlightTime:    ar.EstimatedFlightTime,
		}
		// Add HATEOAS links to each airline route
		if baseURL != "" {
			arResp.Links = BuildAirlineRouteLinks(baseURL, encodedID, ar.Status)
		}
		response.AirlineRoutes = append(response.AirlineRoutes, arResp)
	}

	// Add collection-level links
	if baseURL != "" {
		response.Links = BuildAirlineRouteListLinks(baseURL)
	}

	return response
}
