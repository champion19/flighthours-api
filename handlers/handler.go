package handlers

import (
  "github.com/champion19/flighthours-api/core/interactor"
	"github.com/champion19/flighthours-api/core/ports/input"
)

type handler struct {
	EmployeeService input.Service
	Interactor      *interactor.Interactor
}

func New(service input.Service, interactor *interactor.Interactor) *handler {
	return &handler{
		EmployeeService: service,
		Interactor:      interactor,
	}
}
