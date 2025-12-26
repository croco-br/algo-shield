package branding

// UpdateBrandingRequest represents a request to update branding configuration
type UpdateBrandingRequest struct {
	AppName        string  `json:"app_name" validate:"required,min=1,max=100"`
	IconURL        *string `json:"icon_url,omitempty" validate:"omitempty,uri|filepath"`
	FaviconURL     *string `json:"favicon_url,omitempty" validate:"omitempty,uri|filepath"`
	PrimaryColor   string  `json:"primary_color" validate:"required,hexcolor"`
	SecondaryColor string  `json:"secondary_color" validate:"required,hexcolor"`
	HeaderColor    string  `json:"header_color" validate:"required,hexcolor"`
}
