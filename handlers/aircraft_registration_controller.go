package handlers

import (
	"net/http"

	domain "github.com/champion19/flighthours-api/core/interactor/services/domain"
	"github.com/champion19/flighthours-api/middleware"
	"github.com/champion19/flighthours-api/platform/logger"
	"github.com/gin-gonic/gin"
)

// GetAircraftRegistrationByID godoc
// @Summary      Get aircraft registration by ID
// @Description  Returns aircraft registration information by ID (accepts both UUID and obfuscated ID)
// @Tags         Aircraft Registrations
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Aircraft Registration ID (UUID or obfuscated)"
// @Success      200  {object}  AircraftRegistrationResponse
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /aircraft-registrations/{id} [get]
func (h *handler) GetAircraftRegistrationByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		traceID := middleware.GetTraceIDFromContext(ctx)
		log := Logger.WithTraceID(traceID)

		inputID := c.Param("id")
		if inputID == "" {
			_ = c.Error(domain.ErrInvalidID)
			return
		}

		// Dual-format ID resolution: accept both raw UUID and obfuscated ID
		registrationUUID, responseID := h.resolveID(inputID)
		if registrationUUID == "" {
			log.Warn(logger.LogAircraftRegistrationNotFound, "id", inputID, "client_ip", c.ClientIP())
			_ = c.Error(domain.ErrAircraftRegistrationNotFound)
			return
		}

		registration, err := h.AircraftRegistrationInteractor.GetAircraftRegistrationByID(ctx, registrationUUID)
		if err != nil {
			log.Error(logger.LogAircraftRegistrationGetError, "id", inputID, "error", err, "client_ip", c.ClientIP())
			_ = c.Error(err)
			return
		}

		baseURL := GetBaseURL(c)
		response := FromDomainAircraftRegistration(registration, responseID)
		response.Links = BuildAircraftRegistrationLinks(baseURL, responseID)

		log.Success(logger.LogAircraftRegistrationGetOK, registration.ToLogger())
		c.JSON(http.StatusOK, response)
	}
}

// ListAircraftRegistrations godoc
// @Summary      List all aircraft registrations
// @Description  Returns a list of all aircraft registrations with optional airline_id filter
// @Tags         Aircraft Registrations
// @Produce      json
// @Param        airline_id query string false "Filter by airline ID"
// @Success      200  {object}  AircraftRegistrationListResponse
// @Failure      500  {object}  map[string]interface{}
// @Router       /aircraft-registrations [get]
func (h *handler) ListAircraftRegistrations() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		traceID := middleware.GetTraceIDFromContext(ctx)
		log := Logger.WithTraceID(traceID)

		filters := make(map[string]interface{})

		// Check for airline_id filter
		if airlineID := c.Query("airline_id"); airlineID != "" {
			// Resolve airline ID (could be obfuscated or raw UUID)
			resolvedAirlineID, _ := h.resolveID(airlineID)
			if resolvedAirlineID != "" {
				filters["airline_id"] = resolvedAirlineID
			}
		}

		registrations, err := h.AircraftRegistrationInteractor.ListAircraftRegistrations(ctx, filters)
		if err != nil {
			log.Error(logger.LogAircraftRegistrationListError, "error", err, "client_ip", c.ClientIP())
			_ = c.Error(err)
			return
		}

		baseURL := GetBaseURL(c)
		response := ToAircraftRegistrationListResponse(registrations, h.EncodeID, baseURL)

		log.Success(logger.LogAircraftRegistrationListOK, "count", len(registrations))
		c.JSON(http.StatusOK, response)
	}
}

// CreateAircraftRegistration godoc
// @Summary      Create a new aircraft registration
// @Description  Creates a new aircraft registration (license plate must be unique)
// @Tags         Aircraft Registrations
// @Accept       json
// @Produce      json
// @Param        request body CreateAircraftRegistrationRequest true "Aircraft Registration data"
// @Success      201  {object}  AircraftRegistrationResponse
// @Failure      400  {object}  map[string]interface{}
// @Failure      409  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /aircraft-registrations [post]
func (h *handler) CreateAircraftRegistration() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		traceID := middleware.GetTraceIDFromContext(ctx)
		log := Logger.WithTraceID(traceID)

		var req CreateAircraftRegistrationRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Warn(logger.LogRegJSONBindError, "error", err, "client_ip", c.ClientIP())
			_ = c.Error(domain.ErrInvalidJSONFormat)
			return
		}

		req.Sanitize()

		// Resolve foreign key IDs (could be obfuscated or raw UUID)
		resolvedModelID, _ := h.resolveID(req.AircraftModelID)
		if resolvedModelID == "" {
			log.Warn("Invalid aircraft_model_id format", "id", req.AircraftModelID, "client_ip", c.ClientIP())
			_ = c.Error(domain.ErrAircraftRegistrationInvalidModel)
			return
		}
		req.AircraftModelID = resolvedModelID

		resolvedAirlineID, _ := h.resolveID(req.AirlineID)
		if resolvedAirlineID == "" {
			log.Warn("Invalid airline_id format", "id", req.AirlineID, "client_ip", c.ClientIP())
			_ = c.Error(domain.ErrAircraftRegistrationInvalidAirline)
			return
		}
		req.AirlineID = resolvedAirlineID

		registration := req.ToDomain()

		if err := h.AircraftRegistrationInteractor.CreateAircraftRegistration(ctx, registration); err != nil {
			log.Error(logger.LogAircraftRegistrationCreateError, "error", err, "client_ip", c.ClientIP())
			_ = c.Error(err)
			return
		}

		// Fetch the created registration to get complete data with model_name and airline_name
		createdRegistration, err := h.AircraftRegistrationInteractor.GetAircraftRegistrationByID(ctx, registration.ID)
		if err != nil {
			log.Warn("Could not fetch created registration", "id", registration.ID, "error", err)
			// Fallback to basic response without names
			createdRegistration = &registration
		}

		encodedID, err := h.EncodeID(registration.ID)
		if err != nil {
			h.HandleIDEncodingError(c, registration.ID, err)
			return
		}

		baseURL := GetBaseURL(c)
		response := FromDomainAircraftRegistration(createdRegistration, encodedID)
		response.Links = BuildAircraftRegistrationCreatedLinks(baseURL, encodedID)

		SetLocationHeader(c, baseURL, "aircraft-registrations", encodedID)
		log.Success(logger.LogAircraftRegistrationCreateOK, createdRegistration.ToLogger())
		c.JSON(http.StatusCreated, response)
	}
}

// UpdateAircraftRegistration godoc
// @Summary      Update an existing aircraft registration
// @Description  Updates an aircraft registration by ID (accepts both UUID and obfuscated ID)
// @Tags         Aircraft Registrations
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Aircraft Registration ID (UUID or obfuscated)"
// @Param        request body UpdateAircraftRegistrationRequest true "Aircraft Registration data"
// @Success      200  {object}  AircraftRegistrationResponse
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /aircraft-registrations/{id} [put]
func (h *handler) UpdateAircraftRegistration() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		traceID := middleware.GetTraceIDFromContext(ctx)
		log := Logger.WithTraceID(traceID)

		inputID := c.Param("id")
		if inputID == "" {
			_ = c.Error(domain.ErrInvalidID)
			return
		}

		// Dual-format ID resolution: accept both raw UUID and obfuscated ID
		registrationUUID, responseID := h.resolveID(inputID)
		if registrationUUID == "" {
			log.Warn(logger.LogAircraftRegistrationNotFound, "id", inputID, "client_ip", c.ClientIP())
			_ = c.Error(domain.ErrAircraftRegistrationNotFound)
			return
		}

		var req UpdateAircraftRegistrationRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Warn(logger.LogRegJSONBindError, "error", err, "client_ip", c.ClientIP())
			_ = c.Error(domain.ErrInvalidJSONFormat)
			return
		}

		req.Sanitize()

		// Resolve foreign key IDs (could be obfuscated or raw UUID)
		resolvedModelID, _ := h.resolveID(req.AircraftModelID)
		if resolvedModelID == "" {
			log.Warn("Invalid aircraft_model_id format", "id", req.AircraftModelID, "client_ip", c.ClientIP())
			_ = c.Error(domain.ErrAircraftRegistrationInvalidModel)
			return
		}
		req.AircraftModelID = resolvedModelID

		resolvedAirlineID, _ := h.resolveID(req.AirlineID)
		if resolvedAirlineID == "" {
			log.Warn("Invalid airline_id format", "id", req.AirlineID, "client_ip", c.ClientIP())
			_ = c.Error(domain.ErrAircraftRegistrationInvalidAirline)
			return
		}
		req.AirlineID = resolvedAirlineID

		registration := req.ToDomain(registrationUUID)

		if err := h.AircraftRegistrationInteractor.UpdateAircraftRegistration(ctx, registration); err != nil {
			log.Error(logger.LogAircraftRegistrationUpdateError, "error", err, "client_ip", c.ClientIP())
			_ = c.Error(err)
			return
		}

		// Fetch the updated registration to get complete data with model_name and airline_name
		updatedRegistration, err := h.AircraftRegistrationInteractor.GetAircraftRegistrationByID(ctx, registration.ID)
		if err != nil {
			log.Warn("Could not fetch updated registration", "id", registration.ID, "error", err)
			// Fallback to basic response without names
			updatedRegistration = &registration
		}

		baseURL := GetBaseURL(c)
		response := FromDomainAircraftRegistration(updatedRegistration, responseID)
		response.Links = BuildAircraftRegistrationLinks(baseURL, responseID)

		log.Success(logger.LogAircraftRegistrationUpdateOK, updatedRegistration.ToLogger())
		c.JSON(http.StatusOK, response)
	}
}
