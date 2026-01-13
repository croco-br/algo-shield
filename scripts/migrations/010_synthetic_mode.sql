-- Migration: Add synthetic mode support
-- This migration adds:
-- 1. A synthetic_transactions table (mirror of transactions)
-- 2. A system_config table for global settings including synthetic mode

-- Create system_config table for global settings
CREATE TABLE IF NOT EXISTS system_config (
    key VARCHAR(255) PRIMARY KEY,
    value JSONB NOT NULL DEFAULT '{}',
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Insert default synthetic_mode setting
INSERT INTO system_config (key, value, updated_at)
VALUES ('synthetic_mode', '{"enabled": false}'::jsonb, NOW())
ON CONFLICT (key) DO NOTHING;

-- Create synthetic_transactions table (mirror of transactions)
CREATE TABLE IF NOT EXISTS synthetic_transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    external_id VARCHAR(255) NOT NULL,
    schema_id UUID REFERENCES event_schemas(id),
    amount DECIMAL(20, 2) NOT NULL,
    currency VARCHAR(3) NOT NULL,
    origin VARCHAR(255) NOT NULL,
    destination VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'approved',
    processing_time BIGINT DEFAULT 0,
    matched_rules JSONB DEFAULT '[]',
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    processed_at TIMESTAMP WITH TIME ZONE
);

-- Create indexes for synthetic_transactions (mirror of transactions indexes)
CREATE INDEX IF NOT EXISTS idx_synthetic_transactions_external_id ON synthetic_transactions(external_id);
CREATE INDEX IF NOT EXISTS idx_synthetic_transactions_origin ON synthetic_transactions(origin);
CREATE INDEX IF NOT EXISTS idx_synthetic_transactions_destination ON synthetic_transactions(destination);
CREATE INDEX IF NOT EXISTS idx_synthetic_transactions_status ON synthetic_transactions(status);
CREATE INDEX IF NOT EXISTS idx_synthetic_transactions_created_at ON synthetic_transactions(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_synthetic_transactions_schema_id ON synthetic_transactions(schema_id);
