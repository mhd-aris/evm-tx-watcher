-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Table: addresses
CREATE TABLE IF NOT EXISTS addresses (
    id UUID PRIMARY KEY,                       
    address VARCHAR(255) NOT NULL,           
    chain_id INT NOT NULL,                   
    is_contract BOOLEAN NOT NULL DEFAULT FALSE,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    label VARCHAR(255),                        
    description TEXT,                          
    user_id UUID,                              
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Index for unique address per chain
CREATE UNIQUE INDEX IF NOT EXISTS idx_addresses_address_chain
    ON addresses (address, chain_id);

-- Index for user_id
CREATE INDEX IF NOT EXISTS idx_addresses_user_id
    ON addresses (user_id);

-- Trigger function to update updated_at timestamp
CREATE OR REPLACE FUNCTION set_addresses_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger to call the function before update
CREATE TRIGGER trg_set_addresses_updated_at
BEFORE UPDATE ON addresses
FOR EACH ROW
EXECUTE FUNCTION set_addresses_updated_at();
