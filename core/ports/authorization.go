package ports

import (
	"context"

	"github.com/Nerzal/gocloak/v13"
	"github.com/champion19/Flighthours_backend/core/domain"
)

type AuthorizationService interface {
	LoginUser(ctx context.Context, email string, password string) (*gocloak.JWT, error)

	SyncUserToKeycloak(ctx context.Context, employee *domain.Employee) (string, error)
	DeleteUserFromKeycloak(ctx context.Context, KeycloakUserID string) error

	SetUserPassword(ctx context.Context, KeycloakUserID string, password string) error

	AssignRole(ctx context.Context, employeeID string, roleName string) error
	RemoveRole(ctx context.Context, employeeID string, roleName string) error
	GetUserRoles(ctx context.Context, employeeID string) ([]string, error)

	HasRole(ctx context.Context, employeeID string, roleName string) (bool, error)
	HasPermission(ctx context.Context, employeeID string, resource, action string) (bool, error)

	CreateRole(ctx context.Context, roleName, description string) error
	GetAllRoles(ctx context.Context) ([]*string, error)
}


