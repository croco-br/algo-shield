-- Add schema_id to transactions table for linking transactions to event schemas
-- Nullable initially to support existing transactions without schema association

ALTER TABLE transactions ADD COLUMN IF NOT EXISTS schema_id UUID REFERENCES event_schemas(id);

-- Index for filtering transactions by schema
CREATE INDEX IF NOT EXISTS idx_transactions_schema_id ON transactions(schema_id);

-- Composite index for common filter combinations (status + created_at for dashboard)
CREATE INDEX IF NOT EXISTS idx_transactions_status_created_at ON transactions(status, created_at DESC);

-- Composite index for schema + status filtering
CREATE INDEX IF NOT EXISTS idx_transactions_schema_status ON transactions(schema_id, status);

-- Index for amount range queries
CREATE INDEX IF NOT EXISTS idx_transactions_amount ON transactions(amount);

-- Composite index for temporal dashboard aggregations
CREATE INDEX IF NOT EXISTS idx_transactions_created_at_status ON transactions(created_at, status);
