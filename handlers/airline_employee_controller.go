package handlers

import (
	domain "github.com/champion19/flighthours-api/core/interactor/services/domain"
	"github.com/champion19/flighthours-api/middleware"
	"github.com/champion19/flighthours-api/platform/logger"
	"github.com/gin-gonic/gin"
)

// GetAirlineEmployeeByID godoc
// @Summary      Get airline employee by ID
// @Description  Returns airline employee information by ID
// @Tags         AirlineEmployees
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Airline Employee ID (obfuscated ID)"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /airline-employees/{id} [get]
func (h *handler) GetAirlineEmployeeByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		inputID := c.Param("id")
		if inputID == "" {
			log.Error(logger.LogMessageIDDecodeError, "error", "empty id parameter", "client_ip", c.ClientIP())
			h.Response.Error(c, domain.MsgValIDInvalid)
			return
		}

		log.Info(logger.LogDatabaseAvailable, "operation", "get_airline_employee", "input_id", inputID, "client_ip", c.ClientIP())

		// Resolve ID
		employeeUUID, responseID := h.resolveID(inputID)
		if employeeUUID == "" {
			h.HandleIDDecodingError(c, inputID, domain.ErrInvalidID)
			return
		}

		// Get airline employee from interactor
		employee, err := h.AirlineEmployeeInteractor.GetAirlineEmployeeByID(c.Request.Context(), employeeUUID)
		if err != nil {
			log.Error(logger.LogDatabaseUnavailable, "operation", "get_airline_employee", "employee_id", employeeUUID, "error", err, "client_ip", c.ClientIP())
			if err == domain.ErrAirlineEmployeeNotFound {
				h.Response.Error(c, domain.MsgAirlineEmployeeNotFound)
				return
			}
			h.Response.Error(c, domain.MsgServerError)
			return
		}

		// Encode airline ID
		encodedAirlineID, _ := h.EncodeID(employee.AirlineID)
		response := FromDomainAirlineEmployee(employee, responseID, encodedAirlineID)

		// Build HATEOAS links
		baseURL := GetBaseURL(c)
		response.Links = BuildAirlineEmployeeLinks(baseURL, responseID)

		log.Success(logger.LogDatabaseAvailable, employee.ToLogger())
		h.Response.SuccessWithData(c, domain.MsgAirlineEmployeeGetOK, response)
	}
}

// CreateAirlineEmployee godoc
// @Summary      Create airline employee
// @Description  Creates a new airline employee
// @Tags         AirlineEmployees
// @Accept       json
// @Produce      json
// @Param        request body AirlineEmployeeRequest true "Airline Employee data"
// @Success      201  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /airline-employees [post]
func (h *handler) CreateAirlineEmployee() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		log.Info(logger.LogDatabaseAvailable, "operation", "create_airline_employee", "client_ip", c.ClientIP())

		var req AirlineEmployeeRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Error(logger.LogMessageIDDecodeError, "error", err, "client_ip", c.ClientIP())
			h.Response.Error(c, domain.MsgValInvalidReq)
			return
		}

		req.Sanitize()

		// Decode airline ID
		airlineUUID, _ := h.resolveID(req.AirlineID)
		if airlineUUID == "" {
			h.Response.Error(c, domain.MsgAirlineEmployeeInvalidAirline)
			return
		}
		req.AirlineID = airlineUUID

		// Convert to domain
		employee, err := req.ToDomain()
		if err != nil {
			log.Error(logger.LogDatabaseUnavailable, "operation", "create_airline_employee", "error", err, "client_ip", c.ClientIP())
			if err == domain.ErrInvalidDateFormat {
				h.Response.Error(c, domain.MsgValInvalidDateFormat)
				return
			}
			if err == domain.ErrStartDateAfterEndDate {
				h.Response.Error(c, domain.MsgValStartDateAfterEndDate)
				return
			}
			h.Response.Error(c, domain.MsgValInvalidReq)
			return
		}

		// Create via interactor
		created, err := h.AirlineEmployeeInteractor.CreateAirlineEmployee(c.Request.Context(), employee)
		if err != nil {
			log.Error(logger.LogDatabaseUnavailable, "operation", "create_airline_employee", "error", err, "client_ip", c.ClientIP())
			if err == domain.ErrDuplicateUser {
				h.Response.Error(c, domain.MsgAirlineEmployeeDuplicate)
				return
			}
			if err == domain.ErrInvalidForeignKey {
				h.Response.Error(c, domain.MsgAirlineEmployeeInvalidAirline)
				return
			}
			h.Response.Error(c, domain.MsgAirlineEmployeeSaveError)
			return
		}

		// Encode ID for response
		encodedID, _ := h.EncodeID(created.ID)

		response := AirlineEmployeeCreateResponse{
			ID: encodedID,
		}

		// Build HATEOAS links
		baseURL := GetBaseURL(c)
		response.Links = BuildAirlineEmployeeLinks(baseURL, encodedID)

		log.Success(logger.LogDatabaseAvailable, created.ToLogger())
		h.Response.SuccessWithData(c, domain.MsgAirlineEmployeeCreated, response)
	}
}

// UpdateAirlineEmployee godoc
// @Summary      Update airline employee
// @Description  Updates an existing airline employee
// @Tags         AirlineEmployees
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Airline Employee ID (obfuscated ID)"
// @Param        request body AirlineEmployeeRequest true "Airline Employee data"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /airline-employees/{id} [put]
func (h *handler) UpdateAirlineEmployee() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		inputID := c.Param("id")
		if inputID == "" {
			log.Error(logger.LogMessageIDDecodeError, "error", "empty id parameter", "client_ip", c.ClientIP())
			h.Response.Error(c, domain.MsgValIDInvalid)
			return
		}

		log.Info(logger.LogDatabaseAvailable, "operation", "update_airline_employee", "input_id", inputID, "client_ip", c.ClientIP())

		// Resolve ID
		employeeUUID, responseID := h.resolveID(inputID)
		if employeeUUID == "" {
			h.HandleIDDecodingError(c, inputID, domain.ErrInvalidID)
			return
		}

		var req AirlineEmployeeRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Error(logger.LogMessageIDDecodeError, "error", err, "client_ip", c.ClientIP())
			h.Response.Error(c, domain.MsgValInvalidReq)
			return
		}

		req.Sanitize()

		// Decode airline ID
		airlineUUID, _ := h.resolveID(req.AirlineID)
		if airlineUUID == "" {
			h.Response.Error(c, domain.MsgAirlineEmployeeInvalidAirline)
			return
		}
		req.AirlineID = airlineUUID

		// Convert to domain
		employee, err := req.ToDomain()
		if err != nil {
			log.Error(logger.LogDatabaseUnavailable, "operation", "update_airline_employee", "error", err, "client_ip", c.ClientIP())
			if err == domain.ErrInvalidDateFormat {
				h.Response.Error(c, domain.MsgValInvalidDateFormat)
				return
			}
			if err == domain.ErrStartDateAfterEndDate {
				h.Response.Error(c, domain.MsgValStartDateAfterEndDate)
				return
			}
			h.Response.Error(c, domain.MsgValInvalidReq)
			return
		}

		// Update via interactor
		if err := h.AirlineEmployeeInteractor.UpdateAirlineEmployee(c.Request.Context(), employeeUUID, employee); err != nil {
			log.Error(logger.LogDatabaseUnavailable, "operation", "update_airline_employee", "employee_id", employeeUUID, "error", err, "client_ip", c.ClientIP())
			if err == domain.ErrAirlineEmployeeNotFound {
				h.Response.Error(c, domain.MsgAirlineEmployeeNotFound)
				return
			}
			if err == domain.ErrInvalidForeignKey {
				h.Response.Error(c, domain.MsgAirlineEmployeeInvalidAirline)
				return
			}
			h.Response.Error(c, domain.MsgAirlineEmployeeUpdateError)
			return
		}

		// Get updated employee
		updated, _ := h.AirlineEmployeeInteractor.GetAirlineEmployeeByID(c.Request.Context(), employeeUUID)
		encodedAirlineID, _ := h.EncodeID(updated.AirlineID)
		response := FromDomainAirlineEmployee(updated, responseID, encodedAirlineID)

		// Build HATEOAS links
		baseURL := GetBaseURL(c)
		response.Links = BuildAirlineEmployeeLinks(baseURL, responseID)

		log.Success(logger.LogDatabaseAvailable, updated.ToLogger())
		h.Response.SuccessWithData(c, domain.MsgAirlineEmployeeUpdated, response)
	}
}

// ActivateAirlineEmployee godoc
// @Summary      Activate airline employee
// @Description  Sets airline employee status to active
// @Tags         AirlineEmployees
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Airline Employee ID (obfuscated ID)"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /airline-employees/{id}/activate [patch]
func (h *handler) ActivateAirlineEmployee() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		inputID := c.Param("id")
		if inputID == "" {
			log.Error(logger.LogMessageIDDecodeError, "error", "empty id parameter", "client_ip", c.ClientIP())
			h.Response.Error(c, domain.MsgValIDInvalid)
			return
		}

		log.Info(logger.LogDatabaseAvailable, "operation", "activate_airline_employee", "input_id", inputID, "client_ip", c.ClientIP())

		// Resolve ID
		employeeUUID, responseID := h.resolveID(inputID)
		if employeeUUID == "" {
			h.HandleIDDecodingError(c, inputID, domain.ErrInvalidID)
			return
		}

		// Activate via interactor
		if err := h.AirlineEmployeeInteractor.ActivateAirlineEmployee(c.Request.Context(), employeeUUID); err != nil {
			log.Error(logger.LogDatabaseUnavailable, "operation", "activate_airline_employee", "employee_id", employeeUUID, "error", err, "client_ip", c.ClientIP())
			if err == domain.ErrAirlineEmployeeNotFound {
				h.Response.Error(c, domain.MsgAirlineEmployeeNotFound)
				return
			}
			h.Response.Error(c, domain.MsgAirlineEmployeeActivateErr)
			return
		}

		response := AirlineEmployeeStatusResponse{
			ID:      responseID,
			Active:  true,
			Updated: true,
		}

		// Build HATEOAS links
		baseURL := GetBaseURL(c)
		response.Links = BuildAirlineEmployeeStatusLinks(baseURL, responseID, true)

		log.Success(logger.LogDatabaseAvailable, "operation", "activate_airline_employee", "employee_id", employeeUUID, "client_ip", c.ClientIP())
		h.Response.SuccessWithData(c, domain.MsgAirlineEmployeeActivateOK, response)
	}
}

// DeactivateAirlineEmployee godoc
// @Summary      Deactivate airline employee
// @Description  Sets airline employee status to inactive
// @Tags         AirlineEmployees
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Airline Employee ID (obfuscated ID)"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /airline-employees/{id}/deactivate [patch]
func (h *handler) DeactivateAirlineEmployee() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		inputID := c.Param("id")
		if inputID == "" {
			log.Error(logger.LogMessageIDDecodeError, "error", "empty id parameter", "client_ip", c.ClientIP())
			h.Response.Error(c, domain.MsgValIDInvalid)
			return
		}

		log.Info(logger.LogDatabaseAvailable, "operation", "deactivate_airline_employee", "input_id", inputID, "client_ip", c.ClientIP())

		// Resolve ID
		employeeUUID, responseID := h.resolveID(inputID)
		if employeeUUID == "" {
			h.HandleIDDecodingError(c, inputID, domain.ErrInvalidID)
			return
		}

		// Deactivate via interactor
		if err := h.AirlineEmployeeInteractor.DeactivateAirlineEmployee(c.Request.Context(), employeeUUID); err != nil {
			log.Error(logger.LogDatabaseUnavailable, "operation", "deactivate_airline_employee", "employee_id", employeeUUID, "error", err, "client_ip", c.ClientIP())
			if err == domain.ErrAirlineEmployeeNotFound {
				h.Response.Error(c, domain.MsgAirlineEmployeeNotFound)
				return
			}
			h.Response.Error(c, domain.MsgAirlineEmployeeDeactivateErr)
			return
		}

		response := AirlineEmployeeStatusResponse{
			ID:      responseID,
			Active:  false,
			Updated: true,
		}

		// Build HATEOAS links
		baseURL := GetBaseURL(c)
		response.Links = BuildAirlineEmployeeStatusLinks(baseURL, responseID, false)

		log.Success(logger.LogDatabaseAvailable, "operation", "deactivate_airline_employee", "employee_id", employeeUUID, "client_ip", c.ClientIP())
		h.Response.SuccessWithData(c, domain.MsgAirlineEmployeeDeactivateOK, response)
	}
}

// ListAirlineEmployees godoc
// @Summary      List all airline employees
// @Description  Returns a list of all airline employees with optional filters
// @Tags         AirlineEmployees
// @Produce      json
// @Param        airline_id query string false "Filter by airline ID (obfuscated)"
// @Param        active query string false "Filter by active status (true/false)"
// @Success      200  {object}  AirlineEmployeeListResponse
// @Failure      500  {object}  map[string]interface{}
// @Router       /airline-employees [get]
func (h *handler) ListAirlineEmployees() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		log.Debug(logger.LogDatabaseAvailable, "operation", "list_airline_employees",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"client_ip", c.ClientIP())

		// Parse query parameters for filters
		filters := make(map[string]interface{})

		if airlineID := c.Query("airline_id"); airlineID != "" {
			// Decode airline ID
			airlineUUID, _ := h.resolveID(airlineID)
			if airlineUUID != "" {
				filters["airline_id"] = airlineUUID
			}
		}

		if active := c.Query("active"); active != "" {
			if active == "true" || active == "1" {
				filters["active"] = true
			} else if active == "false" || active == "0" {
				filters["active"] = false
			}
		}

		employees, err := h.AirlineEmployeeInteractor.ListAirlineEmployees(c.Request.Context(), filters)
		if err != nil {
			log.Error(logger.LogDatabaseUnavailable, "operation", "list_airline_employees",
				"error", err,
				"client_ip", c.ClientIP())
			h.Response.Error(c, domain.MsgAirlineEmployeeListError)
			return
		}

		// Convert to response with encoded IDs and HATEOAS links
		baseURL := GetBaseURL(c)
		response := ToAirlineEmployeeListResponse(employees, h.EncodeID, baseURL)

		log.Debug(logger.LogDatabaseAvailable, "operation", "list_airline_employees",
			"count", len(employees),
			"client_ip", c.ClientIP())

		h.Response.DataOnly(c, response)
	}
}
