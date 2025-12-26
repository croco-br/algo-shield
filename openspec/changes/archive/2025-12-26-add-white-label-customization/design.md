# Design: White Label Customization

## Context
AlgoShield needs white label capabilities to support multiple clients with different branding requirements. This feature enables organizations to customize the application's visual identity (name, logo, colors) without requiring code changes or rebuilds. The system must support dynamic theme application while maintaining the <50ms latency requirement for core transaction processing.

**Constraints:**
- Must not impact transaction processing performance
- Should work with existing RBAC system
- Must support dynamic updates without frontend rebuild
- Should be simple to configure and manage

**Stakeholders:**
- System administrators who deploy AlgoShield for different clients
- End users who interact with the branded interface
- Developers maintaining the codebase

## Goals / Non-Goals

**Goals:**
- Enable runtime customization of application name, logo, and color scheme
- Provide admin UI for managing branding configuration
- Support instant theme application without page refresh
- Maintain single codebase for all deployments
- Ensure branding persists across sessions

**Non-Goals:**
- Custom layouts or component structures (only colors/logos)
- Per-user or per-group branding (single deployment = single brand)
- Multi-tenancy with tenant-specific branding
- Advanced theming (gradients, shadows, typography beyond colors)
- Branding version history or rollback functionality (can be added later)

## Decisions

### Decision 1: Database Storage vs Environment Variables
**Choice:** Store branding configuration in PostgreSQL database

**Rationale:**
- Allows dynamic updates via admin UI without deployment
- Persists configuration across container restarts
- Enables future multi-tenancy if needed
- Consistent with project's data storage pattern

**Alternatives considered:**
- Environment variables: Requires container restart, not user-friendly
- Config files: Requires file system access, harder to manage in containers
- Redis only: Risk of data loss, no persistence guarantee

**Implementation:**
- Single table `branding_config` with one row
- UNIQUE constraint or application logic to enforce single configuration
- Default values in migration script

### Decision 2: CSS Custom Properties for Theming
**Choice:** Use CSS custom properties (CSS variables) for dynamic color application

**Rationale:**
- Native browser support, no additional libraries needed
- Can be updated via JavaScript without CSS recompilation
- Works seamlessly with Tailwind CSS
- No build step required for theme changes

**Alternatives considered:**
- Tailwind config at build time: Requires rebuild for color changes
- CSS-in-JS libraries: Adds complexity and bundle size
- Inline styles: Hard to maintain, poor performance

**Implementation:**
```css
:root {
  --color-primary: #3B82F6;
  --color-secondary: #10B981;
}
```

Update via JavaScript:
```typescript
document.documentElement.style.setProperty('--color-primary', newColor);
```

### Decision 3: Public Branding Endpoint
**Choice:** Make GET /api/v1/branding publicly accessible (no authentication)

**Rationale:**
- Branding info is not sensitive (visible to all users anyway)
- Reduces complexity on login page (needs branding before auth)
- Improves performance (no JWT validation overhead)
- Aligns with typical white label implementations

**Security consideration:**
- No sensitive data in branding config
- Update endpoint (PUT) requires admin authentication

### Decision 4: Pinia Store for Frontend State
**Choice:** Use Pinia store to manage branding state in frontend

**Rationale:**
- Consistent with project's state management pattern
- Provides reactive updates across components
- Easy to integrate with Vue 3 Composition API
- Centralizes branding logic

**Implementation:**
- `src/stores/branding.ts` with actions: loadBranding(), updateBranding()
- Loaded on App.vue mount
- Applied via CSS custom properties and document updates

### Decision 5: Single Configuration Model
**Choice:** Support only one active branding configuration per deployment

**Rationale:**
- Simplifies implementation (no tenant routing logic)
- Matches stated requirement (per-deployment, not per-user)
- Reduces complexity and potential bugs
- Can be extended to multi-tenant later if needed

**Migration path for multi-tenancy:**
- Add tenant_id column to branding_config table
- Update API to filter by tenant context
- Add frontend logic to detect tenant from domain/subdomain

### Decision 6: No Caching in First Iteration
**Choice:** Query database directly for branding config, skip Redis caching initially

**Rationale:**
- Branding data is small and rarely changes
- Single query per page load has negligible impact
- Simplifies implementation
- Can add Redis caching later if profiling shows need

**Future optimization:**
- Add Redis cache with TTL (e.g., 5 minutes)
- Invalidate cache on PUT /api/v1/branding

### Decision 7: Validation Strategy
**Choice:** Validate on both frontend (UX) and backend (security)

**Rationale:**
- Frontend validation: Immediate user feedback, better UX
- Backend validation: Security boundary, prevents invalid data
- Use go-playground/validator for backend validation

**Validation rules:**
- app_name: required, max 100 characters
- primary_color/secondary_color: hex format (#RGB or #RRGGBB)
- icon_url/favicon_url: valid URL or relative path format
- All fields have sensible defaults if not provided

## Architecture

### Data Flow

**Initial Load:**
```
Browser → GET /api/v1/branding → API Handler → Service → Repository → PostgreSQL
      ← BrandingConfig JSON ← ← ← ← ←
Browser applies branding (CSS vars, title, favicon)
```

**Admin Update:**
```
Admin UI → PUT /api/v1/branding + JWT → JWT Middleware → RBAC Check → Handler
       → Service (validate) → Repository → PostgreSQL
      ← Success Response ← ← ← ←
Admin UI triggers store reload → GET /api/v1/branding → Apply new branding
```

### Database Schema

```sql
CREATE TABLE branding_config (
    id SERIAL PRIMARY KEY,
    app_name VARCHAR(100) NOT NULL DEFAULT 'AlgoShield',
    icon_url VARCHAR(500),
    favicon_url VARCHAR(500),
    primary_color VARCHAR(7) NOT NULL DEFAULT '#3B82F6',
    secondary_color VARCHAR(7) NOT NULL DEFAULT '#10B981',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Ensure only one configuration exists
CREATE UNIQUE INDEX idx_single_config ON branding_config ((id IS NOT NULL));
```

### API Contract

**GET /api/v1/branding**
- Authentication: None (public endpoint)
- Response: 200 OK
```json
{
  "app_name": "AlgoShield",
  "icon_url": "/assets/logo.svg",
  "favicon_url": "/favicon.ico",
  "primary_color": "#3B82F6",
  "secondary_color": "#10B981"
}
```

**PUT /api/v1/branding**
- Authentication: JWT token required
- Authorization: Admin role only
- Request Body:
```json
{
  "app_name": "Custom Brand",
  "icon_url": "https://example.com/logo.png",
  "favicon_url": "https://example.com/favicon.ico",
  "primary_color": "#FF5733",
  "secondary_color": "#C70039"
}
```
- Response: 200 OK (updated config) or 400/403 errors

### Frontend Structure

```
src/
├── stores/
│   └── branding.ts           # Pinia store for branding state
├── services/
│   └── brandingService.ts    # API client for branding endpoints
├── views/admin/
│   └── BrandingManagement.vue # Admin UI for branding
└── App.vue                    # Load branding on mount
```

## Risks / Trade-offs

### Risk 1: Flash of Unstyled Content (FOUC)
**Risk:** Brief moment where default colors show before branding loads

**Mitigation:**
- Load branding as early as possible in App.vue
- Show loading spinner until branding is applied
- Consider SSR or build-time injection for critical deployments (future)

**Trade-off:** Accepting small FOUC for simplicity of runtime configuration

### Risk 2: Invalid Image URLs
**Risk:** Broken images if icon_url or favicon_url are incorrect

**Mitigation:**
- Validate URLs in backend before saving
- Use img onerror handler to fallback to default logo
- Test image accessibility in admin UI preview

**Trade-off:** Cannot guarantee external URLs remain valid over time

### Risk 3: Color Accessibility
**Risk:** Custom colors may have poor contrast or accessibility issues

**Mitigation:**
- Document recommended color contrast ratios in admin UI
- Provide color picker with preview
- Future: Add automated contrast checking (WCAG AA standard)

**Trade-off:** Trust admin users to choose accessible colors initially

### Risk 4: Browser Compatibility
**Risk:** CSS custom properties not supported in very old browsers

**Mitigation:**
- CSS custom properties supported in all modern browsers (Chrome 49+, Firefox 31+, Safari 9.1+)
- Provide fallback default colors if custom properties fail
- Document minimum browser requirements

**Trade-off:** Not supporting IE11 (acceptable per project context)

## Migration Plan

### Deployment Steps
1. Run database migration to create `branding_config` table with defaults
2. Deploy backend with new branding endpoints
3. Deploy frontend with branding store and theme application
4. Verify default branding appears correctly
5. Test admin UI for branding updates
6. Document configuration process for deployment teams

### Rollback Plan
1. Revert frontend deployment to remove branding loading logic
2. Revert backend deployment to remove branding endpoints
3. Run down migration to drop `branding_config` table (optional)
4. Application returns to hardcoded branding

### Data Migration
- No existing data to migrate
- Default values inserted via migration script
- Future: Export/import branding config for backup/restore

## Performance Considerations

### Impact on Transaction Processing
- Branding endpoints are separate from transaction processing flow
- No impact on <50ms transaction latency requirement
- Database connection pooling ensures efficient branding queries

### Frontend Performance
- Branding loaded once on application mount, cached in Pinia store
- CSS custom property updates are near-instant (no reflow)
- Logo images lazy-loaded as needed
- Estimated overhead: <100ms on initial page load

### Optimization Opportunities (Future)
1. Add Redis caching for branding config (if profiling shows benefit)
2. Preload branding in HTML head via SSR (eliminate FOUC)
3. Bundle default logo to avoid network request
4. Add CDN caching headers for branding endpoint

## Open Questions

1. Should we support custom fonts in addition to colors?
   - **Decision:** No, out of scope for initial version. Can be added later if needed.

2. Should branding history be tracked for audit purposes?
   - **Decision:** No, not required initially. Future enhancement if compliance needs arise.

3. Should we support dark mode variations of branding colors?
   - **Decision:** Yes, but in a future iteration. Initial version uses same colors for light/dark modes.

4. Should logo support multiple sizes/formats (mobile, desktop, high-DPI)?
   - **Decision:** Single logo URL initially. Responsive sizing handled by CSS. Future: multiple size support.

## Testing Strategy

### Unit Tests
- Branding service validation logic (color format, URL validation, length checks)
- Repository CRUD operations with mock database
- Pinia store actions with mocked API responses

### Integration Tests
- GET /api/v1/branding returns correct data
- PUT /api/v1/branding requires admin authentication
- RBAC enforcement on branding update endpoint
- Database constraints (unique config, default values)

### Manual Testing Checklist
- [ ] Default branding displays correctly on fresh install
- [ ] Admin can update all branding fields successfully
- [ ] Non-admin users cannot access PUT endpoint
- [ ] Invalid color formats are rejected with clear error
- [ ] Logo changes appear immediately after save
- [ ] Page title updates correctly
- [ ] Favicon updates without page refresh
- [ ] Broken image URLs fallback to default logo
- [ ] Long application names don't break layout
- [ ] Branding persists after browser refresh
- [ ] Multiple tabs show consistent branding after update

## Future Enhancements

1. **Dark Mode Branding Variants**
   - Add `primary_color_dark`, `secondary_color_dark` fields
   - Automatically adjust for dark mode based on user preference

2. **Advanced Color Customization**
   - Success, warning, error, info color overrides
   - Hover, active, disabled state color variations
   - Gradient support for backgrounds

3. **Multi-Tenancy Support**
   - Add tenant_id to branding_config
   - Domain-based tenant detection
   - Tenant-specific branding application

4. **Branding Templates**
   - Predefined color schemes (e.g., "Corporate Blue", "Financial Green")
   - One-click template application

5. **Import/Export Configuration**
   - JSON export of current branding
   - Import from file for easy migration between environments

6. **Custom Fonts**
   - Font family customization
   - Google Fonts integration
   - Font weight and size adjustments

7. **Branding Analytics**
   - Track which branding configurations are most used
   - A/B testing support for branding variations

8. **Redis Caching**
   - Cache branding config in Redis with TTL
   - Invalidation on update for near-zero latency

## References

- Project conventions: `/openspec/project.md`
- RBAC implementation: `internal/api/auth/middleware.go`
- Pinia stores pattern: `src/stores/`
- Vue 3 Composition API: https://vuejs.org/guide/extras/composition-api-faq.html
- CSS Custom Properties: https://developer.mozilla.org/en-US/docs/Web/CSS/Using_CSS_custom_properties
