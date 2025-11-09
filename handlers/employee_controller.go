package handlers

import (
	"net/http"

	domain "github.com/champion19/flighthours-api/core/interactor/services/domain"
	"github.com/gin-gonic/gin"
)
func (h handler) GetEmployeeByEmail() func(c *gin.Context) {
	return func(c *gin.Context) {
		email := c.Param("email")

		person, err := h.EmployeeService.GetEmployeeByEmail(c,email)
		if err != nil {
			c.Error(err)
			return
		}

		c.JSON(http.StatusOK, person)
	}
}


func (h handler) RegisterEmployee() func(c *gin.Context) {
	return func(c *gin.Context) {
		var employeeRequest EmployeeRequest
		if err := c.ShouldBindJSON(&employeeRequest); err != nil {
			c.Error(domain.ErrInvalidJSONFormat)
			return
		}

		result, err := h.Interactor.RegisterEmployee(c,employeeRequest.ToDomain())
		if err != nil {
			c.Error(err)
			return
		}

		response := RegisterEmployeeResponse{
			User: EmployeeResponse{
				ID:                   result.Employee.ID,
				Name:                 result.Employee.Name,
				Email:                result.Employee.Email,
				Airline:              result.Employee.Airline,
				IdentificationNumber: result.Employee.IdentificationNumber,
				Bp:                   result.Employee.Bp,
				StartDate:            result.Employee.StartDate,
				EndDate:              result.Employee.EndDate,
				Active:               result.Employee.Active,
				Role:                 result.Employee.Role,
			},
			Message: result.Message,
		}

		c.JSON(http.StatusCreated, response)
	}
}
