package dependency

import (
	"log/slog"

	"github.com/champion19/Flighthours_backend/config"
	"github.com/champion19/Flighthours_backend/core/ports"
	"github.com/champion19/Flighthours_backend/core/services"
	keycloak "github.com/champion19/Flighthours_backend/platform/keycloak"
	mysql "github.com/champion19/Flighthours_backend/platform/mysql"
	repo "github.com/champion19/Flighthours_backend/repositories/employee"
)

type Dependencies struct {
	EmployeeService      ports.Service
	EmployeeRepo         ports.Repository
	KeycloakClient       ports.AuthClient
	AuthorizationService ports.AuthorizationService
	Config               *config.Config
}

func Init() (*Dependencies, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		slog.Error("Failed to load config", slog.String("error", err.Error()))
		return nil, err
	}

	slog.Info("Connecting to MySQL database...",
		slog.String("host", cfg.Database.Host),
		slog.String("port", cfg.Database.Port))
	db, err := mysql.GetDB(cfg.Database)
	if err != nil {
		slog.Error("Failed to connect to MySQL", slog.String("error", err.Error()))
		return nil, err
	}
	slog.Info("MySQL connection successful")

	employeeRepo, err := repo.NewRepository(db)
	if err != nil {
		slog.Error("Failed to create employee repository", slog.String("error", err.Error()))
		return nil, err
	}

	slog.Info("Initializing Keycloak client...",
		slog.String("server", cfg.Keycloak.ServerURL),
		slog.String("realm", cfg.Keycloak.Realm))
	keycloakClient, err := keycloak.NewClient(&cfg.Keycloak)
	if err != nil {
		slog.Error("Failed to initialize Keycloak client", slog.String("error", err.Error()))
		return nil, err
	}
	slog.Info("Keycloak client initialized successfully")

	// Crear AuthorizationService con Keycloak y Repository
	authorizationService := services.NewAuthorizationService(keycloakClient, employeeRepo)

	// Crear EmployeeService con todas las dependencias
	employeeService := services.NewService(employeeRepo, authorizationService, cfg)

	return &Dependencies{
		EmployeeService:      employeeService,
		EmployeeRepo:         employeeRepo,
		KeycloakClient:       keycloakClient,
		AuthorizationService: authorizationService,
		Config:               cfg,
	}, nil
}
