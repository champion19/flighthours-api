package handlers

import domain "github.com/champion19/flighthours-api/core/interactor/services/domain"

// RouteResponse - Response DTO for route data
type RouteResponse struct {
	ID                     string `json:"id"`
	OriginAirportID        string `json:"origin_airport_id"`
	OriginIataCode         string `json:"origin_iata_code"`
	OriginAirportName      string `json:"origin_airport_name"`
	DestinationAirportID   string `json:"destination_airport_id"`
	DestinationIataCode    string `json:"destination_iata_code"`
	DestinationAirportName string `json:"destination_airport_name"`
	OriginCountry          string `json:"origin_country,omitempty"`
	DestinationCountry     string `json:"destination_country,omitempty"`
	AirportType            string `json:"airport_type"`
	EstimatedFlightTime    string `json:"estimated_flight_time,omitempty"` // Format: "HH:MM:SS"
	RouteCode              string `json:"route_code"`
	Links                  []Link `json:"_links,omitempty"`
}

// FromDomainRoute converts domain.Route to RouteResponse with encoded IDs
func FromDomainRoute(route *domain.Route, encodedID, encodedOriginAirportID, encodedDestAirportID string) RouteResponse {
	return RouteResponse{
		ID:                     encodedID,
		OriginAirportID:        encodedOriginAirportID,
		OriginIataCode:         route.OriginIataCode,
		OriginAirportName:      route.OriginAirportName,
		DestinationAirportID:   encodedDestAirportID,
		DestinationIataCode:    route.DestinationIataCode,
		DestinationAirportName: route.DestinationAirportName,
		OriginCountry:          route.OriginCountry,
		DestinationCountry:     route.DestinationCountry,
		AirportType:            route.AirportType,
		EstimatedFlightTime:    route.EstimatedFlightTime,
		RouteCode:              route.RouteCode,
	}
}

// RouteListResponse - Response DTO for listing routes
type RouteListResponse struct {
	Routes []RouteResponse `json:"routes"`
	Total  int             `json:"total"`
	Links  []Link          `json:"_links,omitempty"`
}

// ToRouteListResponse converts a slice of domain.Route to RouteListResponse
// baseURL is used to build HATEOAS links for each route
func ToRouteListResponse(routes []domain.Route, encodeFunc func(string) (string, error), baseURL string) RouteListResponse {
	response := RouteListResponse{
		Routes: make([]RouteResponse, 0, len(routes)),
		Total:  len(routes),
	}

	for _, route := range routes {
		encodedID, err := encodeFunc(route.ID)
		if err != nil {
			encodedID = route.ID
		}
		encodedOriginAirportID, err := encodeFunc(route.OriginAirportID)
		if err != nil {
			encodedOriginAirportID = route.OriginAirportID
		}
		encodedDestAirportID, err := encodeFunc(route.DestinationAirportID)
		if err != nil {
			encodedDestAirportID = route.DestinationAirportID
		}
		routeResp := RouteResponse{
			ID:                     encodedID,
			OriginAirportID:        encodedOriginAirportID,
			OriginIataCode:         route.OriginIataCode,
			OriginAirportName:      route.OriginAirportName,
			DestinationAirportID:   encodedDestAirportID,
			DestinationIataCode:    route.DestinationIataCode,
			DestinationAirportName: route.DestinationAirportName,
			OriginCountry:          route.OriginCountry,
			DestinationCountry:     route.DestinationCountry,
			AirportType:            route.AirportType,
			EstimatedFlightTime:    route.EstimatedFlightTime,
			RouteCode:              route.RouteCode,
		}
		// Add HATEOAS links to each route
		if baseURL != "" {
			routeResp.Links = BuildRouteLinks(baseURL, encodedID)
		}
		response.Routes = append(response.Routes, routeResp)
	}

	// Add collection-level links
	if baseURL != "" {
		response.Links = BuildRouteListLinks(baseURL)
	}

	return response
}
