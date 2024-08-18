-- Drop the payments table and associated trigger function and types
DROP TRIGGER IF EXISTS set_timestamp ON payments;
DROP TABLE IF EXISTS payments;
DROP TYPE IF EXISTS payment_type;
DROP TYPE IF EXISTS payment_status;
