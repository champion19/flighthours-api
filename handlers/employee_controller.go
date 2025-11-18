package handlers

import (
	"net/http"

	domain "github.com/champion19/flighthours-api/core/interactor/services/domain"
	"github.com/gin-gonic/gin"

)


func (h handler) RegisterEmployee() func(c *gin.Context) {
	return func(c *gin.Context) {
		var employeeRequest EmployeeRequest
		if err := c.ShouldBindJSON(&employeeRequest); err != nil {
			h.Logger.Error("Error binding JSON", err)
			c.Error(domain.ErrInvalidJSONFormat)
			return
		}

		result, err := h.Interactor.RegisterEmployee(c,employeeRequest.ToDomain())
		if err != nil {
			h.Logger.Error("Error registering employee", err)
			c.Error(err)
			return
		}

		//TODO;TENERLO EN CUENTA, ESTO ES DE COOKIES HTTTPONLY
		c.SetCookie(
            "employee_id",               // name
            result.Employee.ID,          // value
            3600,                        // expira en 1 hora
            "/",                         // path
            c.Request.Host,              // domain
            c.Request.TLS != nil,        // secure
            true,                        // httpOnly
        )

		scheme := "http"
		if c.Request.TLS != nil {
			scheme = "https"
		}
		baseURL := scheme + "://" + c.Request.Host
		links := BuildAccountLinks(baseURL, result.Employee.ID)

		locationURL := baseURL + "/flighthours/api/v1/accounts/" + result.Employee.ID
		c.Header("Location", locationURL)

		response := RegisterEmployeeResponse{
			Message: result.Message,
			Links:links,
		}
		h.Logger.Success("Employee registered successfully", result.Employee.ID)

		c.JSON(http.StatusCreated, response)
	}
}
