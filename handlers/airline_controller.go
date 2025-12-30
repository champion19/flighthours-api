package handlers

import (
	domain "github.com/champion19/flighthours-api/core/interactor/services/domain"
	"github.com/champion19/flighthours-api/middleware"
	"github.com/champion19/flighthours-api/platform/logger"
	"github.com/gin-gonic/gin"
)

// GetAirlineByID godoc
// @Summary      Get airline by ID
// @Description  Returns airline information by ID (accepts both UUID and obfuscated ID)
// @Tags         Airlines
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Airline ID (UUID or obfuscated)"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /airlines/{id} [get]
func (h *handler) GetAirlineByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		inputID := c.Param("id")
		if inputID == "" {
			log.Error(logger.LogMessageIDDecodeError, "error", "empty id parameter", "client_ip", c.ClientIP())
			h.Response.Error(c, domain.MsgValIDInvalid)
			return
		}

		log.Info(logger.LogAirlineGet, "input_id", inputID, "client_ip", c.ClientIP())

		var airlineUUID string
		var responseID string

		// Detect if it's a valid UUID or obfuscated ID
		if isValidUUID(inputID) {
			// It's a direct UUID
			airlineUUID = inputID
			// Encode UUID for response (maintain consistency)
			encodedID, err := h.EncodeID(inputID)
			if err != nil {
				log.Warn(logger.LogIDEncodeError, "uuid", inputID, "error", err)
				responseID = inputID
			} else {
				responseID = encodedID
			}
			log.Debug(logger.LogAirlineGet, "detected_format", "UUID", "uuid", airlineUUID)
		} else {
			// It's an obfuscated ID, decode it
			uuid, err := h.DecodeID(inputID)
			if err != nil {
				h.HandleIDDecodingError(c, inputID, err)
				return
			}
			airlineUUID = uuid
			responseID = inputID
			log.Debug(logger.LogAirlineGet, "detected_format", "encoded", "decoded_uuid", airlineUUID)
		}

		// Get airline from interactor
		airline, err := h.AirlineInteractor.GetAirlineByID(c.Request.Context(), airlineUUID)
		if err != nil {
			log.Error(logger.LogAirlineGetError, "airline_id", airlineUUID, "error", err, "client_ip", c.ClientIP())
			if err == domain.ErrAirlineNotFound {
				h.Response.Error(c, domain.MsgAirlineNotFound)
				return
			}
			h.Response.Error(c, domain.MsgServerError)
			return
		}

		response := FromDomainAirline(airline, responseID)
		log.Success(logger.LogAirlineGetOK, airline.ToLogger())
		h.Response.SuccessWithData(c, domain.MsgAirlineGetOK, response)
	}
}

// ActivateAirline godoc
// @Summary      Activate airline
// @Description  Sets airline status to active (accepts both UUID and obfuscated ID)
// @Tags         Airlines
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Airline ID (UUID or obfuscated)"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /airlines/{id}/activate [patch]
func (h *handler) ActivateAirline() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		inputID := c.Param("id")
		if inputID == "" {
			log.Error(logger.LogMessageIDDecodeError, "error", "empty id parameter", "client_ip", c.ClientIP())
			h.Response.Error(c, domain.MsgValIDInvalid)
			return
		}

		log.Info(logger.LogAirlineActivate, "input_id", inputID, "client_ip", c.ClientIP())

		var airlineUUID string
		var responseID string

		// Detect if it's a valid UUID or obfuscated ID
		if isValidUUID(inputID) {
			airlineUUID = inputID
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
			airlineUUID = uuid
			responseID = inputID
		}

		// Activate airline via interactor
		if err := h.AirlineInteractor.ActivateAirline(c.Request.Context(), airlineUUID); err != nil {
			log.Error(logger.LogAirlineActivateError, "airline_id", airlineUUID, "error", err, "client_ip", c.ClientIP())
			if err == domain.ErrAirlineNotFound {
				h.Response.Error(c, domain.MsgAirlineNotFound)
				return
			}
			h.Response.Error(c, domain.MsgServerError)
			return
		}

		response := AirlineStatusResponse{
			ID:      responseID,
			Status:  "active",
			Updated: true,
		}

		log.Success(logger.LogAirlineActivateOK, "airline_id", airlineUUID, "client_ip", c.ClientIP())
		h.Response.SuccessWithData(c, domain.MsgAirlineActivateOK, response)
	}
}

// DeactivateAirline godoc
// @Summary      Deactivate airline
// @Description  Sets airline status to inactive (accepts both UUID and obfuscated ID)
// @Tags         Airlines
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Airline ID (UUID or obfuscated)"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /airlines/{id}/deactivate [patch]
func (h *handler) DeactivateAirline() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		inputID := c.Param("id")
		if inputID == "" {
			log.Error(logger.LogMessageIDDecodeError, "error", "empty id parameter", "client_ip", c.ClientIP())
			h.Response.Error(c, domain.MsgValIDInvalid)
			return
		}

		log.Info(logger.LogAirlineDeactivate, "input_id", inputID, "client_ip", c.ClientIP())

		var airlineUUID string
		var responseID string

		// Detect if it's a valid UUID or obfuscated ID
		if isValidUUID(inputID) {
			airlineUUID = inputID
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
			airlineUUID = uuid
			responseID = inputID
		}

		// Deactivate airline via interactor
		if err := h.AirlineInteractor.DeactivateAirline(c.Request.Context(), airlineUUID); err != nil {
			log.Error(logger.LogAirlineDeactivateError, "airline_id", airlineUUID, "error", err, "client_ip", c.ClientIP())
			if err == domain.ErrAirlineNotFound {
				h.Response.Error(c, domain.MsgAirlineNotFound)
				return
			}
			h.Response.Error(c, domain.MsgServerError)
			return
		}

		response := AirlineStatusResponse{
			ID:      responseID,
			Status:  "inactive",
			Updated: true,
		}

		log.Success(logger.LogAirlineDeactivateOK, "airline_id", airlineUUID, "client_ip", c.ClientIP())
		h.Response.SuccessWithData(c, domain.MsgAirlineDeactivateOK, response)
	}
}
