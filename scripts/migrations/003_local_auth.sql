-- Add password_hash column to users table
ALTER TABLE users ADD COLUMN IF NOT EXISTS password_hash VARCHAR(255);

-- Update auth_type default to 'local'
ALTER TABLE users ALTER COLUMN auth_type SET DEFAULT 'local';

-- Remove google_id column (optional, can be kept for migration safety)
-- ALTER TABLE users DROP COLUMN IF EXISTS google_id;

-- Update existing users to have local auth_type if they don't have a password
UPDATE users SET auth_type = 'local' WHERE password_hash IS NOT NULL;
