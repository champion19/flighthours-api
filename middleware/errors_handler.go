package middleware

import (
	"net/http"

	"github.com/champion19/flighthours-api/core/interactor/services/domain"
	"github.com/champion19/flighthours-api/platform/logger"
	"github.com/gin-gonic/gin"
)

var mapError = map[error]ErrorResponse{
	domain.ErrDuplicateUser: {
		Code:    "MODE_U_USU_ERR_00001",
		Message: "User already exists",
		Status:  http.StatusConflict,
	},
}

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Status  int    `json:"-"`
}

func ErrorHandler(logger logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			if response, ok := mapError[err]; ok {
				logger.Warn("Error handling error",
					"error", err.Error(),
					"code", response.Code,
					"status", response.Status,
					"path", c.Request.URL.Path,
					"method", c.Request.Method,
					"employee_ip", c.ClientIP())
				c.JSON(response.Status, response)
				return
			}

			logger.Error("Error handling error",
				"error", err.Error(),
				"path", c.Request.URL.Path,
				"method", c.Request.Method,
				"employee_ip", c.ClientIP())
				
			c.JSON(http.StatusInternalServerError, map[string]any{
				"success": false,
				"message": err.Error(),
			})

		}
	}

}
