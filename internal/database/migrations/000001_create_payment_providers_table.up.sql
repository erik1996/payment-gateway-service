-- Create the payment_providers table
CREATE TABLE payment_providers (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Create the trigger function for updated_at
CREATE OR REPLACE FUNCTION trigger_set_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create the trigger to update the updated_at timestamp on modification
CREATE TRIGGER set_timestamp
BEFORE UPDATE ON payment_providers
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();
