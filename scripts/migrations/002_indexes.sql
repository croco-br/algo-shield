-- ============================================================================
-- Migration 002: Database Indexes
-- Creates all indexes for performance optimization
-- ============================================================================

-- ----------------------------------------------------------------------------
-- Event Schemas Indexes
-- ----------------------------------------------------------------------------
CREATE INDEX IF NOT EXISTS idx_event_schemas_name ON event_schemas(name);

-- ----------------------------------------------------------------------------
-- Transactions Indexes
-- ----------------------------------------------------------------------------
CREATE INDEX IF NOT EXISTS idx_transactions_schema_id ON transactions(schema_id);
CREATE INDEX IF NOT EXISTS idx_transactions_status ON transactions(status);
CREATE INDEX IF NOT EXISTS idx_transactions_created_at ON transactions(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_transactions_status_created_at ON transactions(status, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_transactions_schema_status ON transactions(schema_id, status);
CREATE INDEX IF NOT EXISTS idx_transactions_created_at_status ON transactions(created_at, status);

-- ----------------------------------------------------------------------------
-- Synthetic Transactions Indexes
-- ----------------------------------------------------------------------------
CREATE INDEX IF NOT EXISTS idx_synthetic_transactions_schema_id ON synthetic_transactions(schema_id);
CREATE INDEX IF NOT EXISTS idx_synthetic_transactions_status ON synthetic_transactions(status);
CREATE INDEX IF NOT EXISTS idx_synthetic_transactions_created_at ON synthetic_transactions(created_at DESC);

-- ----------------------------------------------------------------------------
-- Rules Indexes
-- ----------------------------------------------------------------------------
CREATE INDEX IF NOT EXISTS idx_rules_enabled ON rules(enabled);
CREATE INDEX IF NOT EXISTS idx_rules_priority ON rules(priority ASC);
CREATE INDEX IF NOT EXISTS idx_rules_schema_id ON rules(schema_id);

-- ----------------------------------------------------------------------------
-- Users Indexes
-- ----------------------------------------------------------------------------
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_google_id ON users(google_id);
CREATE INDEX IF NOT EXISTS idx_users_active ON users(active);

-- ----------------------------------------------------------------------------
-- User-Role Junction Indexes
-- ----------------------------------------------------------------------------
CREATE INDEX IF NOT EXISTS idx_user_roles_user_id ON user_roles(user_id);
CREATE INDEX IF NOT EXISTS idx_user_roles_role_id ON user_roles(role_id);

-- ----------------------------------------------------------------------------
-- User-Group Junction Indexes
-- ----------------------------------------------------------------------------
CREATE INDEX IF NOT EXISTS idx_user_groups_user_id ON user_groups(user_id);
CREATE INDEX IF NOT EXISTS idx_user_groups_group_id ON user_groups(group_id);

-- ----------------------------------------------------------------------------
-- Group-Role Junction Indexes
-- ----------------------------------------------------------------------------
CREATE INDEX IF NOT EXISTS idx_group_roles_group_id ON group_roles(group_id);
CREATE INDEX IF NOT EXISTS idx_group_roles_role_id ON group_roles(role_id);

-- ----------------------------------------------------------------------------
-- Sessions Indexes
-- ----------------------------------------------------------------------------
CREATE INDEX IF NOT EXISTS idx_sessions_user_id ON sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_sessions_token_hash ON sessions(token_hash);
CREATE INDEX IF NOT EXISTS idx_sessions_expires_at ON sessions(expires_at);

-- ----------------------------------------------------------------------------
-- Branding Config Indexes
-- ----------------------------------------------------------------------------
CREATE UNIQUE INDEX IF NOT EXISTS idx_single_branding_config ON branding_config ((id IS NOT NULL));
