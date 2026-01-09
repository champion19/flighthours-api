package aircraft_registration

import domain "github.com/champion19/flighthours-api/core/interactor/services/domain"

// AircraftRegistration is the database entity for aircraft_registration table
type AircraftRegistration struct {
	ID              string `db:"id"`
	LicensePlate    string `db:"license_plate"`
	AircraftModelID string `db:"aircraft_model_id"`
	AirlineID       string `db:"airline_id"`
	// Denormalized fields from JOINs
	ModelName   string `db:"model_name"`
	AirlineName string `db:"airline_name"`
}

// ToDomain converts the database entity to domain model
func (ar *AircraftRegistration) ToDomain() *domain.AircraftRegistration {
	return &domain.AircraftRegistration{
		ID:              ar.ID,
		LicensePlate:    ar.LicensePlate,
		AircraftModelID: ar.AircraftModelID,
		AirlineID:       ar.AirlineID,
		ModelName:       ar.ModelName,
		AirlineName:     ar.AirlineName,
	}
}

// FromDomain converts a domain model to database entity
func FromDomain(domainReg *domain.AircraftRegistration) *AircraftRegistration {
	return &AircraftRegistration{
		ID:              domainReg.ID,
		LicensePlate:    domainReg.LicensePlate,
		AircraftModelID: domainReg.AircraftModelID,
		AirlineID:       domainReg.AirlineID,
		ModelName:       domainReg.ModelName,
		AirlineName:     domainReg.AirlineName,
	}
}
