package handlers

import (
	domain "github.com/champion19/flighthours-api/core/interactor/services/domain"
	"github.com/champion19/flighthours-api/middleware"
	"github.com/champion19/flighthours-api/platform/logger"
	"github.com/gin-gonic/gin"
)

// RegisterEmployee godoc
// @Summary      Registrar nueva cuenta de empleado
// @Description  Crea una nueva cuenta de empleado en el sistema con sincronización a Keycloak. Incluye validación de datos, verificación de duplicados y creación de usuario en el sistema de autenticación.
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param        account  body      EmployeeRequest  true  "Datos del empleado a registrar"
// @Success      201      {object}  middleware.APIResponse{data=RegisterEmployeeResponse}  "Cuenta creada exitosamente"
// @Failure      400      {object}  middleware.APIResponse  "Error de validación - Datos inválidos o incompletos"
// @Failure      409      {object}  middleware.APIResponse  "Conflicto - Email o número de identidad ya registrado"
// @Failure      500      {object}  middleware.APIResponse  "Error interno del servidor"
// @Router       /register [post]
func (h handler) RegisterEmployee() func(c *gin.Context) {

	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		log.Info(logger.LogRegRequestReceived,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"client_ip", c.ClientIP())

		var employeeRequest EmployeeRequest
		if err := c.ShouldBindJSON(&employeeRequest); err != nil {
			log.Error(logger.LogRegJSONParseError, "error", err, "client_ip", c.ClientIP())
			c.Error(domain.ErrInvalidJSONFormat)
			return
		}

		log.Info(logger.LogRegProcessing,
			"email", employeeRequest.Email,
			"role", employeeRequest.Role)

		result, err := h.Interactor.RegisterEmployee(c, employeeRequest.ToDomain())
		if err != nil {
			log.Error(logger.LogRegProcessError,
				"email", employeeRequest.Email,
				"error", err,
				"client_ip", c.ClientIP())
			c.Error(err)
			return
		}

		// Ofuscar el ID antes de exponerlo en la API
		encodedID, err := h.EncodeID(result.Employee.ID)
		if err != nil {
			h.HandleIDEncodingError(c, result.Employee.ID, err)
			return
		}

		//TODO;TENERLO EN CUENTA, ESTO ES DE COOKIES HTTTPONLY
		c.SetCookie(
			"employee_id",        // name
			result.Employee.ID,   // value
			3600,                 // expira en 1 hora
			"/",                  // path
			c.Request.Host,       // domain
			c.Request.TLS != nil, // secure
			true,                 // httpOnly
		)

		// Usar funciones HATEOAS centralizadas
		baseURL := GetBaseURL(c)
		links := BuildAccountLinks(baseURL, encodedID)
		SetLocationHeader(c, baseURL, "accounts", encodedID)

		response := RegisterEmployeeResponse{
			Links: links,
		}

		log.Success("register employee success",
			result.Employee.ToLogger(),
			"encoded_id", encodedID,
			"client_ip", c.ClientIP())

		h.Response.SuccessWithData(c, domain.MsgUserRegistered, response)
	}
}
