package main

import (
	"log/slog"

	"github.com/champion19/flighthours-api/platform/logger"
	"github.com/champion19/flighthours-api/server"
	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	app := gin.New()
	app.Use(gin.Logger())
	app.Use(gin.Recovery())

	dependencies := server.Bootstrap(app)
	serverAddr := dependencies.Config.GetServerAddress()
	slog.Info(logger.LogAppServerStarting, slog.String("address", serverAddr))

	if err := app.Run(serverAddr); err != nil {
		slog.Error(logger.LogAppServerStartError, slog.String("error", err.Error()))
		return
	}

}
