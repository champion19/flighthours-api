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
	log.Info(logger.LogAppStarting)
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Error(logger.LogAppConfigError, err)
		return nil, err
	}

	db, err := mysql.GetDB(cfg.Database, log)
	if err != nil {
		log.Error(logger.LogAppDatabaseError, err)
		return nil, err
	}
	log.Success(logger.LogAppDatabaseConnected)

	keycloakClient, err := keycloak.NewClient(&cfg.Keycloak, log)
	if err != nil {
		log.Error(logger.LogKeycloakClientError, err)
		return nil, err
	}
	log.Success(logger.LogKeycloakClientCreated)

	employeeRepo, err := repo.NewClientRepository(db, keycloakClient)
	if err != nil {
		log.Error(logger.LogEmployeeRepoInitError, err)
		return nil, err
	}

	// Crear EmployeeService con todas las dependencias
	employeeService := services.NewService(employeeRepo, keycloakClient, log)

	interactorFacade := interactor.NewInteractor(employeeService, log)

	log.Success(logger.LogDepInitComplete)

	return &Dependencies{
		EmployeeService: employeeService,
		EmployeeRepo:    employeeRepo,
		Interactor:      interactorFacade,
		KeycloakClient:  keycloakClient,
		Config:          cfg,
		Logger:          log,
	}, nil
}
