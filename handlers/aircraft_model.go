package handlers

import domain "github.com/champion19/flighthours-api/core/interactor/services/domain"

// AircraftModelResponse - Response DTO for aircraft model data
type AircraftModelResponse struct {
	ID               string `json:"id"`
	UUID             string `json:"uuid"`
	ModelName        string `json:"model_name"`
	AircraftTypeName string `json:"aircraft_type_name"`
	EngineTypeName   string `json:"engine_type_name,omitempty"`
	Family           string `json:"family"`
	Manufacturer     string `json:"manufacturer,omitempty"`
	Links            []Link `json:"_links,omitempty"`
}

// FromDomainAircraftModel converts domain.AircraftModel to AircraftModelResponse with encoded ID
func FromDomainAircraftModel(model *domain.AircraftModel, encodedID string) AircraftModelResponse {
	return AircraftModelResponse{
		ID:               encodedID,
		UUID:             model.ID,
		ModelName:        model.ModelName,
		AircraftTypeName: model.AircraftTypeName,
		EngineTypeName:   model.EngineTypeName,
		Family:           model.Family,
		Manufacturer:     model.Manufacturer,
	}
}

// AircraftModelListResponse - Response DTO for listing aircraft models
type AircraftModelListResponse struct {
	AircraftModels []AircraftModelResponse `json:"aircraft_models"`
	Total          int                     `json:"total"`
	Links          []Link                  `json:"_links,omitempty"`
}

// ToAircraftModelListResponse converts a slice of domain.AircraftModel to AircraftModelListResponse
// baseURL is used to build HATEOAS links for each model
func ToAircraftModelListResponse(models []domain.AircraftModel, encodeFunc func(string) (string, error), baseURL string) AircraftModelListResponse {
	response := AircraftModelListResponse{
		AircraftModels: make([]AircraftModelResponse, 0, len(models)),
		Total:          len(models),
	}

	for _, model := range models {
		encodedID, err := encodeFunc(model.ID)
		if err != nil {
			// If encoding fails, use the original UUID
			encodedID = model.ID
		}
		modelResp := AircraftModelResponse{
			ID:               encodedID,
			UUID:             model.ID,
			ModelName:        model.ModelName,
			AircraftTypeName: model.AircraftTypeName,
			EngineTypeName:   model.EngineTypeName,
			Family:           model.Family,
			Manufacturer:     model.Manufacturer,
		}
		// Add HATEOAS links to each model
		if baseURL != "" {
			modelResp.Links = BuildAircraftModelLinks(baseURL, encodedID)
		}
		response.AircraftModels = append(response.AircraftModels, modelResp)
	}

	// Add collection-level links
	if baseURL != "" {
		response.Links = BuildAircraftModelListLinks(baseURL)
	}

	return response
}
