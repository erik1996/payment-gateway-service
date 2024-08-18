package provider

import "context"

// ProviderAdapter is the interface that all provider adapters must implement
type ProviderAdapter interface {
	GetDetails(ctx context.Context, amount float64, transactionType, currencyCode string, countryCode string) (string, string, error)
}
