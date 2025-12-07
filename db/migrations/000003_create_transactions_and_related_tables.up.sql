-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Transactions table
CREATE TABLE transactions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    hash TEXT NOT NULL UNIQUE,
    block_number BIGINT NOT NULL,
    block_hash TEXT NOT NULL,
    transaction_index INT NOT NULL,
    chain_id BIGINT NOT NULL,
    from_address TEXT NOT NULL,
    to_address TEXT,
    value NUMERIC(78,0), -- Ethereum uint256
    gas_used BIGINT,
    gas_price NUMERIC(78,0),
    tx_type INT NOT NULL,
    status INT NOT NULL, -- 1=success, 0=failed
    block_timestamp TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Composite index for block scan
CREATE INDEX idx_transactions_blocknumber_txindex ON transactions(block_number, transaction_index);
-- Lookup by hash
CREATE INDEX idx_transactions_hash ON transactions(hash);
-- Lookup by address
CREATE INDEX idx_transactions_from_address ON transactions(from_address);
CREATE INDEX idx_transactions_to_address ON transactions(to_address);
-- Filter by time
CREATE INDEX idx_transactions_block_timestamp ON transactions(block_timestamp);

-- TokenTransfers table
CREATE TABLE token_transfers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    transaction_id UUID NOT NULL REFERENCES transactions(id) ON DELETE CASCADE,
    log_index INT NOT NULL,
    token_address TEXT NOT NULL,
    from_address TEXT NOT NULL,
    to_address TEXT NOT NULL,
    value NUMERIC(78,0),
    token_decimals INT,
    token_symbol TEXT,
    token_name TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(transaction_id, log_index) -- avoid duplicate events
);

-- Indexes for token transfers
CREATE INDEX idx_token_transfers_transaction_id ON token_transfers(transaction_id);
CREATE INDEX idx_token_transfers_token_address ON token_transfers(token_address);
CREATE INDEX idx_token_transfers_from_address ON token_transfers(from_address);
CREATE INDEX idx_token_transfers_to_address ON token_transfers(to_address);

-- WebhookDeliveries table
CREATE TABLE webhook_deliveries (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    webhook_id UUID NOT NULL,
    transaction_id UUID NOT NULL REFERENCES transactions(id) ON DELETE CASCADE,
    payload JSONB NOT NULL,
    status TEXT NOT NULL, -- pending, delivered, failed
    http_status_code INT,
    response_body TEXT,
    error_message TEXT,
    retry_count INT NOT NULL DEFAULT 0,
    max_retries INT NOT NULL DEFAULT 3,
    next_retry_at TIMESTAMPTZ,
    delivered_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Indexes for webhook deliveries
CREATE INDEX idx_webhook_deliveries_webhook_id ON webhook_deliveries(webhook_id);
CREATE INDEX idx_webhook_deliveries_transaction_id ON webhook_deliveries(transaction_id);
CREATE INDEX idx_webhook_deliveries_status ON webhook_deliveries(status);

-- Trigger function to auto-update updated_at
CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Apply triggers
CREATE TRIGGER trg_set_transactions_updated_at
BEFORE UPDATE ON transactions
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER trg_set_token_transfers_updated_at
BEFORE UPDATE ON token_transfers
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER trg_set_webhook_deliveries_updated_at
BEFORE UPDATE ON webhook_deliveries
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();
