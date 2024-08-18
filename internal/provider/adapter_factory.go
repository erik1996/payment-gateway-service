package provider

import (
	"context"
	"errors"
	"payment-gateway-service/internal/utils"
)

// AdapterFactoryInterface defines the method that the AdapterFactory must implement.
type AdapterFactoryInterface interface {
	GetAdapter(ctx context.Context, currencyCode, countryCode string) (ProviderAdapter, error)
}

// AdapterFactory is responsible for creating provider adapters based on the provider configuration.
type AdapterFactory struct {
	providerService ProviderServiceInterface
}

// NewAdapterFactory initializes a new AdapterFactory with a ProviderServiceInterface.
func NewAdapterFactory(providerService ProviderServiceInterface) *AdapterFactory {
	return &AdapterFactory{providerService: providerService}
}

// GetAdapter returns the appropriate adapter based on the currency code, country code, and priority.
func (f *AdapterFactory) GetAdapter(ctx context.Context, currencyCode, countryCode string) (ProviderAdapter, error) {
	utils.LogWithRequestID(ctx, "AdapterFactory: Attempting to retrieve adapter")

	// Find the appropriate provider configuration based on currency code, country code, and priority.
	providerConfig, err := f.providerService.FindProviderConfig(ctx, currencyCode, countryCode)
	if err != nil {
		utils.LogWithRequestID(ctx, "AdapterFactory: Error retrieving provider configuration")
		return nil, err
	}

	providerName := providerConfig.ProviderName
	utils.LogWithRequestID(ctx, "AdapterFactory: Found provider: "+providerName)

	// Pass the baseURL from the database to the appropriate adapter.
	switch providerName {
	case "HSBC":
		utils.LogWithRequestID(ctx, "AdapterFactory: Creating HSBCAdapter")
		return NewHSBCAdapter(providerConfig.BaseURL), nil
	case "ADCB":
		utils.LogWithRequestID(ctx, "AdapterFactory: Creating ADCBAdapter")
		return NewADCBAdapter(providerConfig.BaseURL), nil
	default:
		utils.LogWithRequestID(ctx, "AdapterFactory: Unsupported provider: "+providerName)
		return nil, errors.New("provider not supported")
	}
}
