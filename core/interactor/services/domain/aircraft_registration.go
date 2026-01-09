package domain

// AircraftRegistration represents the aircraft registration domain model
// It links a license plate (matrícula) with an aircraft model and an airline
type AircraftRegistration struct {
	ID              string `json:"id"`
	LicensePlate    string `json:"license_plate"`
	AircraftModelID string `json:"aircraft_model_id"`
	AirlineID       string `json:"airline_id"`
	// Denormalized fields for display (populated from JOINs)
	ModelName   string `json:"model_name,omitempty"`
	AirlineName string `json:"airline_name,omitempty"`
}

// ToLogger returns a slice of strings for logging aircraft registration information
func (ar *AircraftRegistration) ToLogger() []string {
	return []string{
		"id:" + ar.ID,
		"license_plate:" + ar.LicensePlate,
		"aircraft_model_id:" + ar.AircraftModelID,
		"airline_id:" + ar.AirlineID,
	}
}
