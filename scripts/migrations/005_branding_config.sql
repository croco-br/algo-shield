-- Create branding_config table for white label customization
CREATE TABLE IF NOT EXISTS branding_config (
    id SERIAL PRIMARY KEY,
    app_name VARCHAR(100) NOT NULL DEFAULT 'AlgoShield',
    icon_url VARCHAR(500) DEFAULT '/assets/logo.svg',
    favicon_url VARCHAR(500) DEFAULT '/favicon.ico',
    primary_color VARCHAR(7) NOT NULL DEFAULT '#3B82F6',
    secondary_color VARCHAR(7) NOT NULL DEFAULT '#10B981',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Ensure only one branding configuration exists
CREATE UNIQUE INDEX IF NOT EXISTS idx_single_branding_config ON branding_config ((id IS NOT NULL));

-- Insert default branding configuration
INSERT INTO branding_config (id, app_name, icon_url, favicon_url, primary_color, secondary_color)
VALUES (1, 'AlgoShield', '/assets/logo.svg', '/favicon.ico', '#3B82F6', '#10B981')
ON CONFLICT DO NOTHING;
