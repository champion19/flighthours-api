package domain

import (
	"time"

	"github.com/google/uuid"
)

// DailyLogbook represents the daily logbook entity
type DailyLogbook struct {
	ID         string    `json:"id"`
	LogDate    time.Time `json:"log_date"`
	EmployeeID string    `json:"employee_id"`
	BookPage   *int      `json:"book_page,omitempty"`
	Status     bool      `json:"status"`
}

// SetID generates a new UUID for the logbook
func (d *DailyLogbook) SetID() {
	d.ID = uuid.New().String()
}

// ToLogger returns a slice of strings for logging purposes
func (d *DailyLogbook) ToLogger() []string {
	return []string{
		"id:" + d.ID,
		"employee_id:" + d.EmployeeID,
		"log_date:" + d.LogDate.Format("2006-01-02"),
	}
}

// IsActive returns true if the logbook is active
func (d *DailyLogbook) IsActive() bool {
	return d.Status
}

// StatusString returns "active" or "inactive" based on status
func (d *DailyLogbook) StatusString() string {
	if d.Status {
		return "active"
	}
	return "inactive"
}
