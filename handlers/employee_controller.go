package handlers

import (
	"net/http"

	domain "github.com/champion19/flighthours-api/core/domain"
	"github.com/gin-gonic/gin"
)

func (h handler) GetEmployeeByEmail() func(c *gin.Context) {
	return func(c *gin.Context) {
		email := c.Param("email")

		employee, err := h.EmployeeService.GetEmployeeByEmail(email)
		if err != nil {
			c.Error(err)
			return
		}
		c.JSON(http.StatusOK, employee)
	}
}

func (h handler) RegisterEmployee() func(c *gin.Context) {
	return func(c *gin.Context) {
		var employeeRequest EmployeeRequest
		if err := c.ShouldBindJSON(&employeeRequest); err != nil {
			c.Error(domain.ErrInvalidJSONFormat)
			return
		}

		result, err := h.EmployeeService.RegisterEmployee(employeeRequest.ToDomain())
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
				Emailconfirmed:       result.Employee.Emailconfirmed,
				IdentificationNumber: result.Employee.IdentificationNumber,
				Bp:                   result.Employee.Bp,
				StartDate:            result.Employee.StartDate,
				EndDate:              result.Employee.EndDate,
				Active:               result.Employee.Active,
				Role:                 result.Employee.Role,
				KeycloakUserID:       result.Employee.KeycloakUserID,
			},
			Message: result.Message,
		}

		c.JSON(http.StatusCreated, response)
	}
}


func (h handler) LoginEmployee() func(c *gin.Context) {
	return func(c *gin.Context) {
		var loginRequest LoginRequest
		if err := c.ShouldBindJSON(&loginRequest); err != nil {
			c.Error(domain.ErrInvalidRequest)
			return
		}

		token, err := h.EmployeeService.LoginEmployee(loginRequest.Email, loginRequest.Password)
		if err != nil {
			c.Error(err)
			return
		}

		response:= LoginResponse{
			AccessToken:  token.AccessToken,
			RefreshToken: token.RefreshToken,
			ExpiresIn:    token.ExpiresIn,
			TokenType:    token.TokenType,
		}
		c.JSON(http.StatusOK, response)
	}
}
