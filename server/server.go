package server

import (
	"log/slog"

	"github.com/champion19/flighthours-api/cmd/dependency"
	"github.com/champion19/flighthours-api/handlers"
	"github.com/champion19/flighthours-api/middleware"
	"github.com/champion19/flighthours-api/platform/logger"
	"github.com/champion19/flighthours-api/platform/prometheus"
	"github.com/champion19/flighthours-api/platform/schema"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func routing(app *gin.Engine, dependencies *dependency.Dependencies) {
	dependencies.Logger.Info(logger.LogRouteConfiguring)

	// Endpoint de métricas de Prometheus
	app.GET("/metrics", prometheus.Handler())

	// Documentación Swagger
	app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	dependencies.Logger.Info("Swagger UI disponible en: /swagger/index.html")

	// Middleware de Prometheus para capturar métricas
	app.Use(middleware.PrometheusMiddleware())

	errorHandler := middleware.NewErrorHandler(dependencies.MessagingCache, dependencies.Logger)
	app.Use(errorHandler.Handle())

	handler := handlers.New(dependencies.EmployeeService, dependencies.Interactor, dependencies.Logger, dependencies.IDEncoder, dependencies.ResponseHandler, dependencies.MessageInteractor, dependencies.MessagingCache)

	validators, err := schema.NewValidator(&schema.DefaultFileReader{})
	if err != nil {

		dependencies.Logger.Error(logger.LogRouteValidatorError, err)
		dependencies.Logger.Fatal(logger.LogRouteValidatorError, err)
		return
	}
	validator := middleware.NewMiddlewareValidator(validators, dependencies.Logger)

	// Rutas públicas (sin autenticación)
	public := app.Group("flighthours/api/v1")
	{
		// Registro de usuario
		public.POST("/register", validator.WithValidateRegister(), handler.RegisterEmployee())
		//GET/accounts/:id
		//public.GET("/accounts/:id", handler.GetEmployeeByID())

		// Login - devuelve tokens JWT
		//public.POST("/login", handler.LoginEmployee())
		//public.GET("/user/email/:email", handler.GetEmployeeByEmail())

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

	}

}

func Bootstrap(app *gin.Engine) *dependency.Dependencies {
	// Inicializar métricas de Prometheus
	if err := prometheus.InitMetrics(); err != nil {
		slog.Error(logger.LogPrometheusInitError, slog.String("error", err.Error()))
	}

	dependencies, err := dependency.Init()
	if err != nil {
		slog.Error(logger.LogDepInitError, slog.String("error", err.Error()))
		return nil
	}
	routing(app, dependencies)
	return dependencies
}
