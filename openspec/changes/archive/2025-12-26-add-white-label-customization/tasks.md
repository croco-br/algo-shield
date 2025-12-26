# Implementation Tasks

## 1. Database Schema
- [x] 1.1 Create migration script for `branding_config` table with columns: id, app_name, icon_url, primary_color, secondary_color, favicon_url, created_at, updated_at
- [x] 1.2 Add unique constraint to ensure single branding configuration row
- [x] 1.3 Insert default branding values in migration
- [ ] 1.4 Test migration (up and down)

## 2. Backend API - Models and Repository
- [x] 2.1 Create `BrandingConfig` model in `pkg/models/branding.go` (model exists, location differs from original plan)
- [x] 2.2 Create repository interface in `internal/api/branding/repository.go`
- [x] 2.3 Implement PostgreSQL repository with methods: Get(), Update() (Create() not needed - uses fixed id=1)
- [x] 2.4 Add connection pooling for database access (via pgxpool.Pool)
- [x] 2.5 Add validation struct tags for branding fields (in models.go and UpdateBrandingRequest)

## 3. Backend API - Service Layer
- [x] 3.1 Create branding service in `internal/api/branding/service.go`
- [x] 3.2 Implement GetBranding() method with default fallback
- [x] 3.3 Implement UpdateBranding() method with validation
- [x] 3.4 Add color format validation (hex: #RGB or #RRGGBB)
- [x] 3.5 Add URL validation for icon_url and favicon_url (via validator struct tags)
- [x] 3.6 Add app_name length validation (max 100 chars)

## 4. Backend API - HTTP Handlers
- [x] 4.1 Create handler in `internal/api/branding/handler.go`
- [x] 4.2 Implement GET /api/v1/branding endpoint (public, no auth required)
- [x] 4.3 Implement PUT /api/v1/branding endpoint (admin only, JWT + RBAC)
- [x] 4.4 Add request validation using go-playground/validator
- [x] 4.5 Add error handling with appropriate HTTP status codes
- [x] 4.6 Register routes in main API router

## 5. Frontend - Configuration Service
- [x] 5.1 Create branding API service (integrated in `src/stores/branding.ts` instead of separate service file)
- [x] 5.2 Implement fetchBranding() method (as loadBranding() in store)
- [x] 5.3 Implement updateBranding() method with admin auth
- [x] 5.4 Add TypeScript interface for BrandingConfig

## 6. Frontend - Theme System
- [x] 6.1 Create branding store using Pinia in `src/stores/branding.ts`
- [x] 6.2 Implement loadBranding() action to fetch and apply configuration
- [x] 6.3 Create applyBranding() function to update CSS custom properties
- [x] 6.4 Update document title and favicon dynamically
- [x] 6.5 Add watchers for branding changes to reapply theme (auto-applies on update)
- [x] 6.6 Call loadBranding() in App.vue on mount (auto-initializes in store constructor)

## 7. Frontend - CSS Custom Properties
- [x] 7.1 Define CSS custom properties in global styles: --color-primary, --color-secondary
- [x] 7.2 Update existing components to use CSS custom properties instead of hardcoded colors
- [x] 7.3 Ensure Tailwind CSS respects custom properties for theme colors
- [ ] 7.4 Test color changes in light and dark modes (if applicable)

## 8. Frontend - Logo and Icon Components
- [x] 8.1 Update navigation bar logo to use dynamic icon_url from store (Header.vue)
- [x] 8.2 Add fallback to default logo if icon_url is empty
- [x] 8.3 Handle image loading errors gracefully (@error handler)
- [x] 8.4 Optimize logo display for different screen sizes (responsive sizing with Tailwind classes)

## 9. Admin UI - Branding Management Page
- [x] 9.1 Create BrandingView.vue component in `src/views/` (location differs from original plan)
- [x] 9.2 Add route for /branding with admin guard (route path differs from original plan)
- [x] 9.3 Create form with fields: app_name, icon_url, favicon_url, primary_color, secondary_color
- [x] 9.4 Add color pickers for primary and secondary colors
- [x] 9.5 Implement live preview section showing current branding
- [x] 9.6 Add form validation (required fields, format validation)
- [x] 9.7 Implement save functionality with API call
- [x] 9.8 Show success/error messages after save
- [x] 9.9 Add reset to defaults button

## 10. RBAC Integration
- [x] 10.1 Verify admin role has access to PUT /api/v1/branding
- [x] 10.2 Add RBAC middleware to branding update endpoint (RequireRole("admin"))
- [x] 10.3 Add frontend route guard for /branding (admin only via requiresAdmin meta)
- [x] 10.4 Display branding management menu item only for admins (Sidebar.vue filters by adminOnly)

## 11. Documentation
- [x] 11.1 Update .env.example if any new environment variables added (no new env vars needed - branding stored in DB)
- [x] 11.2 Add API endpoint documentation for GET /api/v1/branding (added to README.md)
- [x] 11.3 Add API endpoint documentation for PUT /api/v1/branding (added to README.md)
- [x] 11.4 Document branding configuration fields and validation rules (added to README.md)

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
- [x] 13.1 Add caching for branding configuration in Redis (10min TTL, cache invalidation on update)
- [ ] 13.2 Test branding load on slow network connections
- [x] 13.3 Ensure no FOUC (flash of unstyled content) on page load (inline script + default branding on store init)
- [ ] 13.4 Test with very long application names (truncation)
- [ ] 13.5 Test with broken image URLs (fallback behavior)
