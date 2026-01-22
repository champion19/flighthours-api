package dependency

import (
	"context"
	"time"

	"github.com/champion19/flighthours-api/config"
	"github.com/champion19/flighthours-api/core/interactor"
	"github.com/champion19/flighthours-api/core/interactor/services"
	"github.com/champion19/flighthours-api/core/ports/input"
	"github.com/champion19/flighthours-api/core/ports/output"
	"github.com/champion19/flighthours-api/middleware"
	messagingCache "github.com/champion19/flighthours-api/platform/cache/messaging"
	mysql "github.com/champion19/flighthours-api/platform/databases/mysql"
	aircraftModelRepo "github.com/champion19/flighthours-api/platform/databases/repositories/aircraft_model"
	aircraftRegistrationRepo "github.com/champion19/flighthours-api/platform/databases/repositories/aircraft_registration"
	airlineRepo "github.com/champion19/flighthours-api/platform/databases/repositories/airline"
	airlineEmployeeRepo "github.com/champion19/flighthours-api/platform/databases/repositories/airline_employee"
	airlineRouteRepo "github.com/champion19/flighthours-api/platform/databases/repositories/airline_route"
	airportRepo "github.com/champion19/flighthours-api/platform/databases/repositories/airport"
	dailyLogbookRepo "github.com/champion19/flighthours-api/platform/databases/repositories/daily_logbook"
	dailyLogbookDetailRepo "github.com/champion19/flighthours-api/platform/databases/repositories/daily_logbook_detail"
	repo "github.com/champion19/flighthours-api/platform/databases/repositories/employee"
	engineRepo "github.com/champion19/flighthours-api/platform/databases/repositories/engine"
	manufacturerRepo "github.com/champion19/flighthours-api/platform/databases/repositories/manufacturer"
	messageRepo "github.com/champion19/flighthours-api/platform/databases/repositories/message"
	routeRepo "github.com/champion19/flighthours-api/platform/databases/repositories/route"
	"github.com/champion19/flighthours-api/platform/identity_provider/keycloak"
	"github.com/champion19/flighthours-api/platform/jwt"
	"github.com/champion19/flighthours-api/platform/logger"
	"github.com/champion19/flighthours-api/tools/idencoder"
)

type Dependencies struct {
	EmployeeService                input.Service
	EmployeeRepo                   output.Repository
	Interactor                     *interactor.Interactor
	KeycloakClient                 output.AuthClient
	Config                         *config.Config
	Logger                         logger.Logger
	IDEncoder                      *idencoder.HashidsEncoder
	ResponseHandler                *middleware.ResponseHandler
	MessagingCache                 *messagingCache.MessageCache
	MessageInteractor              *interactor.MessageInteractor
	AirlineInteractor              *interactor.AirlineInteractor
	AirportInteractor              *interactor.AirportInteractor
	DailyLogbookInteractor         *interactor.DailyLogbookInteractor
	AircraftRegistrationInteractor *interactor.AircraftRegistrationInteractor
	AircraftModelInteractor        *interactor.AircraftModelInteractor
	RouteInteractor                *interactor.RouteInteractor
	AirlineRouteInteractor         *interactor.AirlineRouteInteractor
	DailyLogbookDetailInteractor   *interactor.DailyLogbookDetailInteractor
	EngineInteractor               *interactor.EngineInteractor
	ManufacturerInteractor         *interactor.ManufacturerInteractor
	AirlineEmployeeInteractor      *interactor.AirlineEmployeeInteractor // Release 15
	JWTValidator                   *jwt.JWKSValidator
}

func Init() (*Dependencies, error) {
	log := logger.NewSlogLogger()
	log.Info(logger.LogAppStarting)

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Error(logger.LogAppConfigError, "error", err)
		return nil, err
	}
	log.Info(logger.LogAppConfigLoaded)

	//initialize metrics prometheus
	middleware.PrometheusInit()
	log.Success(logger.LogPrometheusInitOK)

	db, err := mysql.GetDB(cfg.Database, log)
	if err != nil {
		log.Error(logger.LogAppDatabaseError, "error", err)
		return nil, err
	}
	log.Success(logger.LogAppDatabaseConnected)

	keycloakClient, err := keycloak.NewClient(&cfg.Keycloak, log)
	if err != nil {
		log.Error(logger.LogKeycloakClientError, "error", err)
		return nil, err
	}
	log.Success(logger.LogKeycloakClientOK)

	employeeRepo, err := repo.NewClientRepository(db)
	if err != nil {
		log.Error(logger.LogEmployeeRepoInitError, "error", err)
		return nil, err
	}

	// Crear EmployeeService con todas las dependencias
	employeeService := services.NewService(employeeRepo, keycloakClient, log)

	interactorFacade := interactor.NewInteractor(employeeService, log)

	//initialize id encoder
	encoder, err := idencoder.NewHashidsEncoder(idencoder.Config{
		Secret:    cfg.IDEncoder.Secret,
		MinLength: cfg.IDEncoder.MinLength,
	}, log)
	if err != nil {
		log.Error(logger.LogIDEncoderInitError, "error", err)
		return nil, err
	}
	log.Success(logger.LogIDEncodeOK)

	// Inicializar repositorio de mensajes (implementa ambas interfaces)
	msgRepo, err := messageRepo.NewMessageRepository(db)
	if err != nil {
		log.Error(logger.LogRepoMsgInitError, "error", err)
		return nil, err
	}
	log.Success(logger.LogDependencyMessageRepoInit)

	refreshInterval := 5 * time.Minute
	messagingCache := messagingCache.NewMessageCache(msgRepo, refreshInterval)

	if err := messagingCache.LoadMessages(context.Background()); err != nil {
		log.Warn(logger.LogMsgCacheLoadError, "error", err)
		// Don't return error, continue with fallback
	}
	log.Success(logger.LogMsgCacheInit, "messages_loaded", messagingCache.MessageCount())

	// Iniciar auto-refresh en background
	messagingCache.StartAutoRefresh(context.Background())

	responseHandler := middleware.NewResponseHandler(messagingCache)

	// Inicializar servicio de mensajes (msgRepo también implementa output.MessageRepository)
	messageService := services.NewMessageService(msgRepo, log)
	messageInteractor := interactor.NewMessageInteractor(messageService, log)
	log.Success(logger.LogDependencyMessageIntInit)

	// Inicializar repositorio y servicio de aerolíneas
	airlineRepository, err := airlineRepo.NewAirlineRepository(db)
	if err != nil {
		log.Error(logger.LogAirlineRepoInitError, "error", err)
		return nil, err
	}
	log.Success(logger.LogAirlineRepoInitOK)

	airlineService := services.NewAirlineService(airlineRepository, log)
	airlineInteractor := interactor.NewAirlineInteractor(airlineService, log)

	// Inicializar repositorio y servicio de aeropuertos
	airportRepository, err := airportRepo.NewAirportRepository(db)
	if err != nil {
		log.Error(logger.LogAirportRepoInitError, "error", err)
		return nil, err
	}
	log.Success(logger.LogAirportRepoInitOK)

	airportService := services.NewAirportService(airportRepository, log)
	airportInteractor := interactor.NewAirportInteractor(airportService, log)

	// Inicializar repositorio y servicio de bitácoras diarias
	dailyLogbookRepository, err := dailyLogbookRepo.NewDailyLogbookRepository(db)
	if err != nil {
		log.Error(logger.LogDailyLogbookRepoInitError, "error", err)
		return nil, err
	}
	log.Success(logger.LogDailyLogbookRepoInitOK)

	dailyLogbookService := services.NewDailyLogbookService(dailyLogbookRepository, log)
	dailyLogbookInteractor := interactor.NewDailyLogbookInteractor(dailyLogbookService, log)

	// Inicializar repositorio y servicio de matrículas
	aircraftRegistrationRepository, err := aircraftRegistrationRepo.NewAircraftRegistrationRepository(db)
	if err != nil {
		log.Error(logger.LogAircraftRegistrationRepoInitError, "error", err)
		return nil, err
	}
	log.Success(logger.LogAircraftRegistrationRepoInitOK)

	aircraftRegistrationService := services.NewAircraftRegistrationService(aircraftRegistrationRepository, log)
	aircraftRegistrationInteractor := interactor.NewAircraftRegistrationInteractor(aircraftRegistrationService, log)

	// Inicializar repositorio y servicio de modelos de aeronave
	aircraftModelRepository, err := aircraftModelRepo.NewAircraftModelRepository(db)
	if err != nil {
		log.Error(logger.LogAircraftModelRepoInitError, "error", err)
		return nil, err
	}
	log.Success(logger.LogAircraftModelRepoInitOK)

	aircraftModelService := services.NewAircraftModelService(aircraftModelRepository, log)
	aircraftModelInteractor := interactor.NewAircraftModelInteractor(aircraftModelService, log)

	// Inicializar repositorio y servicio de rutas
	routeRepository, err := routeRepo.NewRouteRepository(db)
	if err != nil {
		log.Error(logger.LogRouteRepoInitError, "error", err)
		return nil, err
	}
	log.Success(logger.LogRouteRepoInitOK)

	routeService := services.NewRouteService(routeRepository, log)
	routeInteractor := interactor.NewRouteInteractor(routeService, log)

	// Inicializar repositorio y servicio de rutas aerolínea
	airlineRouteRepository, err := airlineRouteRepo.NewAirlineRouteRepository(db)
	if err != nil {
		log.Error(logger.LogAirlineRouteRepoInitError, "error", err)
		return nil, err
	}
	log.Success(logger.LogAirlineRouteRepoInitOK)

	airlineRouteService := services.NewAirlineRouteService(airlineRouteRepository)
	airlineRouteInteractor := interactor.NewAirlineRouteInteractor(airlineRouteService)

	// Inicializar repositorio y servicio de detalles de bitácora diaria
	dailyLogbookDetailRepository, err := dailyLogbookDetailRepo.NewDailyLogbookDetailRepository(db)
	if err != nil {
		log.Error(logger.LogDailyLogbookDetailRepoInitError, "error", err)
		return nil, err
	}
	log.Success(logger.LogDailyLogbookDetailRepoInitOK)

	dailyLogbookDetailService := services.NewDailyLogbookDetailService(dailyLogbookDetailRepository)
	dailyLogbookDetailInteractor := interactor.NewDailyLogbookDetailInteractor(dailyLogbookDetailService, dailyLogbookService)

	// Inicializar repositorio y servicio de motores (Engine)
	engineRepository, err := engineRepo.NewEngineRepository(db)
	if err != nil {
		log.Error(logger.LogEngineRepoInitError, "error", err)
		return nil, err
	}
	log.Success(logger.LogEngineRepoInitOK)

	engineService := services.NewEngineService(engineRepository, log)
	engineInteractor := interactor.NewEngineInteractor(engineService, log)

	// Inicializar repositorio y servicio de fabricantes (Manufacturer)
	manufacturerRepository, err := manufacturerRepo.NewManufacturerRepository(db)
	if err != nil {
		log.Error(logger.LogManufacturerRepoInitError, "error", err)
		return nil, err
	}
	log.Success(logger.LogManufacturerRepoInitOK)

	manufacturerService := services.NewManufacturerService(manufacturerRepository, log)
	manufacturerInteractor := interactor.NewManufacturerInteractor(manufacturerService, log)

	// Inicializar repositorio y servicio de empleados aerolínea (Release 15)
	airlineEmployeeRepository, err := airlineEmployeeRepo.NewAirlineEmployeeRepository(db)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, "error", err, "repository", "airline_employee")
		return nil, err
	}
	log.Success(logger.LogDatabaseAvailable, "repository", "airline_employee")

	airlineEmployeeService := services.NewAirlineEmployeeService(airlineEmployeeRepository, log)
	airlineEmployeeInteractor := interactor.NewAirlineEmployeeInteractor(airlineEmployeeService, log)

	// JWKS Validator (JWT signature and expiration validation)
	// This fetches Keycloak's public keys for local token validation
	var jwtValidator *jwt.JWKSValidator
	jwtConfig := jwt.JWKSConfig{
		JWKSURL:         cfg.GetKeycloakJWKSURL(),
		Issuer:          cfg.GetKeycloakIssuerURL(),
		RefreshInterval: 15 * time.Minute, // Refresh keys every 15 minutes
	}
	jwtValidator, err = jwt.NewJWKSValidator(context.Background(), jwtConfig)
	if err != nil {
		log.Warn("JWKS validator initialization failed, using fallback validation", "error", err)
		// Don't fail startup - middleware will fall back to simple parsing
		// This allows the app to start even if Keycloak is temporarily unavailable
		jwtValidator = nil
	} else {
		log.Success("JWKS validator initialized", "jwks_url", jwtConfig.JWKSURL)
	}

	return &Dependencies{
		EmployeeService:                employeeService,
		EmployeeRepo:                   employeeRepo,
		Interactor:                     interactorFacade,
		KeycloakClient:                 keycloakClient,
		Config:                         cfg,
		Logger:                         log,
		IDEncoder:                      encoder,
		ResponseHandler:                responseHandler,
		MessagingCache:                 messagingCache,
		MessageInteractor:              messageInteractor,
		AirlineInteractor:              airlineInteractor,
		AirportInteractor:              airportInteractor,
		DailyLogbookInteractor:         dailyLogbookInteractor,
		AircraftRegistrationInteractor: aircraftRegistrationInteractor,
		AircraftModelInteractor:        aircraftModelInteractor,
		RouteInteractor:                routeInteractor,
		AirlineRouteInteractor:         airlineRouteInteractor,
		DailyLogbookDetailInteractor:   dailyLogbookDetailInteractor,
		EngineInteractor:               engineInteractor,
		ManufacturerInteractor:         manufacturerInteractor,
		AirlineEmployeeInteractor:      airlineEmployeeInteractor,
		JWTValidator:                   jwtValidator,
	}, nil
}
