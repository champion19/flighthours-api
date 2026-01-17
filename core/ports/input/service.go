package input

import (
	"context"

	"github.com/Nerzal/gocloak/v13"
	"github.com/champion19/flighthours-api/core/interactor/dto"
	"github.com/champion19/flighthours-api/core/interactor/services/domain"
	"github.com/champion19/flighthours-api/core/ports/output"
)

type Service interface {
	BeginTx(ctx context.Context) (output.Tx, error)

	//employee-validaciones y consultas
	RegisterEmployee(ctx context.Context, employee domain.Employee) (*dto.RegisterEmployee, error)
	GetEmployeeByEmail(ctx context.Context, email string) (*domain.Employee, error)
	GetEmployeeByID(ctx context.Context, id string) (*domain.Employee, error)
	GetEmployeeByKeycloakID(ctx context.Context, keycloakUserID string) (*domain.Employee, error)
	LocateEmployee(ctx context.Context, id string) (*dto.RegisterEmployee, error)
	CheckAndCleanInconsistentState(ctx context.Context, email string) error

	//employee- operaciones transaccionales de BD
	SaveEmployeeToDB(ctx context.Context, tx output.Tx, employee domain.Employee) error
	UpdateEmployeeKeycloakID(ctx context.Context, tx output.Tx, employeeID string, keycloakUserID string) error
	// UpdateEmployee actualiza un empleado en la BD y sincroniza cambios relevantes con Keycloak
	// Sincroniza: estado active (enabled/disabled) y cambios de rol
	UpdateEmployee(ctx context.Context, employee domain.Employee, previousActive bool, previousRole string) error
	// DeleteEmployee elimina un empleado de la BD y de Keycloak
	// Primero elimina de Keycloak, luego de la BD
	DeleteEmployee(ctx context.Context, employeeID string, keycloakUserID string) error

	//employee-operaciones de keycloak
	CreateUserInKeycloak(ctx context.Context, employee *domain.Employee) (string, error)
	SetUserPassword(ctx context.Context, userID string, password string) error
	AssignUserRole(ctx context.Context, userID string, role string) error
	GetUserByEmail(ctx context.Context, email string) (*gocloak.User, error)
	SendVerificationEmail(ctx context.Context, userID string) error
	SendPasswordResetEmail(ctx context.Context, email string) error
	Login(ctx context.Context, email, password string) (*gocloak.JWT, error)
	VerifyEmailByToken(ctx context.Context, token string) (string, error)
	UpdatePassword(ctx context.Context, token, newPassword string) (string, error)
	// ChangePassword changes the password for an authenticated user who knows their current password
	// Returns the email of the user whose password was changed
	ChangePassword(ctx context.Context, email, currentPassword, newPassword string) (string, error)

	//employee- compensaciones (rollback)
	RollbackEmployee(ctx context.Context, employeeID string) error
	RollbackKeycloakUser(ctx context.Context, KeycloakUserID string) error
}
type MessageService interface {
	// Transacciones
	BeginTx(ctx context.Context) (output.Tx, error)

	// Messages - Validaciones y operaciones
	ValidateMessage(ctx context.Context, message domain.Message) error
	GetMessageByID(ctx context.Context, id string) (*domain.Message, error)
	GetMessageByCode(ctx context.Context, code string) (*domain.Message, error)
	ListMessages(ctx context.Context, filters map[string]interface{}) ([]domain.Message, error)
	ListActiveMessages(ctx context.Context) ([]domain.Message, error)

	// Messages - Operaciones transaccionales de BD
	SaveMessageToDB(ctx context.Context, tx output.Tx, message domain.Message) error
	UpdateMessageInDB(ctx context.Context, tx output.Tx, message domain.Message) error
	DeleteMessageFromDB(ctx context.Context, tx output.Tx, id string) error
}

// AirlineService defines the interface for airline business operations
type AirlineService interface {
	BeginTx(ctx context.Context) (output.Tx, error)

	// Airline - queries
	GetAirlineByID(ctx context.Context, id string) (*domain.Airline, error)
	ListAirlines(ctx context.Context, filters map[string]interface{}) ([]domain.Airline, error)

	// Airline - operations
	UpdateAirlineStatus(ctx context.Context, id string, status bool) error
	ActivateAirline(ctx context.Context, id string) error
	DeactivateAirline(ctx context.Context, id string) error
}

// AirportService defines the interface for airport business operations
type AirportService interface {
	BeginTx(ctx context.Context) (output.Tx, error)

	// Airport - queries
	GetAirportByID(ctx context.Context, id string) (*domain.Airport, error)
	ListAirports(ctx context.Context, filters map[string]interface{}) ([]domain.Airport, error)

	// Airport - operations
	UpdateAirportStatus(ctx context.Context, id string, status bool) error
	ActivateAirport(ctx context.Context, id string) error
	DeactivateAirport(ctx context.Context, id string) error
}

// DailyLogbookService defines the interface for daily logbook business operations
type DailyLogbookService interface {
	BeginTx(ctx context.Context) (output.Tx, error)

	// DailyLogbook - queries
	GetDailyLogbookByID(ctx context.Context, id string) (*domain.DailyLogbook, error)
	ListDailyLogbooksByEmployee(ctx context.Context, employeeID string, filters map[string]interface{}) ([]domain.DailyLogbook, error)

	// DailyLogbook - operations
	CreateDailyLogbook(ctx context.Context, logbook domain.DailyLogbook) error
	UpdateDailyLogbook(ctx context.Context, logbook domain.DailyLogbook) error
	DeleteDailyLogbook(ctx context.Context, id string) error
	ActivateDailyLogbook(ctx context.Context, id string) error
	DeactivateDailyLogbook(ctx context.Context, id string) error
}

// AircraftRegistrationService defines the interface for aircraft registration business operations
type AircraftRegistrationService interface {
	BeginTx(ctx context.Context) (output.Tx, error)

	// AircraftRegistration - queries
	GetAircraftRegistrationByID(ctx context.Context, id string) (*domain.AircraftRegistration, error)
	ListAircraftRegistrations(ctx context.Context, filters map[string]interface{}) ([]domain.AircraftRegistration, error)

	// AircraftRegistration - operations
	CreateAircraftRegistration(ctx context.Context, registration domain.AircraftRegistration) error
	UpdateAircraftRegistration(ctx context.Context, registration domain.AircraftRegistration) error
}

// AircraftModelService defines the interface for aircraft model business operations
type AircraftModelService interface {
	// AircraftModel - queries only (read-only module)
	GetAircraftModelByID(ctx context.Context, id string) (*domain.AircraftModel, error)
	ListAircraftModels(ctx context.Context, filters map[string]interface{}) ([]domain.AircraftModel, error)
	// GetAircraftModelsByFamily retrieves all aircraft models for a specific family (HU32)
	GetAircraftModelsByFamily(ctx context.Context, family string) ([]domain.AircraftModel, error)
}

// RouteService defines the interface for route business operations
type RouteService interface {
	// Route - queries only (read-only module for HU39)
	GetRouteByID(ctx context.Context, id string) (*domain.Route, error)
	ListRoutes(ctx context.Context, filters map[string]interface{}) ([]domain.Route, error)
}

// AirlineRouteService defines the interface for airline route business operations
type AirlineRouteService interface {
	BeginTx(ctx context.Context) (output.Tx, error)

	// AirlineRoute - queries
	GetAirlineRouteByID(ctx context.Context, id string) (*domain.AirlineRoute, error)
	ListAirlineRoutes(ctx context.Context, filters map[string]interface{}) ([]domain.AirlineRoute, error)

	// AirlineRoute - operations
	ActivateAirlineRoute(ctx context.Context, id string) error
	DeactivateAirlineRoute(ctx context.Context, id string) error
}

// DailyLogbookDetailService defines the interface for daily logbook detail business operations
// This is the CORE service for flight segment tracking
type DailyLogbookDetailService interface {
	BeginTx(ctx context.Context) (output.Tx, error)

	// DailyLogbookDetail - queries
	GetDailyLogbookDetailByID(ctx context.Context, id string) (*domain.DailyLogbookDetail, error)
	ListDailyLogbookDetailsByLogbook(ctx context.Context, logbookID string) ([]domain.DailyLogbookDetail, error)

	// DailyLogbookDetail - operations
	CreateDailyLogbookDetail(ctx context.Context, detail domain.DailyLogbookDetail) error
	UpdateDailyLogbookDetail(ctx context.Context, detail domain.DailyLogbookDetail) error
	DeleteDailyLogbookDetail(ctx context.Context, id string) error

	// DailyLogbookDetail - validations
	ValidateTimeSequence(outTime, takeoffTime, landingTime, inTime string) error
}

// EngineService defines the interface for engine business operations
type EngineService interface {
	// Engine - queries only (read-only catalog module)
	GetEngineByID(ctx context.Context, id string) (*domain.Engine, error)
	ListEngines(ctx context.Context) ([]domain.Engine, error)
}

// ManufacturerService defines the interface for manufacturer business operations
type ManufacturerService interface {
	// Manufacturer - queries only (read-only catalog module)
	GetManufacturerByID(ctx context.Context, id string) (*domain.Manufacturer, error)
	ListManufacturers(ctx context.Context) ([]domain.Manufacturer, error)
}
