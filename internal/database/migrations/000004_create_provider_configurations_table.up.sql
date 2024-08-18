CREATE TABLE provider_configurations (
    id SERIAL PRIMARY KEY,
    country_id INT REFERENCES countries(id) ON DELETE CASCADE,
    currency_id INT REFERENCES currencies(id) ON DELETE CASCADE,
    provider_id INT REFERENCES payment_providers(id) ON DELETE CASCADE,
    priority INT NOT NULL CHECK (priority >= 1),
    base_url VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE (country_id, currency_id, provider_id)
);

CREATE OR REPLACE FUNCTION trigger_set_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER set_timestamp
BEFORE UPDATE ON provider_configurations
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();
