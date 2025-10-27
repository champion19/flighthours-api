package handlers

import (
	"context"
	"net/http"

	domain "github.com/champion19/flighthours-api/core/domain"
	"github.com/champion19/flighthours-api/core/ports"
	"github.com/gin-gonic/gin"
)
type AuthorizationController struct {
	authService ports.AuthorizationService
}

func NewAuthorizationController(authService ports.AuthorizationService) AuthorizationController {
	return AuthorizationController{
		authService: authService,
	}
}

func (a *AuthorizationController) SyncUserToKeycloak() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request struct {
			EmployeeID string `json:"employee_id" binding:"required"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.Error(domain.ErrPersonNotFound)
			return
		}

		// Aquí necesitarías obtener la persona de tu base de datos
		// person, err := a.personService.GetPersonByID(request.PersonID)
		// if err != nil {
		//     c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
		//     return
		// }

		// ctx := context.Background()
		// keycloakUserID, err := a.authzService.SyncUserToKeycloak(ctx, person)
		// if err != nil {
		//     c.JSON(http.StatusInternalServerError, gin.H{
		//         "error": "Error sincronizando usuario con Keycloak: " + err.Error(),
		//     })
		//     return
		// }

		c.JSON(http.StatusOK, gin.H{
			"message": "Usuario sincronizado con Keycloak exitosamente",
			// "keycloak_user_id": keycloakUserID,
		})
	}
}

func (a *AuthorizationController) AssignRole() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request struct {
			EmployeeID string `json:"employee_id" binding:"required"`
			Role       string `json:"role" binding:"required"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.Error(domain.ErrInvalidRequest)
			return
	}

		ctx := context.Background()
		err := a.authService.AssignRole(ctx, request.EmployeeID, request.Role)
		if err != nil {
			c.Error(domain.ErrRoleAssignmentFailed)
			return
	}

		c.JSON(http.StatusOK, gin.H{
			"message": "Rol asignado exitosamente",
			"employee_id": request.EmployeeID,
			"role": request.Role,
		})
	}
}

func (a *AuthorizationController) RemoveRole() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request struct {
			EmployeeID string `json:"employee_id" binding:"required"`
			Role       string `json:"role" binding:"required"`
		}
    if err := c.ShouldBindJSON(&request); err != nil {
			c.Error(domain.ErrInvalidRequest)
			return
	}
	ctx := context.Background()
err := a.authService.RemoveRole(ctx, request.EmployeeID, request.Role)
if err != nil {
    c.Error(domain.ErrRoleRemovalFailed)
    return
}

c.JSON(http.StatusOK, gin.H{
    "message": "Rol removido exitosamente",
})
	}
}

func (a *AuthorizationController) GetUserRoles() gin.HandlerFunc {
	return func(c *gin.Context) {
	employeeID := c.Param("employee_id")
		if employeeID == "" {
			c.Error(domain.ErrInvalidRequest)
			return
		}

		ctx := context.Background()
		roles, err := a.authService.GetUserRoles(ctx, employeeID)
		if err != nil {
			c.Error(domain.ErrGetUserRolesFailed)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"employee_id": employeeID,
			"roles": roles,
		})
	}
}

func (a *AuthorizationController) CheckRole() gin.HandlerFunc {
	return func(c *gin.Context) {
		employeeID := c.Param("employee_id")
		role := c.Param("role")

		if employeeID == "" || role == "" {
			c.Error(domain.ErrInvalidRequest)
			return
		}

		ctx := context.Background()
		hasRole, err := a.authService.HasRole(ctx, employeeID, role)
		if err != nil {
			c.Error(domain.ErrRoleCheckFailed)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"employee_id": employeeID,
			"role": role,
			"has_role": hasRole,
		})
	}
}
