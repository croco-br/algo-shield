# Implementation Tasks

## 1. Database Schema
- [ ] 1.1 Create migration script for `branding_config` table with columns: id, app_name, icon_url, primary_color, secondary_color, favicon_url, created_at, updated_at
- [ ] 1.2 Add unique constraint to ensure single branding configuration row
- [ ] 1.3 Insert default branding values in migration
- [ ] 1.4 Test migration (up and down)

## 2. Backend API - Models and Repository
- [ ] 2.1 Create `BrandingConfig` model in `internal/api/branding/model.go`
- [ ] 2.2 Create repository interface in `internal/api/branding/repository.go`
- [ ] 2.3 Implement PostgreSQL repository with methods: Get(), Update(), Create()
- [ ] 2.4 Add connection pooling for database access
- [ ] 2.5 Add validation struct tags for branding fields

## 3. Backend API - Service Layer
- [ ] 3.1 Create branding service in `internal/api/branding/service.go`
- [ ] 3.2 Implement GetBranding() method with default fallback
- [ ] 3.3 Implement UpdateBranding() method with validation
- [ ] 3.4 Add color format validation (hex: #RGB or #RRGGBB)
- [ ] 3.5 Add URL validation for icon_url and favicon_url
- [ ] 3.6 Add app_name length validation (max 100 chars)

## 4. Backend API - HTTP Handlers
- [ ] 4.1 Create handler in `internal/api/branding/handler.go`
- [ ] 4.2 Implement GET /api/v1/branding endpoint (public, no auth required)
- [ ] 4.3 Implement PUT /api/v1/branding endpoint (admin only, JWT + RBAC)
- [ ] 4.4 Add request validation using go-playground/validator
- [ ] 4.5 Add error handling with appropriate HTTP status codes
- [ ] 4.6 Register routes in main API router

## 5. Frontend - Configuration Service
- [ ] 5.1 Create branding API service in `src/services/brandingService.ts`
- [ ] 5.2 Implement fetchBranding() method
- [ ] 5.3 Implement updateBranding() method with admin auth
- [ ] 5.4 Add TypeScript interface for BrandingConfig

## 6. Frontend - Theme System
- [ ] 6.1 Create branding store using Pinia in `src/stores/branding.ts`
- [ ] 6.2 Implement loadBranding() action to fetch and apply configuration
- [ ] 6.3 Create applyBranding() function to update CSS custom properties
- [ ] 6.4 Update document title and favicon dynamically
- [ ] 6.5 Add watchers for branding changes to reapply theme
- [ ] 6.6 Call loadBranding() in App.vue on mount

## 7. Frontend - CSS Custom Properties
- [ ] 7.1 Define CSS custom properties in global styles: --color-primary, --color-secondary
- [ ] 7.2 Update existing components to use CSS custom properties instead of hardcoded colors
- [ ] 7.3 Ensure Tailwind CSS respects custom properties for theme colors
- [ ] 7.4 Test color changes in light and dark modes (if applicable)

## 8. Frontend - Logo and Icon Components
- [ ] 8.1 Update navigation bar logo to use dynamic icon_url from store
- [ ] 8.2 Add fallback to default logo if icon_url is empty
- [ ] 8.3 Handle image loading errors gracefully
- [ ] 8.4 Optimize logo display for different screen sizes

## 9. Admin UI - Branding Management Page
- [ ] 9.1 Create BrandingManagement.vue component in `src/views/admin/`
- [ ] 9.2 Add route for /admin/branding with admin guard
- [ ] 9.3 Create form with fields: app_name, icon_url, favicon_url, primary_color, secondary_color
- [ ] 9.4 Add color pickers for primary and secondary colors
- [ ] 9.5 Implement live preview section showing current branding
- [ ] 9.6 Add form validation (required fields, format validation)
- [ ] 9.7 Implement save functionality with API call
- [ ] 9.8 Show success/error messages after save
- [ ] 9.9 Add reset to defaults button

## 10. RBAC Integration
- [ ] 10.1 Verify admin role has access to PUT /api/v1/branding
- [ ] 10.2 Add RBAC middleware to branding update endpoint
- [ ] 10.3 Add frontend route guard for /admin/branding (admin only)
- [ ] 10.4 Display branding management menu item only for admins

## 11. Documentation
- [ ] 11.1 Update .env.example if any new environment variables added
- [ ] 11.2 Add API endpoint documentation for GET /api/v1/branding
- [ ] 11.3 Add API endpoint documentation for PUT /api/v1/branding
- [ ] 11.4 Document branding configuration fields and validation rules

## 12. Testing
- [ ] 12.1 Write unit tests for branding repository methods
- [ ] 12.2 Write unit tests for branding service validation
- [ ] 12.3 Write integration tests for GET /api/v1/branding endpoint
- [ ] 12.4 Write integration tests for PUT /api/v1/branding endpoint with RBAC
- [ ] 12.5 Test frontend branding store actions
- [ ] 12.6 Test CSS custom properties update correctly
- [ ] 12.7 Manual testing of admin UI with different branding configs
- [ ] 12.8 Test with invalid data (bad colors, long names, invalid URLs)

## 13. Performance and Edge Cases
- [ ] 13.1 Add caching for branding configuration in Redis 
- [ ] 13.2 Test branding load on slow network connections
- [ ] 13.3 Ensure no FOUC (flash of unstyled content) on page load
- [ ] 13.4 Test with very long application names (truncation)
- [ ] 13.5 Test with broken image URLs (fallback behavior)
