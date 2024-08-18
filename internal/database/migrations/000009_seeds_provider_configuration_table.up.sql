-- Insert provider configurations into the provider_configurations table

INSERT INTO provider_configurations (country_id, currency_id, provider_id, base_url, priority, created_at, updated_at) VALUES 
-- HSBC configurations
((SELECT id FROM countries WHERE country_code = 'US'), (SELECT id FROM currencies WHERE currency_code = 'USD'), (SELECT id FROM payment_providers WHERE name = 'HSBC'), 'http://hsbc:8081', 1, NOW(), NOW()),
((SELECT id FROM countries WHERE country_code = 'GB'), (SELECT id FROM currencies WHERE currency_code = 'GBP'), (SELECT id FROM payment_providers WHERE name = 'HSBC'), 'http://hsbc:8081', 1, NOW(), NOW()),
((SELECT id FROM countries WHERE country_code = 'JP'), (SELECT id FROM currencies WHERE currency_code = 'JPY'), (SELECT id FROM payment_providers WHERE name = 'HSBC'), 'http://hsbc:8081', 1, NOW(), NOW()),
((SELECT id FROM countries WHERE country_code = 'CN'), (SELECT id FROM currencies WHERE currency_code = 'CNY'), (SELECT id FROM payment_providers WHERE name = 'HSBC'), 'http://hsbc:8081', 1, NOW(), NOW()),

-- ADCB configurations
((SELECT id FROM countries WHERE country_code = 'AE'), (SELECT id FROM currencies WHERE currency_code = 'AED'), (SELECT id FROM payment_providers WHERE name = 'ADCB'), 'http://adcb:8082', 1, NOW(), NOW()),
((SELECT id FROM countries WHERE country_code = 'IN'), (SELECT id FROM currencies WHERE currency_code = 'INR'), (SELECT id FROM payment_providers WHERE name = 'ADCB'), 'http://adcb:8082', 1, NOW(), NOW()),
((SELECT id FROM countries WHERE country_code = 'AU'), (SELECT id FROM currencies WHERE currency_code = 'AUD'), (SELECT id FROM payment_providers WHERE name = 'ADCB'), 'http://adcb:8082', 1, NOW(), NOW());
