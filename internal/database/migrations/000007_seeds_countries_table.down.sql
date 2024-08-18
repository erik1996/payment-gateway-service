-- Remove the inserted countries from the countries table

DELETE FROM countries WHERE country_code IN ('US', 'CA', 'GB', 'AU', 'DE', 'FR', 'JP', 'CN', 'IN', 'AE');
