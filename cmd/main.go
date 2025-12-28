package main

import (
	"log/slog"

	"github.com/champion19/flighthours-api/platform/logger"
	_ "github.com/champion19/flighthours-api/platform/swaggo"
	"github.com/champion19/flighthours-api/server"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// @title           Flighthours Backend API
// @version         1.0
// @description     API RESTful para la plataforma Flighthours, implementada con arquitectura hexagonal y siguiendo los principios del Richardson Maturity Model (Nivel 2-3) con HATEOAS.
// @termsOfService  http://swagger.io/terms/

// @contact.name   Flighthours API Support
// @contact.url    https://flighthours.com/support
// @contact.email  support@flighthours.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8081
// @BasePath  /flighthours/api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/

func main() {
	gin.SetMode(gin.ReleaseMode)
	app := gin.New()
	app.Use(gin.Logger())
	app.Use(gin.Recovery())

	// Configurar CORS para permitir requests desde Swagger UI
	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3001", "http://localhost:8081", "http://localhost:8080"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},  
		AllowCredentials: true,
	}))

	dependencies := server.Bootstrap(app)
	serverAddr := dependencies.Config.GetServerAddress()
	slog.Info(logger.LogAppServerStarting, slog.String("address", serverAddr))

	if err := app.Run(serverAddr); err != nil {
		slog.Error(logger.LogAppServerStartError, slog.String("error", err.Error()))
		return
	}

}
