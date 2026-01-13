-- Migration: Add composite index for velocity queries performance
-- This index optimizes velocityCount and velocitySum queries that filter by origin and created_at

-- Composite index for velocity queries (origin + created_at)
-- This significantly speeds up queries like:
--   SELECT COUNT(*) FROM transactions WHERE origin = $1 AND created_at > NOW() - INTERVAL...
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_transactions_origin_created_at 
ON transactions(origin, created_at DESC);

-- Same composite index for synthetic_transactions
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_synthetic_transactions_origin_created_at 
ON synthetic_transactions(origin, created_at DESC);
