-- Migration: Add comprehensive test data
-- This migration adds:
-- 1. A payment event schema that uses all field types and features
-- 2. Rules that demonstrate all possible combinations of actions, expressions, and priorities

-- Insert comprehensive payment event schema
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

-- Rule 5: Velocity Count Check (Review action, medium-high priority)
INSERT INTO rules (id, name, description, action, priority, enabled, conditions, schema_id, created_at, updated_at)
VALUES (
  '660e8400-e29b-41d4-a716-446655440005'::uuid,
  'High Transaction Frequency',
  'Example: Review accounts with more than 10 transactions in the last hour',
  'review',
  40,
  true,
  '{"custom_expression": "velocityCount(origin, 3600) > 10"}'::jsonb,
  '550e8400-e29b-41d4-a716-446655440000'::uuid,
  NOW(),
  NOW()
)
ON CONFLICT (id) DO NOTHING;

-- Rule 6: Velocity Sum Check (Block action, high priority)
INSERT INTO rules (id, name, description, action, priority, enabled, conditions, schema_id, created_at, updated_at)
VALUES (
  '660e8400-e29b-41d4-a716-446655440006'::uuid,
  'Block High Cumulative Amount',
  'Example: Block accounts with cumulative transaction amount over $50,000 in the last 24 hours',
  'block',
  15,
  true,
  '{"custom_expression": "velocitySum(origin, 86400) > 50000"}'::jsonb,
  '550e8400-e29b-41d4-a716-446655440000'::uuid,
  NOW(),
  NOW()
)
ON CONFLICT (id) DO NOTHING;

-- Rule 7: Polygon Geographic Check (Review action, medium priority)
INSERT INTO rules (id, name, description, action, priority, enabled, conditions, schema_id, created_at, updated_at)
VALUES (
  '660e8400-e29b-41d4-a716-446655440007'::uuid,
  'Review Restricted Geographic Zone',
  'Example: Review transactions originating from a restricted geographic polygon area',
  'review',
  55,
  true,
  '{"custom_expression": "pointInPolygon(location.lat, location.lon, [[37.7749, -122.4194], [37.7849, -122.4094], [37.7649, -122.4294], [37.7549, -122.4394]])"}'::jsonb,
  '550e8400-e29b-41d4-a716-446655440000'::uuid,
  NOW(),
  NOW()
)
ON CONFLICT (id) DO NOTHING;

-- Rule 8: Complex Expression with AND (Review action, medium priority)
INSERT INTO rules (id, name, description, action, priority, enabled, conditions, schema_id, created_at, updated_at)
VALUES (
  '660e8400-e29b-41d4-a716-446655440008'::uuid,
  'Review High Value International',
  'Example: Review high-value transactions to international destinations',
  'review',
  45,
  true,
  '{"custom_expression": "amount > 5000 and currency != \"USD\""}'::jsonb,
  '550e8400-e29b-41d4-a716-446655440000'::uuid,
  NOW(),
  NOW()
)
ON CONFLICT (id) DO NOTHING;

-- Rule 9: Complex Expression with OR (Review action, medium priority)
INSERT INTO rules (id, name, description, action, priority, enabled, conditions, schema_id, created_at, updated_at)
VALUES (
  '660e8400-e29b-41d4-a716-446655440009'::uuid,
  'Review High Risk Combinations',
  'Example: Review transactions that are either high-value USD or from high-risk geographic regions',
  'review',
  50,
  true,
  '{"custom_expression": "(amount > 10000 and currency == \"USD\") or (amount > 5000 and (pointInPolygon(location.lat, location.lon, [[55.0, 30.0], [60.0, 30.0], [60.0, 45.0], [55.0, 45.0]]) or pointInPolygon(location.lat, location.lon, [[35.0, 110.0], [45.0, 110.0], [45.0, 125.0], [35.0, 125.0]])))"}'::jsonb,
  '550e8400-e29b-41d4-a716-446655440000'::uuid,
  NOW(),
  NOW()
)
ON CONFLICT (id) DO NOTHING;

-- Rule 10: Boolean Check (Allow action, low priority - whitelist)
INSERT INTO rules (id, name, description, action, priority, enabled, conditions, schema_id, created_at, updated_at)
VALUES (
  '660e8400-e29b-41d4-a716-446655440010'::uuid,
  'Allow Verified Premium Users',
  'Example: Explicitly allow transactions from verified premium users',
  'allow',
  90,
  true,
  '{"custom_expression": "is_verified == true and user.account_type == \"premium\" and user.verification_status == \"verified\""}'::jsonb,
  '550e8400-e29b-41d4-a716-446655440000'::uuid,
  NOW(),
  NOW()
)
ON CONFLICT (id) DO NOTHING;

-- Rule 11: Array Contains Check (Review action, medium priority)
INSERT INTO rules (id, name, description, action, priority, enabled, conditions, schema_id, created_at, updated_at)
VALUES (
  '660e8400-e29b-41d4-a716-446655440011'::uuid,
  'Review Tagged Transactions',
  'Example: Review transactions tagged as high-value or urgent',
  'review',
  65,
  true,
  '{"custom_expression": "\"high-value\" in tags or \"urgent\" in tags"}'::jsonb,
  '550e8400-e29b-41d4-a716-446655440000'::uuid,
  NOW(),
  NOW()
)
ON CONFLICT (id) DO NOTHING;

-- Rule 12: Nested Field Comparison (Review action, medium priority)
INSERT INTO rules (id, name, description, action, priority, enabled, conditions, schema_id, created_at, updated_at)
VALUES (
  '660e8400-e29b-41d4-a716-446655440012'::uuid,
  'Review High Risk Score',
  'Example: Review transactions with high metadata risk score',
  'review',
  70,
  true,
  '{"custom_expression": "metadata.risk_score > 75"}'::jsonb,
  '550e8400-e29b-41d4-a716-446655440000'::uuid,
  NOW(),
  NOW()
)
ON CONFLICT (id) DO NOTHING;

-- Rule 13: Multiple Conditions with NOT (Review action, medium priority)
INSERT INTO rules (id, name, description, action, priority, enabled, conditions, schema_id, created_at, updated_at)
VALUES (
  '660e8400-e29b-41d4-a716-446655440013'::uuid,
  'Review Unverified High Value',
  'Example: Review high-value transactions from unverified users',
  'review',
  35,
  true,
  '{"custom_expression": "amount > 5000 and not is_verified"}'::jsonb,
  '550e8400-e29b-41d4-a716-446655440000'::uuid,
  NOW(),
  NOW()
)
ON CONFLICT (id) DO NOTHING;

-- Rule 14: Payment Method Check (Review action, medium priority)
INSERT INTO rules (id, name, description, action, priority, enabled, conditions, schema_id, created_at, updated_at)
VALUES (
  '660e8400-e29b-41d4-a716-446655440014'::uuid,
  'Review Non-Standard Payment Methods',
  'Example: Review transactions using non-standard payment methods',
  'review',
  75,
  true,
  '{"custom_expression": "payment_method.type != \"credit_card\" and payment_method.type != \"debit_card\""}'::jsonb,
  '550e8400-e29b-41d4-a716-446655440000'::uuid,
  NOW(),
  NOW()
)
ON CONFLICT (id) DO NOTHING;

-- Rule 15: Disabled Rule Example (Block action, high priority but disabled)
INSERT INTO rules (id, name, description, action, priority, enabled, conditions, schema_id, created_at, updated_at)
VALUES (
  '660e8400-e29b-41d4-a716-446655440015'::uuid,
  'Disabled: Block All International',
  'Example: Example of a disabled rule - would block all international transactions if enabled',
  'block',
  20,
  false,
  '{"custom_expression": "currency != \"USD\""}'::jsonb,
  '550e8400-e29b-41d4-a716-446655440000'::uuid,
  NOW(),
  NOW()
)
ON CONFLICT (id) DO NOTHING;

-- Rule 16: Low Priority Allow Rule (Allow action, very low priority)
INSERT INTO rules (id, name, description, action, priority, enabled, conditions, schema_id, created_at, updated_at)
VALUES (
  '660e8400-e29b-41d4-a716-446655440016'::uuid,
  'Allow Low Value Domestic',
  'Example: Explicitly allow low-value domestic transactions from US geographic region',
  'allow',
  95,
  true,
  '{"custom_expression": "amount < 1000 and currency == \"USD\" and pointInPolygon(location.lat, location.lon, [[25.0, -125.0], [49.0, -125.0], [49.0, -66.0], [25.0, -66.0]])"}'::jsonb,
  '550e8400-e29b-41d4-a716-446655440000'::uuid,
  NOW(),
  NOW()
)
ON CONFLICT (id) DO NOTHING;

-- Rule 17: Complex Multi-Condition with Nested Fields (Review action, medium priority)
INSERT INTO rules (id, name, description, action, priority, enabled, conditions, schema_id, created_at, updated_at)
VALUES (
  '660e8400-e29b-41d4-a716-446655440017'::uuid,
  'Review Complex Risk Pattern',
  'Example: Review transactions matching complex risk patterns combining multiple factors including geographic regions',
  'review',
  30,
  true,
  '{"custom_expression": "(amount > 3000 and metadata.risk_score > 50) or ((pointInPolygon(location.lat, location.lon, [[55.0, 30.0], [60.0, 30.0], [60.0, 45.0], [55.0, 45.0]]) or pointInPolygon(location.lat, location.lon, [[35.0, 110.0], [45.0, 110.0], [45.0, 125.0], [35.0, 125.0]])) and not is_verified) or (velocityCount(origin, 7200) > 5 and amount > 2000)"}'::jsonb,
  '550e8400-e29b-41d4-a716-446655440000'::uuid,
  NOW(),
  NOW()
)
ON CONFLICT (id) DO NOTHING;

