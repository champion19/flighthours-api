package dto

import (
	"github.com/champion19/flighthours-api/core/domain"

)
type RegisterEmployee struct {
	Employee domain.Employee
	Message string
}

type UserSyncStatus struct {
	EmployeeID     string `json:"employee_id"`
	KeycloakUserID string `json:"keycloak_user_id"`
	IsSynced       bool   `json:"is_synced"`
	LastSyncAt     string `json:"last_sync_at"`
}
