package server

import (



	"github.com/champion19/flighthours-api/cmd/dependency"
	"github.com/champion19/flighthours-api/handlers"
	"github.com/champion19/flighthours-api/middleware"
	"github.com/champion19/flighthours-api/platform/schema"
	"github.com/gin-gonic/gin"
)

func routing(app *gin.Engine, dependencies *dependency.Dependencies) {
dependencies.Logger.Info("Setting up routes")

	app.Use(middleware.ErrorHandler(dependencies.Logger))
	handler := handlers.New(dependencies.EmployeeService, dependencies.Interactor, dependencies.Logger)

	validators, err := schema.NewValidator(&schema.DefaultFileReader{})
	if err != nil {

		dependencies.Logger.Error("failed to initialize schema validator", err)
		dependencies.Logger.Fatal("failed to initialize schema validator", err)
		return
	}
	validator := middleware.NewMiddlewareValidator(validators)

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

	}

}

func Bootstrap(app *gin.Engine) *dependency.Dependencies {
	dependencies, err := dependency.Init()
	if err != nil {
		dependencies.Logger.Fatal("failed to init dependencies", err)
		return nil
	}
	routing(app, dependencies)
	return dependencies
}
