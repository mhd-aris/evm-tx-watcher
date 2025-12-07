-- Drop trigger and function
DROP TRIGGER IF EXISTS trg_set_addresses_updated_at ON addresses;
DROP FUNCTION IF EXISTS set_addresses_updated_at;

-- Drop indexes
DROP INDEX IF EXISTS idx_addresses_address_chain;
DROP INDEX IF EXISTS idx_addresses_user_id;

-- Drop table
DROP TABLE IF EXISTS addresses;
