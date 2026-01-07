package middleware

import (
	"strings"

	"github.com/champion19/flighthours-api/core/interactor/services/domain"
	"github.com/champion19/flighthours-api/core/ports/input"
	"github.com/champion19/flighthours-api/platform/cache/messaging"
	"github.com/champion19/flighthours-api/platform/jwt"
	"github.com/gin-gonic/gin"
)

// RequireAuth creates a middleware that validates JWT tokens from Keycloak
// and injects the authenticated user into the Gin context
func RequireAuth(employeeService input.Service, msgCache *messaging.MessageCache) gin.HandlerFunc {
	tokenParser := jwt.NewTokenParser()

	return func(c *gin.Context) {
		// Extract Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Error(domain.ErrInvalidToken)
			c.Abort()
			return
		}

		// Check Bearer prefix
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.Error(domain.ErrInvalidToken)
			c.Abort()
			return
		}

		token := parts[1]

		// Decode JWT and extract claims
		// Note: We extract the "sub" claim which contains the Keycloak User ID
		claims, err := tokenParser.ExtractClaimsFromToken(token)
		if err != nil {
			c.Error(domain.ErrInvalidToken)
			c.Abort()
			return
		}

		// Extract Keycloak User ID from "sub" claim
		keycloakUserID, ok := claims["sub"].(string)
		if !ok || keycloakUserID == "" {
			c.Error(domain.ErrInvalidToken)
			c.Abort()
			return
		}

		// Find user in database by Keycloak ID
		employee, err := employeeService.GetEmployeeByKeycloakID(c.Request.Context(), keycloakUserID)
		if err != nil {
			// User not found in our database
			c.Error(domain.ErrUserNotFound)
			c.Abort()
			return
		}

		// Inject authenticated user into context
		c.Set("authenticated_user", employee)

		c.Next()
	}
}

// GetAuthenticatedUser extracts the authenticated user from the Gin context
func GetAuthenticatedUser(c *gin.Context) (*domain.Employee, bool) {
	user, exists := c.Get("authenticated_user")
	if !exists {
		return nil, false
	}

	employee, ok := user.(*domain.Employee)
	return employee, ok
}
