-- ============================================================================
-- Migration 003: Test Data
-- Inserts default roles, admin user, branding config, system config, and test data
-- ============================================================================

-- ----------------------------------------------------------------------------
-- Default Roles
-- ----------------------------------------------------------------------------
INSERT INTO roles (id, name, description) VALUES
    (gen_random_uuid(), 'admin', 'Administrator with full system access'),
    (gen_random_uuid(), 'rule_editor', 'Can create, edit and delete rules'),
    (gen_random_uuid(), 'viewer', 'Read-only access to view rules and transactions')
ON CONFLICT (name) DO NOTHING;

-- ----------------------------------------------------------------------------
-- Default Admin User
-- Email: admin@admin.com
-- Password: admin@123
-- ----------------------------------------------------------------------------
DO $$
DECLARE
    admin_user_id UUID;
    admin_role_id UUID;
    admin_email VARCHAR(255) := 'admin@admin.com';
BEGIN
    SELECT id INTO admin_role_id FROM roles WHERE name = 'admin';
    IF admin_role_id IS NULL THEN
        RAISE EXCEPTION 'Admin role not found. Please ensure roles have been inserted.';
    END IF;
    SELECT id INTO admin_user_id FROM users WHERE email = admin_email;
    IF admin_user_id IS NULL THEN
        admin_user_id := gen_random_uuid();
        INSERT INTO users (id, email, name, password_hash, auth_type, active, created_at, updated_at)
        VALUES (admin_user_id, admin_email, 'Administrator', '$2a$10$IIbu/Hx8lQJanbd0Rr3OeunWWVDF.m6PdRErfcFpZbaJkSsNoJX0.', 'local', true, NOW(), NOW());
        RAISE NOTICE 'Admin user created successfully with email: % and password: admin@123', admin_email;
    ELSE
        RAISE NOTICE 'Admin user already exists, ensuring admin role is assigned';
    END IF;
    INSERT INTO user_roles (user_id, role_id, assigned_at)
    VALUES (admin_user_id, admin_role_id, NOW())
    ON CONFLICT (user_id, role_id) DO NOTHING;
    RAISE NOTICE 'Admin role assigned to user with email: %', admin_email;
END $$;

-- ----------------------------------------------------------------------------
-- Branding Configuration
-- ----------------------------------------------------------------------------
INSERT INTO branding_config (id, app_name, icon_url, favicon_url, primary_color, secondary_color, header_color)
VALUES (1, 'AlgoShield', '/assets/logo.svg', '/favicon.ico', '#3B82F6', '#10B981', '#1e1e1e')
ON CONFLICT DO NOTHING;

-- ----------------------------------------------------------------------------
-- System Configuration
-- ----------------------------------------------------------------------------
INSERT INTO system_config (key, value, updated_at)
VALUES ('synthetic_mode', '{"enabled": false}'::jsonb, NOW())
ON CONFLICT (key) DO NOTHING;

-- ----------------------------------------------------------------------------
-- Test Event Schema: Payment Transaction Example
-- ----------------------------------------------------------------------------
INSERT INTO event_schemas (id, name, description, sample_json, extracted_fields, created_at, updated_at)
VALUES (
  '550e8400-e29b-41d4-a716-446655440000'::uuid,
  'Payment Transaction Example',
  'Comprehensive payment transaction example schema demonstrating all field types: strings, numbers, booleans, arrays, nested objects, and geographic data',
  '{
    "external_id": "txn_123456789",
    "amount": 5000.50,
    "currency": "USD",
    "origin": "ACC001",
    "destination": "ACC002",
    "type": "transfer",
    "timestamp": 1704067200000,
    "created_at": "2024-01-01T12:00:00Z",
    "is_verified": true,
    "is_suspicious": false,
    "tags": ["high-value", "international", "urgent"],
    "categories": ["payment", "transfer"],
    "metadata": {
      "ip_address": "192.168.1.100",
      "device_id": "device_abc123",
      "user_agent": "Mozilla/5.0",
      "is_suspicious": false,
      "risk_score": 25.5
    },
    "location": {
      "lat": 37.7749,
      "lon": -122.4194,
      "city": "San Francisco",
      "country": "US",
      "timezone": "America/Los_Angeles"
    },
    "user": {
      "id": "user_123",
      "email": "user@example.com",
      "country": "US",
      "account_type": "premium",
      "verification_status": "verified"
    },
    "payment_method": {
      "type": "credit_card",
      "last_four": "1234",
      "issuer": "Visa"
    }
  }'::jsonb,
  '[
    {"path": "external_id", "type": "string", "nullable": false, "sample_value": "txn_123456789"},
    {"path": "amount", "type": "number", "nullable": false, "sample_value": 5000.50},
    {"path": "currency", "type": "string", "nullable": false, "sample_value": "USD"},
    {"path": "origin", "type": "string", "nullable": false, "sample_value": "ACC001"},
    {"path": "destination", "type": "string", "nullable": false, "sample_value": "ACC002"},
    {"path": "type", "type": "string", "nullable": false, "sample_value": "transfer"},
    {"path": "timestamp", "type": "number", "nullable": false, "sample_value": 1704067200000},
    {"path": "created_at", "type": "datetime", "nullable": false, "sample_value": "2024-01-01T12:00:00Z"},
    {"path": "is_verified", "type": "boolean", "nullable": false, "sample_value": true},
    {"path": "is_suspicious", "type": "boolean", "nullable": false, "sample_value": false},
    {"path": "tags", "type": "array", "nullable": false, "sample_value": ["high-value", "international", "urgent"]},
    {"path": "categories", "type": "array", "nullable": false, "sample_value": ["payment", "transfer"]},
    {"path": "metadata.ip_address", "type": "string", "nullable": false, "sample_value": "192.168.1.100"},
    {"path": "metadata.device_id", "type": "string", "nullable": false, "sample_value": "device_abc123"},
    {"path": "metadata.user_agent", "type": "string", "nullable": false, "sample_value": "Mozilla/5.0"},
    {"path": "metadata.is_suspicious", "type": "boolean", "nullable": false, "sample_value": false},
    {"path": "metadata.risk_score", "type": "number", "nullable": false, "sample_value": 25.5},
    {"path": "location.lat", "type": "number", "nullable": false, "sample_value": 37.7749},
    {"path": "location.lon", "type": "number", "nullable": false, "sample_value": -122.4194},
    {"path": "location.city", "type": "string", "nullable": false, "sample_value": "San Francisco"},
    {"path": "location.country", "type": "string", "nullable": false, "sample_value": "US"},
    {"path": "location.timezone", "type": "string", "nullable": false, "sample_value": "America/Los_Angeles"},
    {"path": "user.id", "type": "string", "nullable": false, "sample_value": "user_123"},
    {"path": "user.email", "type": "string", "nullable": false, "sample_value": "user@example.com"},
    {"path": "user.country", "type": "string", "nullable": false, "sample_value": "US"},
    {"path": "user.account_type", "type": "string", "nullable": false, "sample_value": "premium"},
    {"path": "user.verification_status", "type": "string", "nullable": false, "sample_value": "verified"},
    {"path": "payment_method.type", "type": "string", "nullable": false, "sample_value": "credit_card"},
    {"path": "payment_method.last_four", "type": "string", "nullable": false, "sample_value": "1234"},
    {"path": "payment_method.issuer", "type": "string", "nullable": false, "sample_value": "Visa"}
  ]'::jsonb,
  NOW(),
  NOW()
)
ON CONFLICT (name) DO NOTHING;

-- ----------------------------------------------------------------------------
-- Test Rules
-- ----------------------------------------------------------------------------

-- Rule 1: High Value Transaction (Block action, high priority)
INSERT INTO rules (id, name, description, action, priority, enabled, conditions, schema_id, created_at, updated_at)
VALUES (
  '660e8400-e29b-41d4-a716-446655440001'::uuid,
  'Block High Value Transactions',
  'Example: Block transactions over $10,000 immediately',
  'block',
  10,
  true,
  '{"custom_expression": "amount > 10000"}'::jsonb,
  '550e8400-e29b-41d4-a716-446655440000'::uuid,
  NOW(),
  NOW()
)
ON CONFLICT (id) DO NOTHING;

-- Rule 2: Suspicious Flag Check (Review action, medium priority)
INSERT INTO rules (id, name, description, action, priority, enabled, conditions, schema_id, created_at, updated_at)
VALUES (
  '660e8400-e29b-41d4-a716-446655440002'::uuid,
  'Review Suspicious Transactions',
  'Example: Flag transactions with suspicious metadata for manual review',
  'review',
  50,
  true,
  '{"custom_expression": "metadata.is_suspicious == true"}'::jsonb,
  '550e8400-e29b-41d4-a716-446655440000'::uuid,
  NOW(),
  NOW()
)
ON CONFLICT (id) DO NOTHING;

-- Rule 3: Blocklist Rule (Block action, very high priority)
INSERT INTO rules (id, name, description, action, priority, enabled, conditions, schema_id, created_at, updated_at)
VALUES (
  '660e8400-e29b-41d4-a716-446655440003'::uuid,
  'Block Blocklisted Accounts',
  'Example: Block transactions from known fraudulent accounts',
  'block',
  5,
  true,
  '{"custom_expression": "origin in [\"BLOCKED001\", \"BLOCKED002\", \"BLOCKED003\"]"}'::jsonb,
  '550e8400-e29b-41d4-a716-446655440000'::uuid,
  NOW(),
  NOW()
)
ON CONFLICT (id) DO NOTHING;

-- Rule 4: Geographic Restriction (Review action, medium priority)
INSERT INTO rules (id, name, description, action, priority, enabled, conditions, schema_id, created_at, updated_at)
VALUES (
  '660e8400-e29b-41d4-a716-446655440004'::uuid,
  'Review High-Risk Geographic Areas',
  'Example: Review transactions from high-risk geographic regions using polygon coordinates',
  'review',
  60,
  true,
  '{"custom_expression": "pointInPolygon(location.lat, location.lon, [[55.0, 30.0], [60.0, 30.0], [60.0, 45.0], [55.0, 45.0]]) or pointInPolygon(location.lat, location.lon, [[35.0, 110.0], [45.0, 110.0], [45.0, 125.0], [35.0, 125.0]]) or pointInPolygon(location.lat, location.lon, [[38.0, 124.0], [42.0, 124.0], [42.0, 130.0], [38.0, 130.0]]) or pointInPolygon(location.lat, location.lon, [[30.0, 44.0], [40.0, 44.0], [40.0, 63.0], [30.0, 63.0]])"}'::jsonb,
  '550e8400-e29b-41d4-a716-446655440000'::uuid,
  NOW(),
  NOW()
)
ON CONFLICT (id) DO NOTHING;

-- Rule 5: Velocity Check (Review action, medium priority)
INSERT INTO rules (id, name, description, action, priority, enabled, conditions, schema_id, created_at, updated_at)
VALUES (
  '660e8400-e29b-41d4-a716-446655440005'::uuid,
  'Review High Velocity Transactions',
  'Example: Review accounts with more than 10 transactions in the last hour',
  'review',
  55,
  true,
  '{"custom_expression": "velocityCount(\"origin\", 3600) > 10"}'::jsonb,
  '550e8400-e29b-41d4-a716-446655440000'::uuid,
  NOW(),
  NOW()
)
ON CONFLICT (id) DO NOTHING;

-- Rule 6: Amount Velocity Check (Review action, medium priority)
INSERT INTO rules (id, name, description, action, priority, enabled, conditions, schema_id, created_at, updated_at)
VALUES (
  '660e8400-e29b-41d4-a716-446655440006'::uuid,
  'Review High Amount Velocity',
  'Example: Review accounts with transaction sums exceeding $10,000 in the last hour',
  'review',
  56,
  true,
  '{"custom_expression": "velocitySum(\"origin\", 3600) > 10000"}'::jsonb,
  '550e8400-e29b-41d4-a716-446655440000'::uuid,
  NOW(),
  NOW()
)
ON CONFLICT (id) DO NOTHING;

-- Rule 7: Complex Multi-Condition Rule (Block action, high priority)
INSERT INTO rules (id, name, description, action, priority, enabled, conditions, schema_id, created_at, updated_at)
VALUES (
  '660e8400-e29b-41d4-a716-446655440007'::uuid,
  'Block High-Risk Multi-Condition Transactions',
  'Example: Block transactions that meet multiple high-risk criteria simultaneously',
  'block',
  15,
  true,
  '{"custom_expression": "(amount > 5000 and currency != \"USD\") or (amount > 10000 and user.country in [\"RU\", \"CN\", \"KP\"])"}'::jsonb,
  '550e8400-e29b-41d4-a716-446655440000'::uuid,
  NOW(),
  NOW()
)
ON CONFLICT (id) DO NOTHING;

-- Rule 8: Allow Low-Risk Transactions (Allow action, low priority)
INSERT INTO rules (id, name, description, action, priority, enabled, conditions, schema_id, created_at, updated_at)
VALUES (
  '660e8400-e29b-41d4-a716-446655440008'::uuid,
  'Allow Low-Risk Verified Transactions',
  'Example: Automatically approve low-value, verified transactions from trusted sources',
  'allow',
  100,
  true,
  '{"custom_expression": "amount < 1000 and is_verified == true and user.verification_status == \"verified\""}'::jsonb,
  '550e8400-e29b-41d4-a716-446655440000'::uuid,
  NOW(),
  NOW()
)
ON CONFLICT (id) DO NOTHING;
