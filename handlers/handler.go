package handlers

import "github.com/champion19/flighthours-api/core/ports"

type handler struct {
	EmployeeService ports.Service
}

func New(service ports.Service) *handler {
	return &handler{
		EmployeeService: service,
	}
}
