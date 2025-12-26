-- Add header_color column to branding_config table
ALTER TABLE branding_config 
ADD COLUMN IF NOT EXISTS header_color VARCHAR(7) NOT NULL DEFAULT '#1e1e1e';

-- Update existing branding_config rows with default header_color if NULL
UPDATE branding_config 
SET header_color = '#1e1e1e' 
WHERE header_color IS NULL;

