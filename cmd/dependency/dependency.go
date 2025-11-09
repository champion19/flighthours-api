package dependency

import (
	"log/slog"

	"github.com/champion19/flighthours-api/config"
	"github.com/champion19/flighthours-api/core/interactor"
	"github.com/champion19/flighthours-api/core/interactor/services"
	"github.com/champion19/flighthours-api/core/ports/input"
	"github.com/champion19/flighthours-api/core/ports/output"
	mysql "github.com/champion19/flighthours-api/platform/databases/mysql"
	"github.com/champion19/flighthours-api/platform/identity_provider/keycloak"
	repo "github.com/champion19/flighthours-api/platform/databases/repositories/employee"
)

type Dependencies struct {
	EmployeeService      input.Service
	EmployeeRepo         output.Repository
	Interactor           *interactor.Interactor
	KeycloakClient       output.AuthClient
	Config               *config.Config
}

func Init() (*Dependencies, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	db, err := mysql.GetDB(cfg.Database)
	if err != nil {
		return nil, err
	}

	keycloakClient, err := keycloak.NewClient(&cfg.Keycloak)
	if err != nil {
		return nil, err
		}


	employeeRepo, err := repo.NewClientRepository(db,keycloakClient)
	if err != nil {
		slog.Error("Failed to create employee repository", slog.String("error", err.Error()))
		return nil, err
	}



	// Crear EmployeeService con todas las dependencias
	employeeService := services.NewService(employeeRepo,keycloakClient)

	interactorFacade := interactor.NewInteractor(employeeService)

	return &Dependencies{
		EmployeeService:      employeeService,
		EmployeeRepo:         employeeRepo,
		Interactor:           interactorFacade,
		KeycloakClient:       keycloakClient,
		Config:               cfg,
	}, nil
}
