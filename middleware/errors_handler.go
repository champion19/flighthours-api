package middleware

import (
	"errors"
	"net/http"

	"github.com/champion19/flighthours-api/core/interactor/services/domain"
	"github.com/champion19/flighthours-api/platform/logger"
	"github.com/gin-gonic/gin"
)

func ErrorHandler(logger logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			var appErr *domain.AppError

			if errors.As(err, &appErr) {
				logger.Warn("Error handling error",
					"error", err.Error(),
					"code", appErr.Code,
					"status", appErr.StatusCode,
					"path", c.Request.URL.Path,
					"method", c.Request.Method,
					"employee_ip", c.ClientIP())

				c.JSON(appErr.StatusCode, gin.H{
					"code":    appErr.Code,
					"message": appErr.Message,
				})
				return
			}

			logger.Error("Unhandled error",
				"error", err.Error(),
				"path", c.Request.URL.Path,
				"method", c.Request.Method,
				"employee_ip", c.ClientIP())

			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    "INTERNAL_SERVER_ERROR",
				"message": err.Error(),
			})

		}
	}

}
