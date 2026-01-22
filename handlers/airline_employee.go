package handlers

import (
	"time"

	domain "github.com/champion19/flighthours-api/core/interactor/services/domain"
)

// AirlineEmployeeResponse - Response DTO for airline employee data (HU26)
type AirlineEmployeeResponse struct {
	ID                   string    `json:"id"`
	Name                 string    `json:"name"`
	AirlineID            string    `json:"airline_id"`
	AirlineName          string    `json:"airline_name"`
	AirlineCode          string    `json:"airline_code"`
	Email                string    `json:"email"`
	IdentificationNumber string    `json:"identification_number"`
	Bp                   string    `json:"bp,omitempty"`
	StartDate            time.Time `json:"start_date"`
	EndDate              time.Time `json:"end_date"`
	Active               bool      `json:"active"`
	Role                 string    `json:"role"`
	Links                []Link    `json:"_links,omitempty"`
}

// FromDomainAirlineEmployee converts domain.AirlineEmployee to AirlineEmployeeResponse with encoded ID
func FromDomainAirlineEmployee(employee *domain.AirlineEmployee, encodedID, encodedAirlineID string) AirlineEmployeeResponse {
	return AirlineEmployeeResponse{
		ID:                   encodedID,
		Name:                 employee.Name,
		AirlineID:            encodedAirlineID,
		AirlineName:          employee.AirlineName,
		AirlineCode:          employee.AirlineCode,
		Email:                employee.Email,
		IdentificationNumber: employee.IdentificationNumber,
		Bp:                   employee.Bp,
		StartDate:            employee.StartDate,
		EndDate:              employee.EndDate,
		Active:               employee.Active,
		Role:                 employee.Role,
	}
}

// AirlineEmployeeRequest - Request DTO for creating/updating airline employee
type AirlineEmployeeRequest struct {
	Name                 string `json:"name" binding:"required"`
	AirlineID            string `json:"airline_id" binding:"required"`
	Email                string `json:"email" binding:"required,email"`
	IdentificationNumber string `json:"identification_number" binding:"required"`
	Bp                   string `json:"bp"`
	StartDate            string `json:"start_date" binding:"required"`
	EndDate              string `json:"end_date"`
	Active               bool   `json:"active"`
	Role                 string `json:"role" binding:"required"`
}

// Sanitize trims whitespace from AirlineEmployeeRequest fields
func (r *AirlineEmployeeRequest) Sanitize() {
	r.Name = TrimString(r.Name)
	r.AirlineID = TrimString(r.AirlineID)
	r.Email = TrimString(r.Email)
	r.IdentificationNumber = TrimString(r.IdentificationNumber)
	r.Bp = TrimString(r.Bp)
	r.StartDate = TrimString(r.StartDate)
	r.EndDate = TrimString(r.EndDate)
	r.Role = TrimString(r.Role)
}

// ToDomain converts AirlineEmployeeRequest to domain.AirlineEmployee
func (r *AirlineEmployeeRequest) ToDomain() (domain.AirlineEmployee, error) {
	startDate, err := time.Parse("2006-01-02", r.StartDate)
	if err != nil {
		return domain.AirlineEmployee{}, domain.ErrInvalidDateFormat
	}

	var endDate time.Time
	if r.EndDate != "" {
		endDate, err = time.Parse("2006-01-02", r.EndDate)
		if err != nil {
			return domain.AirlineEmployee{}, domain.ErrInvalidDateFormat
		}
	}

	// Validate start date is before end date
	if !endDate.IsZero() && startDate.After(endDate) {
		return domain.AirlineEmployee{}, domain.ErrStartDateAfterEndDate
	}

	return domain.AirlineEmployee{
		Name:                 r.Name,
		AirlineID:            r.AirlineID,
		Email:                r.Email,
		IdentificationNumber: r.IdentificationNumber,
		Bp:                   r.Bp,
		StartDate:            startDate,
		EndDate:              endDate,
		Active:               r.Active,
		Role:                 r.Role,
	}, nil
}

// AirlineEmployeeCreateResponse - Response DTO for created airline employee (HU28)
type AirlineEmployeeCreateResponse struct {
	ID    string `json:"id"`
	Links []Link `json:"_links,omitempty"`
}

// AirlineEmployeeStatusResponse - Response DTO for status update (HU29, HU30)
type AirlineEmployeeStatusResponse struct {
	ID      string `json:"id"`
	Active  bool   `json:"active"`
	Updated bool   `json:"updated"`
	Links   []Link `json:"_links,omitempty"`
}

// AirlineEmployeeListResponse - Response DTO for listing airline employees
type AirlineEmployeeListResponse struct {
	Employees []AirlineEmployeeResponse `json:"employees"`
	Total     int                       `json:"total"`
	Links     []Link                    `json:"_links,omitempty"`
}

// ToAirlineEmployeeListResponse converts a slice of domain.AirlineEmployee to AirlineEmployeeListResponse
func ToAirlineEmployeeListResponse(employees []domain.AirlineEmployee, encodeFunc func(string) (string, error), baseURL string) AirlineEmployeeListResponse {
	response := AirlineEmployeeListResponse{
		Employees: make([]AirlineEmployeeResponse, 0, len(employees)),
		Total:     len(employees),
	}

	for _, employee := range employees {
		encodedID, err := encodeFunc(employee.ID)
		if err != nil {
			encodedID = employee.ID
		}
		encodedAirlineID, err := encodeFunc(employee.AirlineID)
		if err != nil {
			encodedAirlineID = employee.AirlineID
		}

		empResp := AirlineEmployeeResponse{
			ID:                   encodedID,
			Name:                 employee.Name,
			AirlineID:            encodedAirlineID,
			AirlineName:          employee.AirlineName,
			AirlineCode:          employee.AirlineCode,
			Email:                employee.Email,
			IdentificationNumber: employee.IdentificationNumber,
			Bp:                   employee.Bp,
			StartDate:            employee.StartDate,
			EndDate:              employee.EndDate,
			Active:               employee.Active,
			Role:                 employee.Role,
		}
		// Add HATEOAS links to each employee
		if baseURL != "" {
			empResp.Links = BuildAirlineEmployeeLinks(baseURL, encodedID)
		}
		response.Employees = append(response.Employees, empResp)
	}

	// Add collection-level links
	if baseURL != "" {
		response.Links = BuildAirlineEmployeeListLinks(baseURL)
	}

	return response
}
