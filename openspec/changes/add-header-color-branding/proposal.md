# Change: Add Header Color to Branding Configuration

## Why
Currently, the header uses a fixed dark gradient background that cannot be customized. To provide complete white label customization, administrators should be able to customize the header background color to match their brand identity.

## What Changes
- Add `header_color` field to branding configuration (hex color format)
- Update database schema to include `header_color` column
- Extend API to accept and return `header_color` in branding endpoints
- Apply header color dynamically in frontend via CSS custom property
- Update admin UI to include header color picker with preview
- Provide default header color value

## Impact
- Affected specs: `branding-customization` (MODIFIED requirement)
- Affected code:
  - Backend: Database migration, models, repository, service, handler
  - Frontend: Branding store, Header component, BrandingView admin page
  - Database: New column in `branding_config` table

