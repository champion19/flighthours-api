package handlers

import (
	"github.com/champion19/flighthours-api/core/interactor"
	domain "github.com/champion19/flighthours-api/core/interactor/services/domain"
	"github.com/champion19/flighthours-api/core/ports/input"
	"github.com/champion19/flighthours-api/middleware"
	"github.com/champion19/flighthours-api/platform/logger"
	"github.com/champion19/flighthours-api/tools/idencoder"
	"github.com/gin-gonic/gin"
)

type handler struct {
	EmployeeService   input.Service
	Interactor        *interactor.Interactor
	Logger            logger.Logger
	IDEncoder         *idencoder.HashidsEncoder
	Response          *middleware.ResponseHandler
	MessageInteractor *interactor.MessageInteractor
}

func New(
	service input.Service,
	interactor *interactor.Interactor,
	logger logger.Logger,
	idEncoder *idencoder.HashidsEncoder,
	response *middleware.ResponseHandler,
	messageInteractor *interactor.MessageInteractor) *handler {
	return &handler{
		EmployeeService:   service,
		Interactor:        interactor,
		Logger:            logger,
		IDEncoder:         idEncoder,
		Response:          response,
		MessageInteractor: messageInteractor,
	}
}
func (h *handler) EncodeID(uuid string) (string, error) {
	encodedID, err := h.IDEncoder.Encode(uuid)
	if err != nil {
		h.Logger.Error(logger.LogMessageIDEncodeError,
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
		h.Logger.Error(logger.LogMessageIDDecodeError,
			"encoded_id", encodedID,
			"error", err)
		return "", err
	}
	return uuid, nil
}

// HandleIDEncodingError maneja errores de ofuscamiento y envía respuesta apropiada
func (h *handler) HandleIDEncodingError(c *gin.Context, uuid string, err error) {
	h.Logger.Error(logger.LogMessageIDEncodeError,
		"uuid", uuid,
		"error", err,
		"client_ip", c.ClientIP())
	c.Error(domain.ErrInternalServer)
}

// HandleIDDecodingError maneja errores de desofuscamiento y envía respuesta apropiada
func (h *handler) HandleIDDecodingError(c *gin.Context, encodedID string, err error) {
	h.Logger.Error(logger.LogMessageIDDecodeError,
		"encoded_id", encodedID,
		"error", err,
		"client_ip", c.ClientIP())
	c.Error(domain.ErrInvalidID)
}
