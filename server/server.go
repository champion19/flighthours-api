package server

import (
	"log"
	"log/slog"

	"github.com/champion19/Flighthours_backend/cmd/dependency"
	"github.com/champion19/Flighthours_backend/handlers"
	"github.com/champion19/Flighthours_backend/middleware"
	"github.com/champion19/Flighthours_backend/platform/schema"

	"github.com/gin-gonic/gin"
)

func routing(app *gin.Engine, dependencies *dependency.Dependencies) {
	slog.Info("Setting up routes")

	handler := handlers.New(dependencies.EmployeeService)

	validators, err := schema.NewValidator(&schema.DefaultFileReader{})
	if err != nil {
		slog.Error("Error creating validator", slog.String("error", err.Error()))
		return
	}
	validator := middleware.NewMiddlewareValidator(validators)

	// Rutas públicas (sin autenticación)
	public := app.Group("/v1/flighthours")
	{
		// Registro de usuario
		public.POST("/register", validator.WithValidateRegister(), handler.RegisterEmployee())

		// Login - devuelve tokens JWT
		public.POST("/login", handler.LoginEmployee())
		public.GET("/user/email/:email", handler.GetEmployeeByEmail())

	}

}

func Bootstrap(app *gin.Engine) *dependency.Dependencies {
	dependencies, err := dependency.Init()
	if err != nil {
		slog.Error("Failed to initialize dependencies", slog.String("error", err.Error()))
		log.Fatal("failed to init dependencies")
		return nil
	}
	routing(app, dependencies)
	return dependencies
}
