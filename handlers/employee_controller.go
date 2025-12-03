package handlers

import (
	domain "github.com/champion19/flighthours-api/core/interactor/services/domain"
	"github.com/champion19/flighthours-api/platform/logger"
	"github.com/gin-gonic/gin"
)

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

		h.Logger.Success(logger.LogEmployeeRegisterSuccess, result.Employee.ID)
		h.Response.SuccessWithData(c, domain.MsgUserRegistered, response)
	}
}
