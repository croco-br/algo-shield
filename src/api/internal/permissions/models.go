package permissions

import (
	"github.com/google/uuid"
)

type AssignRoleRequest struct {
	RoleID uuid.UUID `json:"role_id" validate:"required,uuid"`
}

type UpdateUserActiveRequest struct {
	Active bool `json:"active" validate:"required"`
}
