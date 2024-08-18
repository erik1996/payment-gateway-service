CREATE TABLE currencies (
    id SERIAL PRIMARY KEY,
    currency_name VARCHAR(255) NOT NULL,
    currency_code VARCHAR(3) NOT NULL UNIQUE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE OR REPLACE FUNCTION trigger_set_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER set_timestamp
BEFORE UPDATE ON currencies
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();
