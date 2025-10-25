package middleware

import (
	"context"
	"net/http"

	"github.com/champion19/Flighthours_backend/core/ports"
	"github.com/gin-gonic/gin"
)


type AuthorizationMiddleware struct {
	authzService ports.AuthorizationService
}

func NewAuthorizationMiddleware(authzService ports.AuthorizationService) *AuthorizationMiddleware {
	return &AuthorizationMiddleware{
		authzService: authzService,
	}
}

func (a *AuthorizationMiddleware) RequireRole(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {

		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Usuario no autenticado",
			})
			c.Abort()
			return
		}

		userIDStr, ok := userID.(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "ID de usuario inválido",
			})
			c.Abort()
			return
		}

		ctx := context.Background()
		hasRole, err := a.authzService.HasRole(ctx, userIDStr, requiredRole)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error verificando permisos",
			})
			c.Abort()
			return
		}

		if !hasRole {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "No tienes permisos para acceder a este recurso",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func (a *AuthorizationMiddleware) RequireAnyRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Usuario no autenticado",
			})
			c.Abort()
			return
		}

		userIDStr, ok := userID.(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "ID de usuario inválido",
			})
			c.Abort()
			return
		}

		ctx := context.Background()

		for _, role := range roles {
			hasRole, err := a.authzService.HasRole(ctx, userIDStr, role)
			if err != nil {
				continue
			}
			if hasRole {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{
			"error": "No tienes permisos para acceder a este recurso",
		})
		c.Abort()
	}
}

func (a *AuthorizationMiddleware) RequirePermission(resource, action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Usuario no autenticado",
			})
			c.Abort()
			return
		}

		userIDStr, ok := userID.(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "ID de usuario inválido",
			})
			c.Abort()
			return
		}

		ctx := context.Background()
		hasPermission, err := a.authzService.HasPermission(ctx, userIDStr, resource, action)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error verificando permisos",
			})
			c.Abort()
			return
		}

		if !hasPermission {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "No tienes permisos para realizar esta acción",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func (a *AuthorizationMiddleware) GetUserRoles() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Usuario no autenticado",
			})
			c.Abort()
			return
		}

		userIDStr, ok := userID.(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "ID de usuario inválido",
			})
			c.Abort()
			return
		}

		ctx := context.Background()
		roles, err := a.authzService.GetUserRoles(ctx, userIDStr)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error obteniendo roles",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"user_id": userIDStr,
			"roles":   roles,
		})
	}
}

func (a *AuthorizationMiddleware) AssignRole() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request struct {
			UserID string `json:"user_id" binding:"required"`
			Role   string `json:"role" binding:"required"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Datos inválidos",
			})
			return
		}

		ctx := context.Background()
		err := a.authzService.AssignRole(ctx, request.UserID, request.Role)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error asignando rol: " + err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Rol asignado exitosamente",
		})
	}
}
