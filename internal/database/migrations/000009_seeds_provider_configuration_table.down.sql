-- Remove the inserted provider configurations from the provider_configurations table

DELETE FROM provider_configurations WHERE country_id IN 
    (SELECT id FROM countries WHERE country_code IN ('US', 'GB', 'JP', 'CN', 'AE', 'IN', 'BR', 'AU')) 
    AND currency_id IN 
    (SELECT id FROM currencies WHERE currency_code IN ('USD', 'GBP', 'JPY', 'CNY', 'AED', 'INR', 'AUD'));
