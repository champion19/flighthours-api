package handlers

import (
	domain "github.com/champion19/flighthours-api/core/interactor/services/domain"
	"github.com/champion19/flighthours-api/middleware"
	"github.com/champion19/flighthours-api/platform/logger"
	"github.com/gin-gonic/gin"
)

// GetAircraftModelByID godoc
// @Summary      Get aircraft model by ID
// @Description  Returns aircraft model information by ID (accepts both UUID and obfuscated ID)
// @Tags         Aircraft Models
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Aircraft Model ID (obfuscated ID)"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /aircraft-models/{id} [get]
func (h *handler) GetAircraftModelByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		inputID := c.Param("id")
		if inputID == "" {
			log.Error(logger.LogMessageIDDecodeError, "error", "empty id parameter", "client_ip", c.ClientIP())
			h.Response.Error(c, domain.MsgValIDInvalid)
			return
		}

		log.Info(logger.LogAircraftModelGet, "input_id", inputID, "client_ip", c.ClientIP())

		// Resolve ID (accepts both UUID and obfuscated ID)
		modelUUID, responseID := h.resolveID(inputID)
		if modelUUID == "" {
			h.HandleIDDecodingError(c, inputID, domain.ErrInvalidID)
			return
		}

		// Get aircraft model from interactor
		model, err := h.AircraftModelInteractor.GetAircraftModelByID(c.Request.Context(), modelUUID)
		if err != nil {
			log.Error(logger.LogAircraftModelGetError, "aircraft_model_id", modelUUID, "error", err, "client_ip", c.ClientIP())
			if err == domain.ErrAircraftModelNotFound {
				h.Response.Error(c, domain.MsgAircraftModelNotFound)
				return
			}
			h.Response.Error(c, domain.MsgServerError)
			return
		}

		response := FromDomainAircraftModel(model, responseID)

		// Build HATEOAS links
		baseURL := GetBaseURL(c)
		response.Links = BuildAircraftModelLinks(baseURL, responseID)

		log.Success(logger.LogAircraftModelGetOK, model.ToLogger())
		h.Response.SuccessWithData(c, domain.MsgAircraftModelGetOK, response)
	}
}

// ListAircraftModels godoc
// @Summary      List all aircraft models
// @Description  Returns a list of all aircraft models with optional engine type filter
// @Tags         Aircraft Models
// @Produce      json
// @Param        engine_type query string false "Filter by engine type (e.g., JET, TUR)"
// @Success      200  {object}  AircraftModelListResponse
// @Failure      500  {object}  map[string]interface{}
// @Router       /aircraft-models [get]
func (h *handler) ListAircraftModels() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		log.Debug(logger.LogAircraftModelList,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"client_ip", c.ClientIP())

		// Parse query parameters for filters
		filters := make(map[string]interface{})
		if engineType := c.Query("engine_type"); engineType != "" {
			filters["engine_type"] = engineType
		}

		models, err := h.AircraftModelInteractor.ListAircraftModels(c.Request.Context(), filters)
		if err != nil {
			log.Error(logger.LogAircraftModelListError,
				"error", err,
				"client_ip", c.ClientIP())
			h.Response.Error(c, domain.MsgServerError)
			return
		}

		// Convert to response with encoded IDs and HATEOAS links
		baseURL := GetBaseURL(c)
		response := ToAircraftModelListResponse(models, h.EncodeID, baseURL)

		log.Debug(logger.LogAircraftModelListOK,
			"count", len(models),
			"client_ip", c.ClientIP())

		h.Response.DataOnly(c, response)
	}
}
