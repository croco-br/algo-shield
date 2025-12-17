package roles

import (
	"github.com/google/uuid"
)

type AssignRoleRequest struct {
	RoleID uuid.UUID `json:"role_id"`
}
