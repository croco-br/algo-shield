-- Create event_schemas table
CREATE TABLE IF NOT EXISTS event_schemas (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    sample_json JSONB NOT NULL,
    extracted_fields JSONB NOT NULL DEFAULT '[]',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create index for name lookups
CREATE INDEX IF NOT EXISTS idx_event_schemas_name ON event_schemas(name);

-- Add schema_id to rules table (nullable initially for migration)
ALTER TABLE rules ADD COLUMN IF NOT EXISTS schema_id UUID REFERENCES event_schemas(id);

-- Create index for schema lookups on rules
CREATE INDEX IF NOT EXISTS idx_rules_schema_id ON rules(schema_id);

