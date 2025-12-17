package permissions

type UpdateUserActiveRequest struct {
	Active bool `json:"active" validate:"required"`
}
