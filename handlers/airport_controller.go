package handlers

import (
	domain "github.com/champion19/flighthours-api/core/interactor/services/domain"
	"github.com/champion19/flighthours-api/middleware"
	"github.com/champion19/flighthours-api/platform/logger"
	"github.com/gin-gonic/gin"
)

// GetAirportByID godoc
// @Summary      Get airport by ID
// @Description  Returns airport information by ID (accepts both UUID and obfuscated ID)
// @Tags         Airports
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Airport ID (UUID or obfuscated)"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /airports/{id} [get]
func (h *handler) GetAirportByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		inputID := c.Param("id")
		if inputID == "" {
			log.Error(logger.LogMessageIDDecodeError, "error", "empty id parameter", "client_ip", c.ClientIP())
			h.Response.Error(c, domain.MsgValIDInvalid)
			return
		}

		log.Info(logger.LogAirportGet, "input_id", inputID, "client_ip", c.ClientIP())

		var airportUUID string
		var responseID string

		// Detect if it's a valid UUID or obfuscated ID
		if isValidUUID(inputID) {
			// It's a direct UUID
			airportUUID = inputID
			// Encode UUID for response (maintain consistency)
			encodedID, err := h.EncodeID(inputID)
			if err != nil {
				log.Warn(logger.LogIDEncodeError, "uuid", inputID, "error", err)
				responseID = inputID
			} else {
				responseID = encodedID
			}
			log.Debug(logger.LogAirportGet, "detected_format", "UUID", "uuid", airportUUID)
		} else {
			// It's an obfuscated ID, decode it
			uuid, err := h.DecodeID(inputID)
			if err != nil {
				h.HandleIDDecodingError(c, inputID, err)
				return
			}
			airportUUID = uuid
			responseID = inputID
			log.Debug(logger.LogAirportGet, "detected_format", "encoded", "decoded_uuid", airportUUID)
		}

		// Get airport from interactor
		airport, err := h.AirportInteractor.GetAirportByID(c.Request.Context(), airportUUID)
		if err != nil {
			log.Error(logger.LogAirportGetError, "airport_id", airportUUID, "error", err, "client_ip", c.ClientIP())
			if err == domain.ErrAirportNotFound {
				h.Response.Error(c, domain.MsgAirportNotFound)
				return
			}
			h.Response.Error(c, domain.MsgServerError)
			return
		}

		response := FromDomainAirport(airport, responseID)
		log.Success(logger.LogAirportGetOK, airport.ToLogger())
		h.Response.SuccessWithData(c, domain.MsgAirportGetOK, response)
	}
}

// ActivateAirport godoc
// @Summary      Activate airport
// @Description  Sets airport status to active (accepts both UUID and obfuscated ID)
// @Tags         Airports
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Airport ID (UUID or obfuscated)"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /airports/{id}/activate [patch]
func (h *handler) ActivateAirport() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		inputID := c.Param("id")
		if inputID == "" {
			log.Error(logger.LogMessageIDDecodeError, "error", "empty id parameter", "client_ip", c.ClientIP())
			h.Response.Error(c, domain.MsgValIDInvalid)
			return
		}

		log.Info(logger.LogAirportActivate, "input_id", inputID, "client_ip", c.ClientIP())

		var airportUUID string
		var responseID string

		// Detect if it's a valid UUID or obfuscated ID
		if isValidUUID(inputID) {
			airportUUID = inputID
			encodedID, err := h.EncodeID(inputID)
			if err != nil {
				log.Warn(logger.LogIDEncodeError, "uuid", inputID, "error", err)
				responseID = inputID
			} else {
				responseID = encodedID
			}
		} else {
			uuid, err := h.DecodeID(inputID)
			if err != nil {
				h.HandleIDDecodingError(c, inputID, err)
				return
			}
			airportUUID = uuid
			responseID = inputID
		}

		// Activate airport via interactor
		if err := h.AirportInteractor.ActivateAirport(c.Request.Context(), airportUUID); err != nil {
			log.Error(logger.LogAirportActivateError, "airport_id", airportUUID, "error", err, "client_ip", c.ClientIP())
			if err == domain.ErrAirportNotFound {
				h.Response.Error(c, domain.MsgAirportNotFound)
				return
			}
			h.Response.Error(c, domain.MsgServerError)
			return
		}

		response := AirportStatusResponse{
			ID:      responseID,
			Status:  "active",
			Updated: true,
		}

		log.Success(logger.LogAirportActivateOK, "airport_id", airportUUID, "client_ip", c.ClientIP())
		h.Response.SuccessWithData(c, domain.MsgAirportActivateOK, response)
	}
}

// DeactivateAirport godoc
// @Summary      Deactivate airport
// @Description  Sets airport status to inactive (accepts both UUID and obfuscated ID)
// @Tags         Airports
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Airport ID (UUID or obfuscated)"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /airports/{id}/deactivate [patch]
func (h *handler) DeactivateAirport() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		inputID := c.Param("id")
		if inputID == "" {
			log.Error(logger.LogMessageIDDecodeError, "error", "empty id parameter", "client_ip", c.ClientIP())
			h.Response.Error(c, domain.MsgValIDInvalid)
			return
		}

		log.Info(logger.LogAirportDeactivate, "input_id", inputID, "client_ip", c.ClientIP())

		var airportUUID string
		var responseID string

		// Detect if it's a valid UUID or obfuscated ID
		if isValidUUID(inputID) {
			airportUUID = inputID
			encodedID, err := h.EncodeID(inputID)
			if err != nil {
				log.Warn(logger.LogIDEncodeError, "uuid", inputID, "error", err)
				responseID = inputID
			} else {
				responseID = encodedID
			}
		} else {
			uuid, err := h.DecodeID(inputID)
			if err != nil {
				h.HandleIDDecodingError(c, inputID, err)
				return
			}
			airportUUID = uuid
			responseID = inputID
		}

		// Deactivate airport via interactor
		if err := h.AirportInteractor.DeactivateAirport(c.Request.Context(), airportUUID); err != nil {
			log.Error(logger.LogAirportDeactivateError, "airport_id", airportUUID, "error", err, "client_ip", c.ClientIP())
			if err == domain.ErrAirportNotFound {
				h.Response.Error(c, domain.MsgAirportNotFound)
				return
			}
			h.Response.Error(c, domain.MsgServerError)
			return
		}

		response := AirportStatusResponse{
			ID:      responseID,
			Status:  "inactive",
			Updated: true,
		}

		log.Success(logger.LogAirportDeactivateOK, "airport_id", airportUUID, "client_ip", c.ClientIP())
		h.Response.SuccessWithData(c, domain.MsgAirportDeactivateOK, response)
	}
}
