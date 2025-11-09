package middleware

import (
	"net/http"

	"github.com/champion19/flighthours-api/core/interactor/services/domain"
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

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			if response, ok := mapError[err]; ok {
				c.JSON(response.Status, response)
				return
			}

			c.JSON(http.StatusInternalServerError, map[string]any{
				"success": false,
				"message": err.Error(),
			})

		}
	}

}
