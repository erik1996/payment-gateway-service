-- Create the uuid-ossp extension if not already enabled
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create the ENUM types for payment type and status
CREATE TYPE payment_type AS ENUM ('DEPOSIT', 'WITHDRAWAL');
CREATE TYPE payment_status AS ENUM ('INITIALIZED', 'PENDING', 'SUCCESS', 'FAILED');

-- Create the payments table
CREATE TABLE payments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    amount NUMERIC(12, 2) NOT NULL,
    payment_type payment_type NOT NULL,
    status payment_status DEFAULT 'INITIALIZED', 
    currency_code VARCHAR(3) NOT NULL, 
    user_id INT NOT NULL,
    provider_id INT REFERENCES payment_providers(id) ON DELETE SET NULL,
    external_id VARCHAR(255), 
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(external_id, provider_id)
);

-- Trigger to automatically update the updated_at column on row update
CREATE OR REPLACE FUNCTION trigger_set_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER set_timestamp
BEFORE UPDATE ON payments
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();
