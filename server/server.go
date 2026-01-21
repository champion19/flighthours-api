package server

import (
	"log/slog"
	"time"

	"github.com/champion19/flighthours-api/cmd/dependency"
	"github.com/champion19/flighthours-api/handlers"
	"github.com/champion19/flighthours-api/middleware"
	"github.com/champion19/flighthours-api/platform/logger"
	"github.com/champion19/flighthours-api/platform/schema"
	_ "github.com/champion19/flighthours-api/platform/swaggo"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func routing(app *gin.Engine, dependencies *dependency.Dependencies) {
	dependencies.Logger.Info(logger.LogRouteConfiguring)

	// CORS configuration - Allow requests from Keycloak (localhost:8080) and other origins
	// This is required for the email verification flow from Keycloak's theme pages
	corsConfig := cors.Config{
		AllowOrigins:     []string{"http://localhost:8080", "http://localhost:8081", "http://localhost:3001"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Request-ID"},
		ExposeHeaders:    []string{"Content-Length", "Location"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	app.Use(cors.New(corsConfig))
	dependencies.Logger.Info("CORS middleware configured")

	// Endpoint de métricas de Prometheus
	app.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Documentación Swagger
	app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Middleware de Prometheus para capturar métricas
	app.Use(middleware.RequestID())

	// Apply Prometheus metrics tracking middleware
	app.Use(middleware.TrackMetrics())

	errorHandler := middleware.NewErrorHandler(dependencies.MessagingCache)
	app.Use(errorHandler.Handle())

	handler := handlers.New(
		dependencies.EmployeeService,
		dependencies.Interactor,
		dependencies.IDEncoder,
		dependencies.ResponseHandler,
		dependencies.MessageInteractor,
		dependencies.MessagingCache,
		dependencies.AirlineInteractor,
		dependencies.AirportInteractor,
		dependencies.DailyLogbookInteractor,
		dependencies.AircraftRegistrationInteractor,
		dependencies.AircraftModelInteractor,
		dependencies.RouteInteractor,
		dependencies.AirlineRouteInteractor,
		dependencies.DailyLogbookDetailInteractor,
		dependencies.EngineInteractor,
		dependencies.ManufacturerInteractor,
		dependencies.AirlineEmployeeInteractor,
	)

	validators, err := schema.NewValidator(&schema.DefaultFileReader{})
	if err != nil {
		dependencies.Logger.Error(logger.LogRouteValidatorError, err)
		dependencies.Logger.Fatal(logger.LogRouteValidatorError, err)
		return
	}
	dependencies.Logger.Success(logger.LogRouteValidatorOK)
	validator := middleware.NewMiddlewareValidator(validators)

	// Health check endpoint (no authentication required)
	app.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "flighthours-backend",
		})
	})

	app.NoRoute(middleware.NotFoundHandler())

	// ===========================================
	// PUBLIC ROUTES (no authentication required)
	// ===========================================
	public := app.Group("flighthours/api/v1")
	{
		// ---- Authentication ----
		// POST /register - New user registration
		public.POST("/register", validator.WithValidateRegister(), handler.RegisterEmployee())

		// POST /login - Returns JWT tokens
		public.POST("/login", handler.Login())

		// ---- Email Verification & Password Reset (Public) ----
		// POST /auth/resend-verification - Resend verification email
		public.POST("/auth/resend-verification", validator.WithValidateResendVerificationEmail(), handler.ResendVerificationEmail())

		// POST /auth/password-reset - Request password reset
		public.POST("/auth/password-reset", validator.WithValidatePasswordResetRequest(), handler.RequestPasswordReset())

		// POST /auth/verify-email - Verify email with token
		public.POST("/auth/verify-email", handler.VerifyEmailByToken())

		// POST /auth/update-password - Update password with token (from reset email)
		public.POST("/auth/update-password", validator.WithValidateUpdatePassword(), handler.UpdatePassword())

		// ---- Airlines (Public - Read Only) ----
		// GET /airlines - List all airlines (with optional status filter)
		// Query params: ?status=true (active) or ?status=false (inactive)
		public.GET("/airlines", handler.ListAirlines())

		// GET /airlines/:id - Get airline information by ID
		public.GET("/airlines/:id", handler.GetAirlineByID())

		// ---- Airports (Public - Read Only) ----
		// GET /airports - List all airports (with optional status filter)
		// Query params: ?status=true (active) or ?status=false (inactive)
		public.GET("/airports", handler.ListAirports())

		// GET /airports/:id - Get airport information by ID
		public.GET("/airports/:id", handler.GetAirportByID())

		// ---- Aircraft Models (Public - Read Only) ----
		// GET /aircraft-models - List all aircraft models (HU43)
		// Query params: ?engine_type=JET (filter by engine type)
		public.GET("/aircraft-models", handler.ListAircraftModels())

		// GET /aircraft-models/:id - Get aircraft model by ID (HU36)
		public.GET("/aircraft-models/:id", handler.GetAircraftModelByID())

		// ---- Aircraft Families (Public - Read Only) - HU32 ----
		// GET /aircraft-families/:family - Get all models belonging to a family
		// Example: GET /aircraft-families/A320 returns A320, A321, A319 models
		public.GET("/aircraft-families/:family", handler.GetAircraftModelsByFamily())

		// ---- Routes (Public - Read Only) ----
		// GET /routes - List all routes (HU36)
		// Query params: ?airport_type=Nacional, ?origin_country=Colombia
		public.GET("/routes", handler.ListRoutes())

		// GET /routes/:id - Get route by ID (HU36)
		public.GET("/routes/:id", handler.GetRouteByID())

		// ---- Airline Routes (Public - Read Only) ----
		// GET /airline-routes - List all airline routes
		// Query params: ?airline_code=AV, ?status=true
		public.GET("/airline-routes", handler.ListAirlineRoutes())

		// GET /airline-routes/:id - Get airline route by ID (HU37)
		public.GET("/airline-routes/:id", handler.GetAirlineRoute())

		// ---- Engines (Public - Read Only) - HU35 ----
		// GET /engines - List all engine types
		public.GET("/engines", handler.ListEngines())

		// GET /engines/:id - Get engine by ID
		public.GET("/engines/:id", handler.GetEngineByID())

		// ---- Manufacturers (Public - Read Only) - HU29 ----
		// GET /manufacturers - List all manufacturers
		public.GET("/manufacturers", handler.ListManufacturers())

		// GET /manufacturers/:id - Get manufacturer by ID
		public.GET("/manufacturers/:id", handler.GetManufacturerByID())

		// ---- Cities (Public - Read Only) ----
		// GET /cities/:city_name - Get all airports in a city
		// Example: GET /cities/Bogota returns all airports in Bogota
		// Note: No new tables needed - queries airport.city field directly
		public.GET("/cities/:city_name", handler.GetAirportsByCity())

		// ---- Countries (Public - Read Only) ----
		// GET /countries/:country_name - Get all airports in a country
		// Example: GET /countries/Colombia returns all airports in Colombia
		// Note: No new tables needed - queries airport.country field directly
		public.GET("/countries/:country_name", handler.GetAirportsByCountry())

		// ---- Airport Types (Public - Read Only) ----
		// GET /airport-types/:airport_type - Get all airports of a specific type
		// Example: GET /airport-types/INTERNACIONAL returns all international airports
		// Note: No new tables needed - queries airport.airport_type field directly
		public.GET("/airport-types/:airport_type", handler.GetAirportsByType())

		// ---- Crew Member Types (Public - Read Only) ----
		// GET /crew-member-types/:role - Get all employees of a specific role
		// Example: GET /crew-member-types/PILOTO returns all pilots
		// Note: No new tables needed - queries employee.role field directly
		public.GET("/crew-member-types/:role", handler.GetEmployeesByRole())
	}

	// ===========================================
	// PROTECTED ROUTES (authentication required)
	// ===========================================
	protected := app.Group("flighthours/api/v1")
	// Use the RequireAuth middleware from jwt_middleware.go
	// This validates JWT tokens and injects the authenticated user into context
	protected.Use(middleware.RequireAuth(dependencies.EmployeeService, dependencies.MessagingCache, dependencies.JWTValidator))
	{
		// ---- Authenticated User Endpoints ----
		// POST /auth/change-password - Change password (authenticated user knows current password)
		protected.POST("/auth/change-password", validator.WithValidateChangePassword(), handler.ChangePassword())

		// ---- Employee Self-Service (Protected) ----
		// GET /employees/me - Get current authenticated employee's information
		protected.GET("/employees/me", handler.GetMe())

		// PUT /employees/me - Update current authenticated employee's information
		// Does not modify email or password (handled in separate endpoints)
		// Syncs active and role changes with Keycloak
		protected.PUT("/employees/me", validator.WithValidateUpdateEmployee(), handler.UpdateMe())

		// DELETE /employees/me - Delete current authenticated employee (DB and Keycloak)
		// This operation is irreversible
		protected.DELETE("/employees/me", handler.DeleteMe())

		// ---- Messages Management (Protected - Admin) ----
		// POST /messages - Create new message
		protected.POST("/messages", validator.WithValidateMessage(), handler.CreateMessage())

		// PUT /messages/:id - Update existing message
		protected.PUT("/messages/:id", validator.WithValidateMessage(), handler.UpdateMessage())

		// DELETE /messages/:id - Delete message
		protected.DELETE("/messages/:id", handler.DeleteMessage())

		// GET /messages/:id - Get message by ID
		protected.GET("/messages/:id", handler.GetMessageByID())

		// GET /messages - List messages (with optional filters)
		// Query params: ?module=users&type=ERROR&category=usuario_final&active=true
		protected.GET("/messages", handler.ListMessages())

		// POST /messages/cache/reload - Reload message cache from DB
		// Administrative endpoint to force reload after manual changes
		protected.POST("/messages/cache/reload", handler.ReloadMessageCache())

		// ---- Airlines Management (Protected - Write Operations) ----
		// PATCH /airlines/:id/activate - Activate an airline
		protected.PATCH("/airlines/:id/activate", handler.ActivateAirline())

		// PATCH /airlines/:id/deactivate - Deactivate an airline
		protected.PATCH("/airlines/:id/deactivate", handler.DeactivateAirline())

		// ---- Airports Management (Protected - Write Operations) ----
		// PATCH /airports/:id/activate - Activate an airport
		protected.PATCH("/airports/:id/activate", handler.ActivateAirport())

		// PATCH /airports/:id/deactivate - Deactivate an airport
		protected.PATCH("/airports/:id/deactivate", handler.DeactivateAirport())

		// ---- Daily Logbooks Management (Protected) ----
		// GET /daily-logbooks - List daily logbooks for authenticated employee
		// Query params: ?status=true (active) or ?status=false (inactive)
		protected.GET("/daily-logbooks", handler.ListDailyLogbooks())

		// GET /daily-logbooks/:id - Get a specific daily logbook by ID
		protected.GET("/daily-logbooks/:id", handler.GetDailyLogbookByID())

		// POST /daily-logbooks - Create a new daily logbook
		protected.POST("/daily-logbooks", handler.CreateDailyLogbook())

		// PUT /daily-logbooks/:id - Update an existing daily logbook
		protected.PUT("/daily-logbooks/:id", handler.UpdateDailyLogbook())

		// DELETE /daily-logbooks/:id - Delete a daily logbook
		protected.DELETE("/daily-logbooks/:id", handler.DeleteDailyLogbook())

		// PATCH /daily-logbooks/:id/activate - Activate a daily logbook
		protected.PATCH("/daily-logbooks/:id/activate", handler.ActivateDailyLogbook())

		// PATCH /daily-logbooks/:id/deactivate - Deactivate a daily logbook
		protected.PATCH("/daily-logbooks/:id/deactivate", handler.DeactivateDailyLogbook())

		// ---- Aircraft Registrations Management (Protected) ----
		// GET /aircraft-registrations - List all aircraft registrations
		// Query params: ?airline_id=xxx (filter by airline)
		protected.GET("/aircraft-registrations", handler.ListAircraftRegistrations())

		// GET /aircraft-registrations/:id - Get aircraft registration by ID (HU33)
		protected.GET("/aircraft-registrations/:id", handler.GetAircraftRegistrationByID())

		// POST /aircraft-registrations - Create a new aircraft registration (HU34)
		protected.POST("/aircraft-registrations", validator.WithValidateCreateAircraftRegistration(), handler.CreateAircraftRegistration())

		// PUT /aircraft-registrations/:id - Update an existing aircraft registration (HU35)
		protected.PUT("/aircraft-registrations/:id", validator.WithValidateUpdateAircraftRegistration(), handler.UpdateAircraftRegistration())

		// ---- Airline Routes Management (Protected - Write Operations) ----
		// PATCH /airline-routes/:id/activate - Activate an airline route (HU42)
		protected.PATCH("/airline-routes/:id/activate", handler.ActivateAirlineRoute())

		// PATCH /airline-routes/:id/deactivate - Deactivate an airline route (HU41)
		protected.PATCH("/airline-routes/:id/deactivate", handler.DeactivateAirlineRoute())

		// ---- Daily Logbook Details Management (Protected - HU15-HU18) ----
		// GET /daily-logbook-details/:id - Get a specific detail (HU15)
		protected.GET("/daily-logbook-details/:id", handler.GetDailyLogbookDetail())

		// PUT /daily-logbook-details/:id - Update a specific detail (HU17)
		protected.PUT("/daily-logbook-details/:id", handler.UpdateDailyLogbookDetail())

		// DELETE /daily-logbook-details/:id - Delete a specific detail (HU18)
		protected.DELETE("/daily-logbook-details/:id", handler.DeleteDailyLogbookDetail())

		// POST /daily-logbooks/:id/details - Add a new detail to logbook (HU16)
		protected.POST("/daily-logbooks/:id/details", handler.CreateDailyLogbookDetail())

		// GET /daily-logbooks/:id/details - List all details for a logbook
		//el id debe ser el id del logbook , no el id del detail, se debe tener en cuenta que el employee_id es el id del employee que esta autenticado y que el detail es un registro de la tabla daily_logbook_details
		protected.GET("/daily-logbooks/:id/details", handler.ListDailyLogbookDetails())

		// ---- Airline Employees Management (Protected) ----
		// GET /airline-employees - List all airline employees (employees with airline assigned)
		// Query params: ?airline_id=xxx (filter by airline), ?active=true/false (filter by status)
		protected.GET("/airline-employees", handler.ListAirlineEmployees())

		// GET /airline-employees/:id - Get airline employee by ID (HU26)
		protected.GET("/airline-employees/:id", handler.GetAirlineEmployeeByID())

		// POST /airline-employees - Create a new airline employee (HU28)
		protected.POST("/airline-employees", handler.CreateAirlineEmployee())

		// PUT /airline-employees/:id - Update an existing airline employee (HU27)
		protected.PUT("/airline-employees/:id", handler.UpdateAirlineEmployee())

		// PATCH /airline-employees/:id/activate - Activate an airline employee (HU29)
		protected.PATCH("/airline-employees/:id/activate", handler.ActivateAirlineEmployee())

		// PATCH /airline-employees/:id/deactivate - Deactivate an airline employee (HU30)
		protected.PATCH("/airline-employees/:id/deactivate", handler.DeactivateAirlineEmployee())
	}

	dependencies.Logger.Success(logger.LogRouteConfigured)
}

func Bootstrap(app *gin.Engine) *dependency.Dependencies {
	// Initialize Prometheus metrics
	dependencies, err := dependency.Init()
	if err != nil {
		slog.Error(logger.LogDepInitError, slog.String("error", err.Error()))
		panic(err)
	}
	routing(app, dependencies)
	return dependencies
}
