package handlers

import (
	"github.com/champion19/flighthours-api/core/interactor"
	"github.com/champion19/flighthours-api/core/interactor/services/domain"
	"github.com/champion19/flighthours-api/core/ports/input"
	"github.com/champion19/flighthours-api/platform/logger"
	"github.com/champion19/flighthours-api/tools/idencoder"
)

type handler struct {
	EmployeeService input.Service
	Interactor      *interactor.Interactor
	Logger          logger.Logger
	MessageManager  *domain.MessageManager // Add this field to access MessageManager
	IDEncoder       *idencoder.HashidsEncoder
}

func New(service input.Service, interactor *interactor.Interactor, logger logger.Logger, messageManager *domain.MessageManager, encoder *idencoder.HashidsEncoder) *handler {
	return &handler{
		EmployeeService: service,
		Interactor:      interactor,
		Logger:          logger,
		MessageManager:  messageManager,
		IDEncoder:       encoder,
	}
}
