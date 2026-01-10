package handlers

import (
	domain "github.com/champion19/flighthours-api/core/interactor/services/domain"
	"github.com/champion19/flighthours-api/middleware"
	"github.com/champion19/flighthours-api/platform/logger"
	"github.com/gin-gonic/gin"
)

// GetRouteByID godoc
// @Summary      Get route by ID
// @Description  Returns route information by ID (accepts both UUID and obfuscated ID)
// @Tags         Routes
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Route ID (UUID or obfuscated)"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /routes/{id} [get]
func (h *handler) GetRouteByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		inputID := c.Param("id")
		if inputID == "" {
			log.Error(logger.LogMessageIDDecodeError, "error", "empty id parameter", "client_ip", c.ClientIP())
			h.Response.Error(c, domain.MsgValIDInvalid)
			return
		}

		log.Info(logger.LogRouteGet, "input_id", inputID, "client_ip", c.ClientIP())

		var routeUUID string
		var responseID string

		// Detect if it's a valid UUID or obfuscated ID
		if isValidUUID(inputID) {
			// It's a direct UUID
			routeUUID = inputID
			// Encode UUID for response (maintain consistency)
			encodedID, err := h.EncodeID(inputID)
			if err != nil {
				log.Warn(logger.LogIDEncodeError, "uuid", inputID, "error", err)
				responseID = inputID
			} else {
				responseID = encodedID
			}
			log.Debug(logger.LogRouteGet, "detected_format", "UUID", "uuid", routeUUID)
		} else {
			// It's an obfuscated ID, decode it
			uuid, err := h.DecodeID(inputID)
			if err != nil {
				h.HandleIDDecodingError(c, inputID, err)
				return
			}
			routeUUID = uuid
			responseID = inputID
			log.Debug(logger.LogRouteGet, "detected_format", "encoded", "decoded_uuid", routeUUID)
		}

		// Get route from interactor
		route, err := h.RouteInteractor.GetRouteByID(c.Request.Context(), routeUUID)
		if err != nil {
			log.Error(logger.LogRouteGetError, "route_id", routeUUID, "error", err, "client_ip", c.ClientIP())
			if err == domain.ErrRouteNotFound {
				h.Response.Error(c, domain.MsgRouteNotFound)
				return
			}
			h.Response.Error(c, domain.MsgServerError)
			return
		}

		response := FromDomainRoute(route, responseID)

		// Build HATEOAS links
		baseURL := GetBaseURL(c)
		response.Links = BuildRouteLinks(baseURL, responseID)

		log.Success(logger.LogRouteGetOK, route.ToLogger())
		h.Response.SuccessWithData(c, domain.MsgRouteGetOK, response)
	}
}

// ListRoutes godoc
// @Summary      List all routes
// @Description  Returns a list of all routes with optional filters
// @Tags         Routes
// @Produce      json
// @Param        airport_type query string false "Filter by airport type (e.g., Nacional, Internacional)"
// @Param        origin_country query string false "Filter by origin country (e.g., Colombia)"
// @Success      200  {object}  RouteListResponse
// @Failure      500  {object}  map[string]interface{}
// @Router       /routes [get]
func (h *handler) ListRoutes() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		log.Debug(logger.LogRouteList,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"client_ip", c.ClientIP())

		// Parse query parameters for filters
		filters := make(map[string]interface{})
		if airportType := c.Query("airport_type"); airportType != "" {
			filters["airport_type"] = airportType
		}
		if originCountry := c.Query("origin_country"); originCountry != "" {
			filters["origin_country"] = originCountry
		}

		routes, err := h.RouteInteractor.ListRoutes(c.Request.Context(), filters)
		if err != nil {
			log.Error(logger.LogRouteListError,
				"error", err,
				"client_ip", c.ClientIP())
			h.Response.Error(c, domain.MsgServerError)
			return
		}

		// Convert to response with encoded IDs and HATEOAS links
		baseURL := GetBaseURL(c)
		response := ToRouteListResponse(routes, h.EncodeID, baseURL)

		log.Debug(logger.LogRouteListOK,
			"count", len(routes),
			"client_ip", c.ClientIP())

		h.Response.DataOnly(c, response)
	}
}
