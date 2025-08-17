-- Create addresses table
CREATE TABLE IF NOT EXISTS addresses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    address VARCHAR(42) NOT NULL,
    chain_id INTEGER NOT NULL,
    is_contract BOOLEAN NOT NULL DEFAULT FALSE,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    label VARCHAR(100),
    description VARCHAR(255),
    user_id UUID,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    
    -- Constraints
    CONSTRAINT unique_address_per_chain UNIQUE (address, chain_id),
    CONSTRAINT valid_address_format CHECK (address ~ '^0x[a-fA-F0-9]{40}$'),
    CONSTRAINT positive_chain_id CHECK (chain_id > 0)
);

-- Create webhooks table
CREATE TABLE IF NOT EXISTS webhooks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    address_id UUID NOT NULL,
    url VARCHAR(2048) NOT NULL,
    secret VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    
    -- Foreign key constraint
    CONSTRAINT fk_webhooks_address_id FOREIGN KEY (address_id) REFERENCES addresses(id) ON DELETE CASCADE,
    
    -- Constraints
    CONSTRAINT valid_webhook_url CHECK (url ~ '^https?://'),
    CONSTRAINT min_secret_length CHECK (LENGTH(secret) >= 10)
);

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_addresses_address ON addresses(address);
CREATE INDEX IF NOT EXISTS idx_addresses_chain_id ON addresses(chain_id);
CREATE INDEX IF NOT EXISTS idx_addresses_is_active ON addresses(is_active);
CREATE INDEX IF NOT EXISTS idx_addresses_created_at ON addresses(created_at);
CREATE INDEX IF NOT EXISTS idx_webhooks_address_id ON webhooks(address_id);

-- Create updated_at trigger function
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create triggers to automatically update updated_at
CREATE TRIGGER update_addresses_updated_at
    BEFORE UPDATE ON addresses
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_webhooks_updated_at
    BEFORE UPDATE ON webhooks
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();