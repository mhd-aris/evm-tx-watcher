-- DROP TRIGGER AND FUNCTION
DROP TRIGGER IF EXISTS trg_set_webhooks_updated_at ON webhooks;
DROP FUNCTION IF EXISTS set_webhooks_updated_at;

-- DROP INDEX
DROP INDEX IF EXISTS idx_webhooks_address_id;

-- DROP TABLE
DROP TABLE IF EXISTS webhooks;
