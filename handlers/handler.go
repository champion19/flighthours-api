package handlers

import "github.com/champion19/flighthours-api/core/ports/input"

type handler struct {
	EmployeeService input.Service
}

func New(service input.Service) *handler {
	return &handler{
		EmployeeService: service,
	}
}
