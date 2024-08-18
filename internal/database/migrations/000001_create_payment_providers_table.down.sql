-- Drop the payment_providers table and associated trigger function
DROP TRIGGER IF EXISTS set_timestamp ON payment_providers;
DROP FUNCTION IF EXISTS trigger_set_timestamp();
DROP TABLE IF EXISTS payment_providers;
