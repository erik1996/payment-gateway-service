package provider_test

import (
	"context"
	"errors"
	"payment-gateway-service/internal/provider"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Define a custom type for context keys to prevent collisions
type contextKey string

const requestIDKey = contextKey("requestID")

// MockProviderService is a mock implementation of the ProviderServiceInterface.
type MockProviderService struct {
	mock.Mock
}

func (m *MockProviderService) FindProviderByName(ctx context.Context, name string) (*provider.Provider, error) {
	args := m.Called(ctx, name)
	if args.Get(0) != nil {
		return args.Get(0).(*provider.Provider), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockProviderService) FindProviderConfig(ctx context.Context, currencyCode, countryCode string) (*provider.ProviderConfiguration, error) {
	args := m.Called(ctx, currencyCode, countryCode)
	if args.Get(0) != nil {
		return args.Get(0).(*provider.ProviderConfiguration), args.Error(1)
	}
	return nil, args.Error(1)
}

// setupTest initializes the AdapterFactory with a mocked ProviderServiceInterface.
func setupTest(_ *testing.T) (*provider.AdapterFactory, *MockProviderService) {
	// Create a mock provider service
	mockProviderService := new(MockProviderService)

	// Inject the mock service into the AdapterFactory
	factory := provider.NewAdapterFactory(mockProviderService)

	return factory, mockProviderService
}

func TestAdapterFactory_GetAdapter_HSBC(t *testing.T) {
	// Add RequestID to the context using a custom key type
	ctx := context.WithValue(context.Background(), requestIDKey, "test-request-id")

	// Set up the AdapterFactory and mock service
	factory, mockProviderService := setupTest(t)

	// Define the expected provider configuration for HSBC
	expectedConfig := &provider.ProviderConfiguration{
		ProviderName: "HSBC",
		BaseURL:      "https://hsbc.example.com",
	}

	// Set up the mock expectation for FindProviderConfig
	mockProviderService.On("FindProviderConfig", ctx, "USD", "US").Return(expectedConfig, nil)

	// Call the method under test
	adapter, err := factory.GetAdapter(ctx, "USD", "US")

	// Assert that the adapter was created successfully and without error
	assert.NoError(t, err)
	assert.NotNil(t, adapter)

	// Verify the correct adapter type was returned
	_, ok := adapter.(*provider.HSBCAdapter)
	assert.True(t, ok, "Expected adapter to be of type HSBCAdapter")

	// Ensure all expectations were met
	mockProviderService.AssertExpectations(t)
}

func TestAdapterFactory_GetAdapter_ADCB(t *testing.T) {
	// Add RequestID to the context using a custom key type
	ctx := context.WithValue(context.Background(), requestIDKey, "test-request-id")

	// Set up the AdapterFactory and mock service
	factory, mockProviderService := setupTest(t)

	// Define the expected provider configuration for ADCB
	expectedConfig := &provider.ProviderConfiguration{
		ProviderName: "ADCB",
		BaseURL:      "https://adcb.example.com",
	}

	// Set up the mock expectation for FindProviderConfig
	mockProviderService.On("FindProviderConfig", ctx, "AED", "AE").Return(expectedConfig, nil)

	// Call the method under test
	adapter, err := factory.GetAdapter(ctx, "AED", "AE")

	// Assert that the adapter was created successfully and without error
	assert.NoError(t, err)
	assert.NotNil(t, adapter)

	// Verify the correct adapter type was returned
	_, ok := adapter.(*provider.ADCBAdapter)
	assert.True(t, ok, "Expected adapter to be of type ADCBAdapter")

	// Ensure all expectations were met
	mockProviderService.AssertExpectations(t)
}

func TestAdapterFactory_GetAdapter_UnsupportedProvider(t *testing.T) {
	// Add RequestID to the context using a custom key type
	ctx := context.WithValue(context.Background(), requestIDKey, "test-request-id")

	// Set up the AdapterFactory and mock service
	factory, mockProviderService := setupTest(t)

	// Define the expected provider configuration for an unsupported provider
	expectedConfig := &provider.ProviderConfiguration{
		ProviderName: "Unsupported",
		BaseURL:      "https://unsupported.example.com",
	}

	// Set up the mock expectation for FindProviderConfig
	mockProviderService.On("FindProviderConfig", ctx, "GBP", "GB").Return(expectedConfig, nil)

	// Call the method under test
	adapter, err := factory.GetAdapter(ctx, "GBP", "GB")

	// Assert that the adapter creation failed due to unsupported provider
	assert.Error(t, err)
	assert.Nil(t, adapter)
	assert.EqualError(t, err, "provider not supported")

	// Ensure all expectations were met
	mockProviderService.AssertExpectations(t)
}

func TestAdapterFactory_GetAdapter_ErrorInProviderService(t *testing.T) {
	// Add RequestID to the context using a custom key type
	ctx := context.WithValue(context.Background(), requestIDKey, "test-request-id")

	// Set up the AdapterFactory and mock service
	factory, mockProviderService := setupTest(t)

	// Set up the mock expectation for FindProviderConfig to return an error
	mockProviderService.On("FindProviderConfig", ctx, "JPY", "JP").Return(nil, errors.New("database error"))

	// Call the method under test
	adapter, err := factory.GetAdapter(ctx, "JPY", "JP")

	// Assert that the adapter creation failed due to service error
	assert.Error(t, err)
	assert.Nil(t, adapter)
	assert.EqualError(t, err, "database error")

	// Ensure all expectations were met
	mockProviderService.AssertExpectations(t)
}
