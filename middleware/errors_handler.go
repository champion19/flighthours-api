package middleware

import (
	"errors"
	"log"
	"net/http"

	"github.com/champion19/flighthours-api/core/domain"
	json_schema "github.com/champion19/flighthours-api/platform/schema"
	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Status  int    `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

var (
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
		"MOD_V_VAL_ERR_00003": http.StatusInternalServerError,
		"MOD_V_VAL_ERR_00004": http.StatusInternalServerError,
		"MOD_V_VAL_ERR_00005": http.StatusInternalServerError,
		"MOD_V_VAL_ERR_00006": http.StatusBadRequest,
		"MOD_V_VAL_ERR_00007": http.StatusBadRequest,
		"MOD_V_VAL_ERR_00008": http.StatusBadRequest,
		"MOD_V_VAL_ERR_00009": http.StatusBadRequest,
		"MOD_V_VAL_ERR_00010": http.StatusBadRequest,
		"MOD_V_VAL_ERR_00011": http.StatusBadRequest,

		"MOD_A_AUT_ERR_00001": http.StatusInternalServerError,
		"MOD_A_AUT_ERR_00002": http.StatusInternalServerError,
		"MOD_A_AUT_ERR_00003": http.StatusInternalServerError,
		"MOD_A_AUT_ERR_00004": http.StatusInternalServerError,
	}
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			// Try SchemaError first
			var schemaErr *json_schema.SchemaError
			if errors.As(err, &schemaErr) {
				statusCode, exists := errorStatusMap[schemaErr.Code]
				if !exists {
					statusCode = http.StatusInternalServerError
					log.Printf("Unknown schema error code: %s", schemaErr.Code)
				}

				response := ErrorResponse{
					Status:  statusCode,
					Code:    schemaErr.Code,
					Message: schemaErr.Message,
				}

				c.AbortWithStatusJSON(statusCode, response)
				return
			}

			// Try DomainError second
			var domainErr *domain.DomainError
			if errors.As(err, &domainErr) {
				statusCode, exists := errorStatusMap[domainErr.Code]
				if !exists {
					statusCode = http.StatusInternalServerError
					log.Printf("Unknown domain error code: %s", domainErr.Code)
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


