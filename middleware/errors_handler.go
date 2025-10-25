package middleware

import (
	"log"
	"net/http"

	"github.com/champion19/Flighthours_backend/core/domain"
	"github.com/gin-gonic/gin"
)
var(
	errorsMap = map[error]ErrorResponse{
		domain.ErrDuplicateUser:{
            Code: http.StatusConflict,
            Message: "User already exists",
        },
	}
)

type 	ErrorResponse struct{
	Code int `json:"code"`
	Message string `json:"message"`
}

// ErrorHandler captures errors and returns a consistent JSON error response
func ErrorHandler() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next() // Step1: Process the request first.

        // Step2: Check if any errors were added to the context
        if len(c.Errors) > 0 {
            // Step3: Use the last error
            err := c.Errors.Last().Err
            errResponse, ok := errorsMap[err]
            if !ok{
                log.Println("Error desconocido",err)
                c.JSON(http.StatusInternalServerError, ErrorResponse{
                    Code: http.StatusInternalServerError,
                    Message: "Error desconocido",
                })
                return
            }
            // Step4: Respond with a generic error message
            c.JSON(http.StatusInternalServerError,ErrorResponse{
                Code: errResponse.Code,
                Message: errResponse.Message,
            })
        }
    }
}


