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

	// Rutas públicas (sin autenticación)
	public := app.Group("flighthours/api/v1")
	{
		// Registro de usuario
		public.POST("/register", validator.WithValidateRegister(), handler.RegisterEmployee())

		// GET /employees/:id - Obtener información de un empleado por ID (sin contraseña)
		public.GET("/employees/:id", handler.GetEmployeeByID())

		// PUT /employees/:id - Actualizar información general de un empleado
		// No modifica email ni password (se manejan en endpoints separados)
		// Sincroniza cambios de active y role con Keycloak
		// Valida: longitud de campos, formato de fechas, rol válido
		public.PUT("/employees/:id", validator.WithValidateUpdateEmployee(), handler.UpdateEmployee())

		// DELETE /employees/:id - Eliminar empleado (BD y Keycloak)
		// Esta operación es irreversible
		public.DELETE("/employees/:id", handler.DeleteEmployee())

		// Login - devuelve tokens JWT
		public.POST("/login", handler.Login())
		//public.GET("/user/email/:email", handler.GetEmployeeByEmail())

		// Email verification and password reset
		//POST /auth/resend-verification - Reenviar correo de verificación
		public.POST("/auth/resend-verification", validator.WithValidateResendVerificationEmail(), handler.ResendVerificationEmail())

		//POST /auth/password-reset - Solicitar restablecimiento de contraseña
		public.POST("/auth/password-reset", validator.WithValidatePasswordResetRequest(), handler.RequestPasswordReset())

		//POST /auth/verify-email - Verificar correo con token
		public.POST("/auth/verify-email", handler.VerifyEmailByToken())

		//POST /auth/update-password - Actualizar contraseña con token
		public.POST("/auth/update-password", validator.WithValidateUpdatePassword(), handler.UpdatePassword())

		//POST /auth/change-password - Cambiar contraseña (usuario autenticado conoce su contraseña actual)
		public.POST("/auth/change-password", validator.WithValidateChangePassword(), handler.ChangePassword())

		//Messages Endpoints
		// POST /messages - Crear nuevo mensaje
		public.POST("/messages", validator.WithValidateMessage(), handler.CreateMessage())

		// PUT /messages/:id - Actualizar mensaje existente
		public.PUT("/messages/:id", validator.WithValidateMessage(), handler.UpdateMessage())

		// DELETE /messages/:id - Eliminar mensaje
		public.DELETE("/messages/:id", handler.DeleteMessage())

		// GET /messages/:id - Obtener mensaje por ID
		public.GET("/messages/:id", handler.GetMessageByID())

		// GET /messages - Listar mensajes (con filtros opcionales)
		// Query params: ?module=users&type=ERROR&category=usuario_final&active=true
		public.GET("/messages", handler.ListMessages())

		// POST /messages/cache/reload - Recargar caché de mensajes desde BD
		// Endpoint administrativo para forzar recarga después de cambios manuales
		public.POST("/messages/cache/reload", handler.ReloadMessageCache())

		// Airlines Endpoints
		// GET /airlines/:id - Obtener información de una aerolínea por ID
		public.GET("/airlines/:id", handler.GetAirlineByID())

		// PATCH /airlines/:id/activate - Activar una aerolínea
		public.PATCH("/airlines/:id/activate", handler.ActivateAirline())

		// PATCH /airlines/:id/deactivate - Desactivar una aerolínea
		public.PATCH("/airlines/:id/deactivate", handler.DeactivateAirline())

		dependencies.Logger.Success(logger.LogRouteConfigured)

	}

}

func Bootstrap(app *gin.Engine) *dependency.Dependencies {
	// Inicializar métricas de Prometheus
	dependencies, err := dependency.Init()
	if err != nil {
		slog.Error(logger.LogDepInitError, slog.String("error", err.Error()))
		panic(err)
	}
	routing(app, dependencies)
	return dependencies
}
