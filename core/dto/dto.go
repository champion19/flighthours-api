package dto

import (
	"github.com/Nerzal/gocloak/v13"
	"github.com/champion19/Flighthours_backend/core/domain"
)
type RegisterEmployee struct {
	Employee domain.Employee
	Token    *gocloak.JWT
}

type UserSyncStatus struct {
	EmployeeID     string `json:"employee_id"`
	KeycloakUserID string `json:"keycloak_user_id"`
	IsSynced       bool   `json:"is_synced"`
	LastSyncAt     string `json:"last_sync_at"`
}
