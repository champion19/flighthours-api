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
			AccessToken:  result.Token.AccessToken,
			RefreshToken: result.Token.RefreshToken,
			ExpiresIn:    result.Token.ExpiresIn,
			TokenType:    result.Token.TokenType,
		}

		c.JSON(http.StatusCreated, response)
	}
}


func (h handler) LoginEmployee() func(c *gin.Context) {
	return func(c *gin.Context) {
		var loginRequest struct {
			Email    string `json:"email" binding:"required,email"`
			Password string `json:"password" binding:"required"`
		}

		if err := c.ShouldBindJSON(&loginRequest); err != nil {
			c.Error(domain.ErrInvalidRequest)
			return
		}

		token, err := h.EmployeeService.LoginEmployee(loginRequest.Email, loginRequest.Password)
		if err != nil {
			c.Error(domain.ErrUserCannotFound)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"access_token":  token.AccessToken,
			"refresh_token": token.RefreshToken,
			"expires_in":    token.ExpiresIn,
			"token_type":    token.TokenType,
		})
	}
}
