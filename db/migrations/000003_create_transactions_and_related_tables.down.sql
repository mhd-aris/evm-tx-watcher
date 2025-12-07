-- Drop triggers if any were created on transactions
DROP TRIGGER IF EXISTS transactions_updated_at ON transactions;
DROP FUNCTION IF EXISTS set_transactions_updated_at();

-- Drop indexes
DROP INDEX IF EXISTS idx_transactions_hash;
DROP INDEX IF EXISTS idx_transactions_block_number;
DROP INDEX IF EXISTS idx_token_transfers_tx_id;
DROP INDEX IF EXISTS idx_token_transfers_token_address;

-- Drop tables (reverse order of dependencies)
DROP TABLE IF EXISTS webhook_deliveries;
DROP TABLE IF EXISTS token_transfers;
DROP TABLE IF EXISTS transactions;
