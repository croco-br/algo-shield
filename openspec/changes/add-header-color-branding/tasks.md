# Implementation Tasks

## 1. Database Schema
- [x] 1.1 Create migration script to add `header_color` column to `branding_config` table (006_add_header_color.sql)
- [x] 1.2 Set default value for `header_color` to '#1e1e1e' (in migration)
- [x] 1.3 Update existing branding_config rows with default header_color value (in migration)
- [ ] 1.4 Test migration (up and down)

## 2. Backend API - Models
- [x] 2.1 Add `HeaderColor` field to `BrandingConfig` model in `pkg/models/branding.go`
- [x] 2.2 Add `HeaderColor` field to `UpdateBrandingRequest` in `internal/api/branding/models.go`
- [x] 2.3 Add validation struct tag for header_color (required, hexcolor)

## 3. Backend API - Repository
- [x] 3.1 Update repository Get() method to include header_color in SELECT query
- [x] 3.2 Update repository Update() method to include header_color in UPDATE query
- [x] 3.3 Ensure cache invalidation works correctly with new field (cache invalidation already handles all fields)

## 4. Backend API - Service Layer
- [x] 4.1 Add DefaultHeaderColor constant to service
- [x] 4.2 Update GetBranding() to return default header_color when config doesn't exist
- [x] 4.3 Add header_color validation in UpdateBranding() method
- [x] 4.4 Ensure header_color is included in returned config

## 5. Frontend - Store
- [x] 5.1 Add `header_color` field to BrandingConfig TypeScript interface
- [x] 5.2 Add `header_color` field to UpdateBrandingRequest interface
- [x] 5.3 Update applyBranding() function to set `--color-header-background` CSS custom property
- [x] 5.4 Update applyDefaultBranding() to include default header color

## 6. Frontend - Header Component
- [x] 6.1 Replace hardcoded background gradient with CSS custom property
- [x] 6.2 Use `--color-header-background` for header background color
- [x] 6.3 Ensure header updates dynamically when branding changes (via store watcher)

## 7. Frontend - Admin UI
- [x] 7.1 Add header color field to BrandingView form
- [x] 7.2 Add color picker for header color
- [x] 7.3 Update live preview to show header with custom color (added header preview section)
- [x] 7.4 Add header color to form validation (pattern validation in input)

## 8. CSS Custom Properties
- [x] 8.1 Define `--color-header-background` CSS custom property in global styles
- [x] 8.2 Update index.html inline script to set default header color
- [x] 8.3 Ensure header color is applied before FOUC (inline script + default in store)

## 9. Documentation
- [x] 9.1 Update README.md API documentation to include header_color field
- [x] 9.2 Document header_color validation rules and default value

## 10. Testing
- [ ] 10.1 Test header color updates dynamically after branding change
- [ ] 10.2 Test default header color is applied when no config exists
- [ ] 10.3 Test header color validation (invalid hex format)
- [ ] 10.4 Test header color in admin UI preview

