package handlers

import (
	domain "github.com/champion19/flighthours-api/core/interactor/services/domain"
	"github.com/champion19/flighthours-api/platform/logger"
	"github.com/champion19/flighthours-api/platform/prometheus"
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
		var employeeRequest EmployeeRequest
		if err := c.ShouldBindJSON(&employeeRequest); err != nil {
			h.Logger.Error(logger.LogRegJSONBindError, err)
			c.Error(domain.ErrInvalidJSONFormat)
			return
		}

		result, err := h.Interactor.RegisterEmployee(c, employeeRequest.ToDomain())
		if err != nil {
			h.Logger.Error(logger.LogEmployeeRegisterError, err)
			c.Error(err)
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
		SetLocationHeader(c, baseURL, "accounts", result.Employee.ID)

		response := RegisterEmployeeResponse{
			Message: result.Message,
			Links:   BuildAccountLinks(baseURL, result.Employee.ID),
		}

		// Registrar métrica de empleado registrado
		prometheus.EmployeesRegistered.Inc()

		h.Logger.Success(logger.LogEmployeeRegisterSuccess, result.Employee.ID)
		h.Response.SuccessWithData(c, domain.MsgUserRegistered, response)
	}
}
