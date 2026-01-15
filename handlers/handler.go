package handlers

import (
	"github.com/champion19/flighthours-api/core/interactor"
	domain "github.com/champion19/flighthours-api/core/interactor/services/domain"
	"github.com/champion19/flighthours-api/core/ports/input"
	"github.com/champion19/flighthours-api/middleware"
	"github.com/champion19/flighthours-api/platform/cache/messaging"
	"github.com/champion19/flighthours-api/platform/logger"
	"github.com/champion19/flighthours-api/tools/idencoder"
	"github.com/gin-gonic/gin"
)

type handler struct {
	EmployeeService                input.Service
	Interactor                     *interactor.Interactor
	IDEncoder                      *idencoder.HashidsEncoder
	Response                       *middleware.ResponseHandler
	MessageInteractor              *interactor.MessageInteractor
	MessagingCache                 *messaging.MessageCache
	AirlineInteractor              *interactor.AirlineInteractor
	AirportInteractor              *interactor.AirportInteractor
	DailyLogbookInteractor         *interactor.DailyLogbookInteractor
	AircraftRegistrationInteractor *interactor.AircraftRegistrationInteractor
	AircraftModelInteractor        *interactor.AircraftModelInteractor
	RouteInteractor                *interactor.RouteInteractor
	AirlineRouteInteractor         *interactor.AirlineRouteInteractor
	DailyLogbookDetailInteractor   *interactor.DailyLogbookDetailInteractor
}

func New(
	service input.Service,
	interactor *interactor.Interactor,
	idEncoder *idencoder.HashidsEncoder,
	response *middleware.ResponseHandler,
	messageInteractor *interactor.MessageInteractor,
	messagingCache *messaging.MessageCache,
	airlineInteractor *interactor.AirlineInteractor,
	airportInteractor *interactor.AirportInteractor,
	dailyLogbookInteractor *interactor.DailyLogbookInteractor,
	aircraftRegistrationInteractor *interactor.AircraftRegistrationInteractor,
	aircraftModelInteractor *interactor.AircraftModelInteractor,
	routeInteractor *interactor.RouteInteractor,
	airlineRouteInteractor *interactor.AirlineRouteInteractor,
	dailyLogbookDetailInteractor *interactor.DailyLogbookDetailInteractor) *handler {
	return &handler{
		EmployeeService:                service,
		Interactor:                     interactor,
		IDEncoder:                      idEncoder,
		Response:                       response,
		MessageInteractor:              messageInteractor,
		MessagingCache:                 messagingCache,
		AirlineInteractor:              airlineInteractor,
		AirportInteractor:              airportInteractor,
		DailyLogbookInteractor:         dailyLogbookInteractor,
		AircraftRegistrationInteractor: aircraftRegistrationInteractor,
		AircraftModelInteractor:        aircraftModelInteractor,
		RouteInteractor:                routeInteractor,
		AirlineRouteInteractor:         airlineRouteInteractor,
		DailyLogbookDetailInteractor:   dailyLogbookDetailInteractor,
	}
}

var Logger = logger.NewSlogLogger()

func (h *handler) EncodeID(uuid string) (string, error) {
	encodedID, err := h.IDEncoder.Encode(uuid)
	if err != nil {
		Logger.Error(logger.LogMessageIDEncodeError,
			"uuid", uuid,
			"error", err)
		return "", err
	}
	return encodedID, nil
}

// DecodeID desofusca un ID ofuscado a UUID usando el encoder del handler
// Retorna el UUID o un error si falla
func (h *handler) DecodeID(encodedID string) (string, error) {
	uuid, err := h.IDEncoder.Decode(encodedID)
	if err != nil {
		Logger.Error(logger.LogMessageIDDecodeError,
			"encoded_id", encodedID,
			"error", err)
		return "", err
	}
	return uuid, nil
}

// HandleIDEncodingError maneja errores de ofuscamiento y envía respuesta apropiada
func (h *handler) HandleIDEncodingError(c *gin.Context, uuid string, err error) {
	Logger.Error(logger.LogMessageIDEncodeError,
		"uuid", uuid,
		"error", err,
		"client_ip", c.ClientIP())
	c.Error(domain.ErrInternalServer)
}

// HandleIDDecodingError maneja errores de desofuscamiento y envía respuesta apropiada
func (h *handler) HandleIDDecodingError(c *gin.Context, encodedID string, err error) {
	Logger.Error(logger.LogMessageIDDecodeError,
		"encoded_id", encodedID,
		"error", err,
		"client_ip", c.ClientIP())
	c.Error(domain.ErrInvalidID)
}

// resolveID accepts both raw UUID and obfuscated ID formats
// Returns (uuid, responseID) where:
// - uuid: the decoded UUID for internal use
// - responseID: the ID to use in the response (maintains consistency with input format)
// If the ID cannot be resolved, returns empty strings
func (h *handler) resolveID(inputID string) (string, string) {
	if inputID == "" {
		return "", ""
	}

	// Check if it's a valid UUID
	if isValidUUID(inputID) {
		// It's a direct UUID - encode it for consistent response format
		encodedID, err := h.EncodeID(inputID)
		if err != nil {
			Logger.Warn(logger.LogIDEncodeError, "uuid", inputID, "error", err)
			return inputID, inputID // Use raw UUID if encoding fails
		}
		return inputID, encodedID
	}

	// It's an obfuscated ID - decode it
	uuid, err := h.DecodeID(inputID)
	if err != nil {
		Logger.Warn(logger.LogMessageIDDecodeError, "encoded_id", inputID, "error", err)
		return "", "" // Return empty if decoding fails
	}
	return uuid, inputID // Return decoded UUID and keep original obfuscated ID for response
}

// isValidUUID verifica si un string es un UUID válido
func isValidUUID(str string) bool {
	// UUID tiene formato: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx (36 caracteres con guiones)
	if len(str) != 36 {
		return false
	}
	// Verificar posiciones de los guiones
	if str[8] != '-' || str[13] != '-' || str[18] != '-' || str[23] != '-' {
		return false
	}
	// Verificar que los demás caracteres sean hexadecimales
	for i, c := range str {
		if i == 8 || i == 13 || i == 18 || i == 23 {
			continue // Saltar guiones
		}
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
			return false
		}
	}
	return true
}
