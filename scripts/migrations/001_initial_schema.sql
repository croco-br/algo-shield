-- Create transactions table
CREATE TABLE IF NOT EXISTS transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    external_id VARCHAR(255) UNIQUE NOT NULL,
    amount DECIMAL(20, 2) NOT NULL,
    currency VARCHAR(3) NOT NULL,
    from_account VARCHAR(255) NOT NULL,
    to_account VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL,
    status VARCHAR(50) NOT NULL,
    risk_score DECIMAL(5, 2) DEFAULT 0,
    risk_level VARCHAR(20) DEFAULT 'low',
    processing_time BIGINT DEFAULT 0,
    matched_rules JSONB DEFAULT '[]',
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    processed_at TIMESTAMP WITH TIME ZONE
);

-- Create rules table
CREATE TABLE IF NOT EXISTS rules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    type VARCHAR(50) NOT NULL,
    action VARCHAR(50) NOT NULL,
    priority INTEGER DEFAULT 0,
    enabled BOOLEAN DEFAULT true,
    conditions JSONB DEFAULT '{}',
    score DECIMAL(5, 2) DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_transactions_external_id ON transactions(external_id);
CREATE INDEX IF NOT EXISTS idx_transactions_from_account ON transactions(from_account);
CREATE INDEX IF NOT EXISTS idx_transactions_to_account ON transactions(to_account);
CREATE INDEX IF NOT EXISTS idx_transactions_status ON transactions(status);
CREATE INDEX IF NOT EXISTS idx_transactions_created_at ON transactions(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_rules_enabled ON rules(enabled);
CREATE INDEX IF NOT EXISTS idx_rules_priority ON rules(priority ASC);

-- Insert sample rules
INSERT INTO rules (id, name, description, type, action, priority, enabled, conditions, score) VALUES
    (gen_random_uuid(), 'High Amount Transaction', 'Flag transactions above $10,000', 'amount', 'score', 10, true, '{"amount_threshold": 10000}', 30),
    (gen_random_uuid(), 'Very High Amount Transaction', 'Block transactions above $50,000', 'amount', 'block', 5, true, '{"amount_threshold": 50000}', 100),
    (gen_random_uuid(), 'Transaction Velocity Check', 'Flag accounts with more than 10 transactions per hour', 'velocity', 'score', 20, true, '{"transaction_count": 10, "time_window_seconds": 3600}', 40),
    (gen_random_uuid(), 'Blocklist Check', 'Block transactions from blocklisted accounts', 'blocklist', 'block', 1, true, '{"blocklisted_accounts": []}', 100)
ON CONFLICT DO NOTHING;

