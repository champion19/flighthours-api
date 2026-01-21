package output

import (
	"context"

	"github.com/champion19/flighthours-api/core/interactor/services/domain"
)

type Tx interface {
	Commit() error
	Rollback() error
}

type Repository interface {
	BeginTx(ctx context.Context) (Tx, error)

	//employee operations - transactional
	Save(ctx context.Context, tx Tx, employee domain.Employee) error
	UpdateEmployee(ctx context.Context, tx Tx, employee domain.Employee) error
	PatchEmployee(ctx context.Context, tx Tx, id string, keycloakUserID string) error
	DeleteEmployee(ctx context.Context, tx Tx, id string) error

	//employee operations - read
	GetEmployeeByEmail(ctx context.Context, email string) (*domain.Employee, error)
	GetEmployeeByID(ctx context.Context, id string) (*domain.Employee, error)
	GetEmployeeByKeycloakID(ctx context.Context, keycloakUserID string) (*domain.Employee, error)
	// HU47 - Get employees by role (Virtual Entity pattern)
	GetEmployeesByRole(ctx context.Context, role string) ([]domain.Employee, error)
}

type MessageRepository interface {
	BeginTx(ctx context.Context) (Tx, error)

	// Message operations - transactional
	SaveMessage(ctx context.Context, tx Tx, message domain.Message) error
	UpdateMessage(ctx context.Context, tx Tx, message domain.Message) error
	DeleteMessage(ctx context.Context, tx Tx, id string) error

	// Message operations - read
	GetAllActive(ctx context.Context) ([]domain.Message, error)
	GetByID(ctx context.Context, id string) (*domain.Message, error)
	GetByCode(ctx context.Context, code string) (*domain.Message, error)
	GetByType(ctx context.Context, msgType string) ([]domain.Message, error)
	GetByModule(ctx context.Context, module string) ([]domain.Message, error)
}

// AirlineRepository defines the interface for airline data persistence
type AirlineRepository interface {
	BeginTx(ctx context.Context) (Tx, error)

	// Airline operations - read
	GetAirlineByID(ctx context.Context, id string) (*domain.Airline, error)
	ListAirlines(ctx context.Context, filters map[string]interface{}) ([]domain.Airline, error)

	// Airline operations - transactional
	UpdateAirlineStatus(ctx context.Context, tx Tx, id string, status bool) error
}

// AirportRepository defines the interface for airport data persistence
type AirportRepository interface {
	BeginTx(ctx context.Context) (Tx, error)

	// Airport operations - read
	GetAirportByID(ctx context.Context, id string) (*domain.Airport, error)
	ListAirports(ctx context.Context, filters map[string]interface{}) ([]domain.Airport, error)
	// GetAirportsByCity retrieves all airports for a specific city
	GetAirportsByCity(ctx context.Context, city string) ([]domain.Airport, error)
	// GetAirportsByCountry retrieves all airports for a specific country
	GetAirportsByCountry(ctx context.Context, country string) ([]domain.Airport, error)
	// GetAirportsByType retrieves all airports for a specific type
	GetAirportsByType(ctx context.Context, airportType string) ([]domain.Airport, error)

	// Airport operations - transactional
	UpdateAirportStatus(ctx context.Context, tx Tx, id string, status bool) error
}

// DailyLogbookRepository defines the interface for daily logbook data persistence
type DailyLogbookRepository interface {
	BeginTx(ctx context.Context) (Tx, error)

	// DailyLogbook operations - read
	GetDailyLogbookByID(ctx context.Context, id string) (*domain.DailyLogbook, error)
	ListDailyLogbooksByEmployee(ctx context.Context, employeeID string, filters map[string]interface{}) ([]domain.DailyLogbook, error)

	// DailyLogbook operations - transactional
	SaveDailyLogbook(ctx context.Context, tx Tx, logbook domain.DailyLogbook) error
	UpdateDailyLogbook(ctx context.Context, tx Tx, logbook domain.DailyLogbook) error
	DeleteDailyLogbook(ctx context.Context, tx Tx, id string) error
	UpdateDailyLogbookStatus(ctx context.Context, tx Tx, id string, status bool) error
}

// AircraftRegistrationRepository defines the interface for aircraft registration data persistence
type AircraftRegistrationRepository interface {
	BeginTx(ctx context.Context) (Tx, error)

	// AircraftRegistration operations - read
	GetAircraftRegistrationByID(ctx context.Context, id string) (*domain.AircraftRegistration, error)
	ListAircraftRegistrations(ctx context.Context, filters map[string]interface{}) ([]domain.AircraftRegistration, error)

	// AircraftRegistration operations - transactional
	SaveAircraftRegistration(ctx context.Context, tx Tx, registration domain.AircraftRegistration) error
	UpdateAircraftRegistration(ctx context.Context, tx Tx, registration domain.AircraftRegistration) error
}

// AircraftModelRepository defines the interface for aircraft model data persistence
type AircraftModelRepository interface {
	// AircraftModel operations - read only (no transactions needed)
	GetAircraftModelByID(ctx context.Context, id string) (*domain.AircraftModel, error)
	ListAircraftModels(ctx context.Context, filters map[string]interface{}) ([]domain.AircraftModel, error)
	// GetAircraftModelsByFamily retrieves all aircraft models for a specific family (HU32)
	GetAircraftModelsByFamily(ctx context.Context, family string) ([]domain.AircraftModel, error)
}

// RouteRepository defines the interface for route data persistence
type RouteRepository interface {
	// Route operations - read only (no transactions needed for HU39)
	GetRouteByID(ctx context.Context, id string) (*domain.Route, error)
	ListRoutes(ctx context.Context, filters map[string]interface{}) ([]domain.Route, error)
}

// AirlineRouteRepository defines the interface for airline route data persistence
type AirlineRouteRepository interface {
	BeginTx(ctx context.Context) (Tx, error)

	// AirlineRoute operations - read
	GetAirlineRouteByID(ctx context.Context, id string) (*domain.AirlineRoute, error)
	ListAirlineRoutes(ctx context.Context, filters map[string]interface{}) ([]domain.AirlineRoute, error)

	// AirlineRoute operations - transactional
	UpdateAirlineRouteStatus(ctx context.Context, tx Tx, id string, status bool) error
}

// DailyLogbookDetailRepository defines the interface for daily logbook detail data persistence
// This is the CORE repository for flight segment tracking
type DailyLogbookDetailRepository interface {
	BeginTx(ctx context.Context) (Tx, error)

	// DailyLogbookDetail operations - read
	GetDailyLogbookDetailByID(ctx context.Context, id string) (*domain.DailyLogbookDetail, error)
	ListDailyLogbookDetailsByLogbook(ctx context.Context, logbookID string) ([]domain.DailyLogbookDetail, error)

	// DailyLogbookDetail operations - transactional
	SaveDailyLogbookDetail(ctx context.Context, tx Tx, detail domain.DailyLogbookDetail) error
	UpdateDailyLogbookDetail(ctx context.Context, tx Tx, detail domain.DailyLogbookDetail) error
	DeleteDailyLogbookDetail(ctx context.Context, tx Tx, id string) error
}

// EngineRepository defines the interface for engine type data persistence
type EngineRepository interface {
	// Engine operations - read only (catalog table)
	GetEngineByID(ctx context.Context, id string) (*domain.Engine, error)
	ListEngines(ctx context.Context) ([]domain.Engine, error)
}

// ManufacturerRepository defines the interface for manufacturer data persistence
type ManufacturerRepository interface {
	// Manufacturer operations - read only (catalog table)
	GetManufacturerByID(ctx context.Context, id string) (*domain.Manufacturer, error)
	ListManufacturers(ctx context.Context) ([]domain.Manufacturer, error)
}

// AirlineEmployeeRepository defines the interface for airline employee data persistence (Release 15)
// This is a specialized view of employees where airline IS NOT NULL
type AirlineEmployeeRepository interface {
	BeginTx(ctx context.Context) (Tx, error)

	// AirlineEmployee operations - read (HU26)
	GetAirlineEmployeeByID(ctx context.Context, id string) (*domain.AirlineEmployee, error)
	ListAirlineEmployees(ctx context.Context, filters map[string]interface{}) ([]domain.AirlineEmployee, error)

	// AirlineEmployee operations - transactional
	SaveAirlineEmployee(ctx context.Context, tx Tx, employee domain.AirlineEmployee) error   // HU28
	UpdateAirlineEmployee(ctx context.Context, tx Tx, employee domain.AirlineEmployee) error // HU27
	UpdateAirlineEmployeeStatus(ctx context.Context, tx Tx, id string, status bool) error    // HU29, HU30
}
