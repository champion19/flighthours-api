// handlers/messages.go
package handlers

import (
	"net/http"

	domain "github.com/champion19/flighthours-api/core/interactor/services/domain"
	"github.com/gin-gonic/gin"
)

func (h *handler) ReloadMessages() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := h.MessageManager.Reload(c.Request.Context()); err != nil {
			h.Logger.Error("Failed to reload messages", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to reload messages",
			})
			return
		}

		// Re-inicializar errores globales con nuevos mensajes
		domain.InitErrors()

		h.Logger.Success("Messages reloaded successfully")
		c.JSON(http.StatusOK, gin.H{
			"message": "messages reloaded successfully",
		})
	}
}
