package permissions

import (
	"github.com/google/uuid"
)

type AssignRoleRequest struct {
	RoleID uuid.UUID `json:"role_id"`
}

type UpdateUserActiveRequest struct {
	Active bool `json:"active"`
}
