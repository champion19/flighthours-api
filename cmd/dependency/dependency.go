package dependency

import (
	"github.com/champion19/flighthours-api/config"
	"github.com/champion19/flighthours-api/core/interactor"
	"github.com/champion19/flighthours-api/core/interactor/services"
	"github.com/champion19/flighthours-api/core/ports/input"
	"github.com/champion19/flighthours-api/core/ports/output"
	mysql "github.com/champion19/flighthours-api/platform/databases/mysql"
	repo "github.com/champion19/flighthours-api/platform/databases/repositories/employee"
	"github.com/champion19/flighthours-api/platform/identity_provider/keycloak"
	"github.com/champion19/flighthours-api/platform/logger"
)

type Dependencies struct {
	EmployeeService input.Service
	EmployeeRepo    output.Repository
	Interactor      *interactor.Interactor
	KeycloakClient  output.AuthClient
	Config          *config.Config
	Logger          logger.Logger
}

func Init() (*Dependencies, error) {
	log := logger.NewSlogLogger()
	log.Info("Starting application")
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Error("Failed to load config", err)
		return nil, err
	}
	log.Info("Config loaded successfully")

	db, err := mysql.GetDB(cfg.Database, log)
	if err != nil {
		log.Error("Failed to connect to database", err)
		return nil, err
	}
	log.Success("Database connection established successfully")

	keycloakClient, err := keycloak.NewClient(&cfg.Keycloak, log)
	if err != nil {
		log.Error("Failed to create keycloak client", err)
		return nil, err
	}
	log.Success("Keycloak client created successfully")

	employeeRepo, err := repo.NewClientRepository(db)
	if err != nil {
		log.Error("Failed to create employee repository", err)
		return nil, err
	}

	// Crear EmployeeService con todas las dependencias
	employeeService := services.NewService(employeeRepo, keycloakClient, log)

	interactorFacade := interactor.NewInteractor(employeeService, log)

	log.Success("Dependencies initialized successfully")

	return &Dependencies{
		EmployeeService: employeeService,
		EmployeeRepo:    employeeRepo,
		Interactor:      interactorFacade,
		KeycloakClient:  keycloakClient,
		Config:          cfg,
		Logger:          log,
	}, nil
}
