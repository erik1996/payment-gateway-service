package payment

import (
	"context"
	"fmt"
	"testing"
	"time"

	"payment-gateway-service/internal/provider"
	"payment-gateway-service/internal/utils"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Helper function to create a mock GORM DB and service
func setupTest(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("an error '%s' was not expected when initializing gorm", err)
	}

	return gormDB, mock, func() {
		db.Close()
	}
}

func TestCreatePayment_Success(t *testing.T) {
	gormDB, mock, teardown := setupTest(t)
	defer teardown()

	// Setup mock dependencies
	providerSvc := new(MockProviderService)
	adapterFactory := new(MockAdapterFactory)
	paymentService := NewPaymentService(gormDB, providerSvc, adapterFactory)

	// Setup expectations for SQL queries
	sqlRows := sqlmock.NewRows([]string{"id"}).AddRow(1)
	mock.ExpectBegin()
	mock.ExpectQuery(`^INSERT INTO "payments" \("amount","payment_type","status","currency_code","user_id","provider_id","external_id","created_at","updated_at"\) VALUES \(\$1,\$2,\$3,\$4,\$5,\$6,\$7,\$8,\$9\) RETURNING "id"$`).
		WithArgs(
			float64(100),     // Amount
			"DEPOSIT",        // PaymentType
			"INITIALIZED",    // Status
			"USD",            // CurrencyCode
			1,                // UserID
			1,                // ProviderID
			"",               // ExternalID
			sqlmock.AnyArg(), // CreatedAt
			sqlmock.AnyArg(), // UpdatedAt
		).
		WillReturnRows(sqlRows)
	mock.ExpectExec(`^UPDATE "payments" SET "amount"=\$1,"payment_type"=\$2,"status"=\$3,"currency_code"=\$4,"user_id"=\$5,"provider_id"=\$6,"external_id"=\$7,"created_at"=\$8,"updated_at"=\$9 WHERE "id" = \$10$`).
		WithArgs(
			float64(100),     // Amount
			"DEPOSIT",        // PaymentType
			"PENDING",        // Status
			"USD",            // CurrencyCode
			1,                // UserID
			1,                // ProviderID
			"external-id",    // ExternalID
			sqlmock.AnyArg(), // CreatedAt
			sqlmock.AnyArg(), // UpdatedAt
			"1",              // ID
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	// Mock expectations for adapter and provider service
	mockAdapter := new(MockProviderAdapter)
	mockAdapter.On("GetDetails", context.TODO(), float64(100), "DEPOSIT", "USD", "US").Return("http://payment.url", "external-id", nil)
	adapterFactory.On("GetAdapter", context.TODO(), "USD", "US").Return(mockAdapter, nil)
	providerSvc.On("FindProviderConfig", context.TODO(), "USD", "US").Return(&provider.ProviderConfiguration{
		ProviderID: 1,
	}, nil)

	// Call the method under test
	url, err := paymentService.CreatePayment(context.TODO(), &PaymentRequest{
		UserID:       1,
		Amount:       float64(100),
		CurrencyCode: "USD",
		CountryCode:  "US",
	}, utils.PaymentTypeDeposit)

	assert.NoError(t, err)
	assert.Equal(t, "http://payment.url", url)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreatePayment_Failure_Insert(t *testing.T) {
	gormDB, mock, teardown := setupTest(t)
	defer teardown()

	// Setup mock dependencies
	providerSvc := new(MockProviderService)
	adapterFactory := new(MockAdapterFactory)
	paymentService := NewPaymentService(gormDB, providerSvc, adapterFactory)

	// Setup expectations for SQL queries
	mock.ExpectBegin()
	mock.ExpectQuery(`^INSERT INTO "payments" \("amount","payment_type","status","currency_code","user_id","provider_id","external_id","created_at","updated_at"\) VALUES \(\$1,\$2,\$3,\$4,\$5,\$6,\$7,\$8,\$9\) RETURNING "id"$`).
		WithArgs(
			float64(100),     // Amount
			"DEPOSIT",        // PaymentType
			"INITIALIZED",    // Status
			"USD",            // CurrencyCode
			1,                // UserID
			1,                // ProviderID
			"",               // ExternalID
			sqlmock.AnyArg(), // CreatedAt
			sqlmock.AnyArg(), // UpdatedAt
		).
		WillReturnError(fmt.Errorf("insert error"))

	mock.ExpectRollback()

	// Setup mock expectations for provider service
	providerSvc.On("FindProviderConfig", context.TODO(), "USD", "US").Return(&provider.ProviderConfiguration{
		ProviderID: 1,
	}, nil)

	// Setup mock expectations for adapter
	mockAdapter := new(MockProviderAdapter)
	adapterFactory.On("GetAdapter", context.TODO(), "USD", "US").Return(mockAdapter, nil)
	mockAdapter.On("GetDetails", context.TODO(), float64(100), "DEPOSIT", "USD", "US").Return("", "", nil)

	// Call the method under test
	url, err := paymentService.CreatePayment(context.TODO(), &PaymentRequest{
		UserID:       1,
		Amount:       float64(100),
		CurrencyCode: "USD",
		CountryCode:  "US",
	}, utils.PaymentTypeDeposit)

	assert.Error(t, err)
	assert.Empty(t, url)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreatePayment_Failure_FindProviderConfig(t *testing.T) {
	gormDB, _, teardown := setupTest(t)
	defer teardown()

	// Setup mock dependencies
	providerSvc := new(MockProviderService)
	adapterFactory := new(MockAdapterFactory)
	paymentService := NewPaymentService(gormDB, providerSvc, adapterFactory)

	// Setup mock to return an error for FindProviderConfig
	providerSvc.On("FindProviderConfig", context.TODO(), "USD", "US").Return(nil, fmt.Errorf("find provider config error"))

	// Call the method under test
	url, err := paymentService.CreatePayment(context.TODO(), &PaymentRequest{
		UserID:       1,
		Amount:       float64(100),
		CurrencyCode: "USD",
		CountryCode:  "US",
	}, utils.PaymentTypeDeposit)

	assert.Error(t, err)
	assert.Empty(t, url)
}

func TestCreatePayment_Failure_GetAdapter(t *testing.T) {
	gormDB, _, teardown := setupTest(t)
	defer teardown()

	// Setup mock dependencies
	providerSvc := new(MockProviderService)
	adapterFactory := new(MockAdapterFactory)
	paymentService := NewPaymentService(gormDB, providerSvc, adapterFactory)

	// Setup mock expectations
	providerSvc.On("FindProviderConfig", context.TODO(), "USD", "US").Return(&provider.ProviderConfiguration{
		ProviderID: 1,
	}, nil)
	adapterFactory.On("GetAdapter", context.TODO(), "USD", "US").Return(nil, fmt.Errorf("get adapter error"))

	// Call the method under test
	url, err := paymentService.CreatePayment(context.TODO(), &PaymentRequest{
		UserID:       1,
		Amount:       float64(100),
		CurrencyCode: "USD",
		CountryCode:  "US",
	}, utils.PaymentTypeDeposit)

	assert.Error(t, err)
	assert.Empty(t, url)
}

func TestCreatePayment_Failure_GetDetails(t *testing.T) {
	gormDB, _, teardown := setupTest(t)
	defer teardown()

	// Setup mock dependencies
	providerSvc := new(MockProviderService)
	adapterFactory := new(MockAdapterFactory)
	paymentService := NewPaymentService(gormDB, providerSvc, adapterFactory)

	// Setup mock expectations
	providerSvc.On("FindProviderConfig", context.TODO(), "USD", "US").Return(&provider.ProviderConfiguration{
		ProviderID: 1,
	}, nil)
	mockAdapter := new(MockProviderAdapter)
	adapterFactory.On("GetAdapter", context.TODO(), "USD", "US").Return(mockAdapter, nil)
	mockAdapter.On("GetDetails", context.TODO(), float64(100), "DEPOSIT", "USD", "US").Return("", "", fmt.Errorf("get details error"))

	// Call the method under test
	url, err := paymentService.CreatePayment(context.TODO(), &PaymentRequest{
		UserID:       1,
		Amount:       float64(100),
		CurrencyCode: "USD",
		CountryCode:  "US",
	}, utils.PaymentTypeDeposit)

	assert.Error(t, err)
	assert.Empty(t, url)
}

func TestCreatePayment_Failure_Update(t *testing.T) {
	gormDB, mock, teardown := setupTest(t)
	defer teardown()

	// Setup mock dependencies
	providerSvc := new(MockProviderService)
	adapterFactory := new(MockAdapterFactory)
	paymentService := NewPaymentService(gormDB, providerSvc, adapterFactory)

	// Setup mock expectations
	sqlRows := sqlmock.NewRows([]string{"id"}).AddRow(1)
	mock.ExpectBegin()
	mock.ExpectQuery(`^INSERT INTO "payments" \("amount","payment_type","status","currency_code","user_id","provider_id","external_id","created_at","updated_at"\) VALUES \(\$1,\$2,\$3,\$4,\$5,\$6,\$7,\$8,\$9\) RETURNING "id"$`).
		WithArgs(
			float64(100),     // Amount
			"DEPOSIT",        // PaymentType
			"INITIALIZED",    // Status
			"USD",            // CurrencyCode
			1,                // UserID
			1,                // ProviderID
			"",               // ExternalID
			sqlmock.AnyArg(), // CreatedAt
			sqlmock.AnyArg(), // UpdatedAt
		).
		WillReturnRows(sqlRows)
	mock.ExpectExec(`^UPDATE "payments" SET "amount"=\$1,"payment_type"=\$2,"status"=\$3,"currency_code"=\$4,"user_id"=\$5,"provider_id"=\$6,"external_id"=\$7,"created_at"=\$8,"updated_at"=\$9 WHERE "id" = \$10$`).
		WithArgs(
			float64(100),     // Amount
			"DEPOSIT",        // PaymentType
			"PENDING",        // Status
			"USD",            // CurrencyCode
			1,                // UserID
			1,                // ProviderID
			"external-id",    // ExternalID
			sqlmock.AnyArg(), // CreatedAt
			sqlmock.AnyArg(), // UpdatedAt
			"1",              // ID
		).
		WillReturnError(fmt.Errorf("update error"))
	mock.ExpectRollback()

	// Setup mock expectations for provider service and adapter
	mockAdapter := new(MockProviderAdapter)
	mockAdapter.On("GetDetails", context.TODO(), float64(100), "DEPOSIT", "USD", "US").Return("http://payment.url", "external-id", nil)
	adapterFactory.On("GetAdapter", context.TODO(), "USD", "US").Return(mockAdapter, nil)
	providerSvc.On("FindProviderConfig", context.TODO(), "USD", "US").Return(&provider.ProviderConfiguration{
		ProviderID: 1,
	}, nil)

	// Call the method under test
	url, err := paymentService.CreatePayment(context.TODO(), &PaymentRequest{
		UserID:       1,
		Amount:       float64(100),
		CurrencyCode: "USD",
		CountryCode:  "US",
	}, utils.PaymentTypeDeposit)

	assert.Error(t, err)
	assert.Empty(t, url)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdatePayment_Success(t *testing.T) {
	gormDB, mock, teardown := setupTest(t)
	defer teardown()

	// Setup mock expectations
	sqlmock.NewRows([]string{"id"}).AddRow(1)
	mock.ExpectBegin()
	mock.ExpectExec(`^UPDATE "payments" SET "amount"=\$1,"payment_type"=\$2,"status"=\$3,"currency_code"=\$4,"user_id"=\$5,"provider_id"=\$6,"external_id"=\$7,"created_at"=\$8,"updated_at"=\$9 WHERE "id" = \$10$`).
		WithArgs(
			100.0,            // Amount
			"DEPOSIT",        // PaymentType
			"PENDING",        // Status
			"USD",            // CurrencyCode
			1,                // UserID
			1,                // ProviderID
			"external-id",    // ExternalID
			sqlmock.AnyArg(), // CreatedAt
			sqlmock.AnyArg(), // UpdatedAt
			"1",              // ID
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Create a payment object with test data
	payment := &Payment{
		ID:           "1",
		Amount:       100.0,
		PaymentType:  "DEPOSIT",
		Status:       "PENDING",
		CurrencyCode: "USD",
		UserID:       1,
		ProviderID:   1,
		ExternalID:   "external-id",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Setup the payment service
	paymentService := NewPaymentService(gormDB, nil, nil)

	// Call the method under test
	err := paymentService.UpdatePayment(payment)

	// Check if there were no errors
	assert.NoError(t, err)

	// Check if all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdatePayment_Failure(t *testing.T) {
	gormDB, mock, teardown := setupTest(t)
	defer teardown()

	// Setup mock expectations for a failure scenario
	mock.ExpectBegin()
	mock.ExpectExec(`^UPDATE "payments" SET "amount"=\$1,"payment_type"=\$2,"status"=\$3,"currency_code"=\$4,"user_id"=\$5,"provider_id"=\$6,"external_id"=\$7,"created_at"=\$8,"updated_at"=\$9 WHERE "id" = \$10$`).
		WithArgs(
			100.0,            // Amount
			"DEPOSIT",        // PaymentType
			"PENDING",        // Status
			"USD",            // CurrencyCode
			1,                // UserID
			1,                // ProviderID
			"external-id",    // ExternalID
			sqlmock.AnyArg(), // CreatedAt
			sqlmock.AnyArg(), // UpdatedAt
			"1",              // ID
		).
		WillReturnError(fmt.Errorf("update error"))
	mock.ExpectRollback()

	// Create a payment object with test data
	payment := &Payment{
		ID:           "1",
		Amount:       100.0,
		PaymentType:  "DEPOSIT",
		Status:       "PENDING",
		CurrencyCode: "USD",
		UserID:       1,
		ProviderID:   1,
		ExternalID:   "external-id",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Setup the payment service
	paymentService := NewPaymentService(gormDB, nil, nil)

	// Call the method under test
	err := paymentService.UpdatePayment(payment)

	// Check if an error occurred
	assert.Error(t, err)

	// Check if all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

// Mock implementations
type MockProviderService struct {
	mock.Mock
}

func (m *MockProviderService) FindProviderConfig(ctx context.Context, currencyCode, countryCode string) (*provider.ProviderConfiguration, error) {
	args := m.Called(ctx, currencyCode, countryCode)
	return args.Get(0).(*provider.ProviderConfiguration), args.Error(1)
}

type MockAdapterFactory struct {
	mock.Mock
}

func (m *MockAdapterFactory) GetAdapter(ctx context.Context, currencyCode, countryCode string) (provider.ProviderAdapter, error) {
	args := m.Called(ctx, currencyCode, countryCode)
	return args.Get(0).(provider.ProviderAdapter), args.Error(1)
}

type MockProviderAdapter struct {
	mock.Mock
}

func (m *MockProviderAdapter) GetDetails(ctx context.Context, amount float64, paymentType string, currencyCode, countryCode string) (string, string, error) {
	args := m.Called(ctx, amount, paymentType, currencyCode, countryCode)
	return args.String(0), args.String(1), args.Error(2)
}
