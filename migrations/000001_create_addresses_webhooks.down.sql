-- Drop triggers
DROP TRIGGER IF EXISTS update_addresses_updated_at ON addresses;
DROP TRIGGER IF EXISTS update_webhooks_updated_at ON webhooks;

-- Drop trigger function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop indexes
DROP INDEX IF EXISTS idx_addresses_address;
DROP INDEX IF EXISTS idx_addresses_chain_id;
DROP INDEX IF EXISTS idx_addresses_is_active;
DROP INDEX IF EXISTS idx_addresses_created_at;
DROP INDEX IF EXISTS idx_webhooks_address_id;

-- Drop tables (order matters due to foreign key constraints)
DROP TABLE IF EXISTS webhooks;
DROP TABLE IF EXISTS addresses;