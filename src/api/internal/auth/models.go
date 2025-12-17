package auth

// LoginRequest represents a login request DTO
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email_regex,max=254"`
	Password string `json:"password" validate:"required,min=8"`
}

// RegisterRequest represents a registration request DTO
type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email_regex,max=254"`
	Password string `json:"password" validate:"required,password"`
	Name     string `json:"name" validate:"required,min=1,max=255"`
}
