package domain

// AircraftModel represents the aircraft model domain model
// Contains the model name, aircraft type and engine type
type AircraftModel struct {
	ID               string `json:"id"`
	ModelName        string `json:"model_name"`
	AircraftTypeName string `json:"aircraft_type_name"`
	EngineTypeName   string `json:"engine_type_name"`
}

// ToLogger returns a slice of strings for logging aircraft model information
func (am *AircraftModel) ToLogger() []string {
	return []string{
		"id:" + am.ID,
		"model_name:" + am.ModelName,
		"aircraft_type_name:" + am.AircraftTypeName,
		"engine_type_name:" + am.EngineTypeName,
	}
}
