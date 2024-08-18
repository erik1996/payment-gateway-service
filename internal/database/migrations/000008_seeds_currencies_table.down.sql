-- Remove the inserted currencies from the currencies table

DELETE FROM currencies WHERE currency_code IN ('USD', 'EUR', 'GBP', 'AUD', 'CAD', 'JPY', 'CNY', 'INR', 'BRL', 'AED');
