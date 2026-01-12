package handlers

import (
	"context"

	domain "github.com/champion19/flighthours-api/core/interactor/services/domain"
	"github.com/champion19/flighthours-api/middleware"
	"github.com/gin-gonic/gin"
)

// GetAirlineRoute handles GET /airline-routes/:id (HU40)
func (h *handler) GetAirlineRoute() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		ctx := context.Background()

		// Resolve ID (supports both obfuscated and raw UUID)
		id := c.Param("id")
		uuid, responseID := h.resolveID(id)

		if uuid == "" {
			c.Error(domain.ErrInvalidID)
			return
		}

		// Get airline route from interactor
		airlineRoute, err := h.AirlineRouteInteractor.GetAirlineRouteByID(ctx, traceID, uuid)
		if err != nil {
			c.Error(err)
			return
		}

		// Build response
		response := FromDomainAirlineRoute(airlineRoute, responseID)
		baseURL := GetBaseURL(c)
		response.Links = BuildAirlineRouteLinks(baseURL, responseID, airlineRoute.Status)

		h.Response.SuccessWithData(c, domain.MsgAirlineRouteGetOK, response)
	}
}

// ListAirlineRoutes handles GET /airline-routes
func (h *handler) ListAirlineRoutes() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		ctx := context.Background()

		// Build filters from query params
		filters := make(map[string]interface{})

		if airlineCode := c.Query("airline_code"); airlineCode != "" {
			filters["airline_code"] = airlineCode
		}

		if status := c.Query("status"); status != "" {
			if status == "true" || status == "1" {
				filters["status"] = true
			} else if status == "false" || status == "0" {
				filters["status"] = false
			}
		}

		// Get airline routes from interactor
		airlineRoutes, err := h.AirlineRouteInteractor.ListAirlineRoutes(ctx, traceID, filters)
		if err != nil {
			c.Error(err)
			return
		}

		// Build response
		baseURL := GetBaseURL(c)
		response := ToAirlineRouteListResponse(airlineRoutes, h.EncodeID, baseURL)

		h.Response.SuccessWithData(c, domain.MsgAirlineRouteListOK, response)
	}
}

// ActivateAirlineRoute handles PATCH /airline-routes/:id/activate (HU42)
func (h *handler) ActivateAirlineRoute() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		ctx := context.Background()

		// Resolve ID (supports both obfuscated and raw UUID)
		id := c.Param("id")
		uuid, responseID := h.resolveID(id)

		if uuid == "" {
			c.Error(domain.ErrInvalidID)
			return
		}

		// Activate airline route
		err := h.AirlineRouteInteractor.ActivateAirlineRoute(ctx, traceID, uuid)
		if err != nil {
			c.Error(err)
			return
		}

		// Build response with HATEOAS links
		baseURL := GetBaseURL(c)
		links := BuildAirlineRouteStatusLinks(baseURL, responseID, true)

		h.Response.SuccessWithData(c, domain.MsgAirlineRouteActivateOK, gin.H{
			"id":     responseID,
			"status": true,
			"_links": links,
		})
	}
}

// DeactivateAirlineRoute handles PATCH /airline-routes/:id/deactivate (HU41)
func (h *handler) DeactivateAirlineRoute() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		ctx := context.Background()

		// Resolve ID (supports both obfuscated and raw UUID)
		id := c.Param("id")
		uuid, responseID := h.resolveID(id)

		if uuid == "" {
			c.Error(domain.ErrInvalidID)
			return
		}

		// Deactivate airline route
		err := h.AirlineRouteInteractor.DeactivateAirlineRoute(ctx, traceID, uuid)
		if err != nil {
			c.Error(err)
			return
		}

		// Build response with HATEOAS links
		baseURL := GetBaseURL(c)
		links := BuildAirlineRouteStatusLinks(baseURL, responseID, false)

		h.Response.SuccessWithData(c, domain.MsgAirlineRouteDeactivateOK, gin.H{
			"id":     responseID,
			"status": false,
			"_links": links,
		})
	}
}
