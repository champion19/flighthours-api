package handlers

import (
	"github.com/champion19/flighthours-api/core/interactor"
	"github.com/champion19/flighthours-api/core/ports/input"
	"github.com/champion19/flighthours-api/platform/logger"
)

type handler struct {
	EmployeeService input.Service
	Interactor      *interactor.Interactor
	Logger          logger.Logger
}

func New(service input.Service, interactor *interactor.Interactor, logger logger.Logger) *handler {
	return &handler{
		EmployeeService: service,
		Interactor:      interactor,
	}
}
