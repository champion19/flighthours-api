package middleware

import (
	"log"
	"net/http"
    "errors"
	"github.com/champion19/flighthours-api/core/domain"
	"github.com/gin-gonic/gin"
)
var(
	errorStatusMap = map[string]int{
		"MOD_U_USU_ERR_00001": http.StatusConflict,
        "MOD_U_USU_ERR_00002": http.StatusInternalServerError,
        "MOD_U_USU_ERR_00003": http.StatusNotFound,
        "MOD_U_USU_ERR_00004": http.StatusNotFound,
        "MOD_U_USU_ERR_00005": http.StatusNotFound,
        "MOD_U_USU_ERR_00006": http.StatusNotFound,
        "MOD_U_USU_ERR_00007": http.StatusInternalServerError,
        "MOD_U_USU_ERR_00008": http.StatusForbidden,
        "MOD_U_USU_ERR_00009": http.StatusNotFound,
        "MOD_U_USU_ERR_00010": http.StatusGone,
        "MOD_U_USU_ERR_00011": http.StatusConflict,
        "MOD_U_USU_ERR_00012": http.StatusInternalServerError,
        "MOD_U_USU_ERR_00013": http.StatusBadRequest,
        "Mod_U_USU_ERR_00014": http.StatusInternalServerError,


    "MOD_V_VAL_ERR_00001": http.StatusBadRequest,
	"MOD_V_VAL_ERR_00002": http.StatusBadRequest,


	"MOD_A_AUT_ERR_00001": http.StatusInternalServerError,
	"MOD_A_AUT_ERR_00002": http.StatusInternalServerError,
	"MOD_A_AUT_ERR_00003": http.StatusInternalServerError,
	"MOD_A_AUT_ERR_00004": http.StatusInternalServerError,

    }
)

type 	ErrorResponse struct{
    Status int `json:"status"`
	Code string `json:"code"`
	Message string `json:"message"`
}


// ErrorHandler captures errors and returns a consistent JSON error response
func ErrorHandler() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			var domainErr *domain.DomainError
			if errors.As(err, &domainErr) {
				statusCode, exists := errorStatusMap[domainErr.Code]
				if !exists {
					statusCode = http.StatusInternalServerError
					log.Printf("Unknown error code: %s", domainErr.Code)
				}

				response := ErrorResponse{
					Status:  statusCode,
					Code:    domainErr.Code,
					Message: domainErr.Message,
				}

				c.AbortWithStatusJSON(statusCode, response)
				return
			}

			log.Printf("Non-domain error: %v", err)
			response := ErrorResponse{
				Status:  http.StatusInternalServerError,
				Code:    "MOD_G_GEN_ERR_00001",
				Message: "Error interno del servidor",
			}
			c.AbortWithStatusJSON(http.StatusInternalServerError, response)
		}

    }
}


