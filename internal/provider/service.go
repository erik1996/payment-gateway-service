package provider

import (
	"context"
	"fmt"
	"payment-gateway-service/internal/utils"

	"gorm.io/gorm"
)

// ProviderServiceInterface defines the methods that the ProviderService must implement.
type ProviderServiceInterface interface {
	FindProviderByName(ctx context.Context, name string) (*Provider, error)
	FindProviderConfig(ctx context.Context, currencyCode, countryCode string) (*ProviderConfiguration, error)
}

// ProviderService handles operations related to payment providers.
type ProviderService struct {
	db *gorm.DB
}

// NewProviderService initializes a new ProviderService with the provided database connection.
func NewProviderService(db *gorm.DB) *ProviderService {
	return &ProviderService{db: db}
}

// FindProviderByName retrieves a provider by name from the database.
func (s *ProviderService) FindProviderByName(ctx context.Context, name string) (*Provider, error) {
	var provider Provider

	utils.LogWithRequestID(ctx, fmt.Sprintf("ProviderService: Attempting to find provider with name: %s", name))

	// Query the database for the provider with the specified name
	if err := s.db.Where("name = ?", name).First(&provider).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.LogWithRequestID(ctx, fmt.Sprintf("ProviderService: Provider with name %s not found.", name))
		} else {
			utils.LogWithRequestID(ctx, fmt.Sprintf("ProviderService: Error finding provider with name %s: %v", name, err))
		}
		return nil, err
	}

	utils.LogWithRequestID(ctx, fmt.Sprintf("ProviderService: Successfully found provider: %+v", provider))
	return &provider, nil
}

// FindProviderConfig retrieves the provider configuration based on currency code, country code, and priority.
func (s *ProviderService) FindProviderConfig(ctx context.Context, currencyCode, countryCode string) (*ProviderConfiguration, error) {
	var providerConfig ProviderConfiguration

	utils.LogWithRequestID(ctx, fmt.Sprintf("ProviderService: Attempting to find provider configuration for CurrencyCode: %s, CountryCode: %s", currencyCode, countryCode))

	// Perform the join query safely with parameterized inputs
	err := s.db.
		Table("provider_configurations").
		Joins("JOIN currencies ON currencies.id = provider_configurations.currency_id").
		Joins("JOIN countries ON countries.id = provider_configurations.country_id").
		Joins("JOIN payment_providers ON payment_providers.id = provider_configurations.provider_id").
		Select("provider_configurations.*, payment_providers.name as provider_name").
		Where("currencies.currency_code = ? AND countries.country_code = ?", currencyCode, countryCode).
		Order("provider_configurations.priority ASC, provider_configurations.id").
		First(&providerConfig).Error

	// Handle the case where no matching configuration is found
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.LogWithRequestID(ctx, fmt.Sprintf("ProviderService: No provider configuration found for CurrencyCode: %s, CountryCode: %s", currencyCode, countryCode))
		} else {
			utils.LogWithRequestID(ctx, fmt.Sprintf("ProviderService: Error retrieving provider configuration: %v", err))
		}
		return nil, err
	}

	utils.LogWithRequestID(ctx, fmt.Sprintf("ProviderService: Successfully found provider configuration for CurrencyCode: %s, CountryCode: %s with Provider: %s", currencyCode, countryCode, providerConfig.ProviderName))
	return &providerConfig, nil
}

// Ensure ProviderService implements ProviderServiceInterface.
var _ ProviderServiceInterface = (*ProviderService)(nil)
