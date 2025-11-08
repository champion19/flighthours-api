package dependency

import (
	"log/slog"

	"github.com/champion19/flighthours-api/config"
	"github.com/champion19/flighthours-api/core/interactor/services"
	"github.com/champion19/flighthours-api/platform/keycloak"
  "github.com/champion19/flighthours-api/core/ports/input"
  "github.com/champion19/flighthours-api/core/ports/output"
	mysql "github.com/champion19/flighthours-api/platform/mysql"
	repo "github.com/champion19/flighthours-api/repositories/employee"

)

type Dependencies struct {
	EmployeeService      input.Service
	EmployeeRepo         output.Repository
	KeycloakClient       output.AuthClient
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

	keycloakClient, err := keycloak.NewClient(&cfg.Keycloak)
	if err != nil {
		slog.Error("Failed to initialize Keycloak client", slog.String("error", err.Error()))
		return nil, err
	}
	slog.Info("Keycloak client initialized successfully")

	employeeRepo, err := repo.NewRepository(db,keycloakClient)
	if err != nil {
		slog.Error("Failed to create employee repository", slog.String("error", err.Error()))
		return nil, err
	}

	slog.Info("Initializing Keycloak client...",
		slog.String("server", cfg.Keycloak.ServerURL),
		slog.String("realm", cfg.Keycloak.Realm))
	keycloakClient, err = keycloak.NewClient(&cfg.Keycloak)
	if err != nil {
		slog.Error("Failed to initialize Keycloak client", slog.String("error", err.Error()))
		return nil, err
	}
	slog.Info("Keycloak client initialized successfully")


	// Crear EmployeeService con todas las dependencias
	employeeService := services.NewService(employeeRepo, keycloakClient, cfg)

	return &Dependencies{
		EmployeeService:      employeeService,
		EmployeeRepo:         employeeRepo,
		KeycloakClient:       keycloakClient,
		Config:               cfg,
	}, nil
}
