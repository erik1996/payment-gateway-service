package provider

import (
	"context"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
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

func TestFindProviderByName(t *testing.T) {
	gormDB, mock, teardown := setupTest(t)
	defer teardown()

	providerService := NewProviderService(gormDB)

	// Set up mock expectations
	sqlRows := sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "TestProvider")
	mock.ExpectQuery(`^SELECT \* FROM "payment_providers" WHERE name = \$1 ORDER BY "payment_providers"."id" LIMIT \$2$`).
		WithArgs("TestProvider", 1).
		WillReturnRows(sqlRows)

	// Call the service method
	result, err := providerService.FindProviderByName(context.TODO(), "TestProvider")

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "TestProvider", result.Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFindProviderByName_NotFound(t *testing.T) {
	gormDB, mock, teardown := setupTest(t)
	defer teardown()

	providerService := NewProviderService(gormDB)

	// Set up mock expectations for not found case
	mock.ExpectQuery(`^SELECT \* FROM "payment_providers" WHERE name = \$1 ORDER BY "payment_providers"."id" LIMIT \$2$`).
		WithArgs("TestProvider", 1).
		WillReturnError(gorm.ErrRecordNotFound)

	// Call the service method
	result, err := providerService.FindProviderByName(context.TODO(), "TestProvider")

	// Assertions
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
	assert.Nil(t, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFindProviderByName_DBError(t *testing.T) {
	gormDB, mock, teardown := setupTest(t)
	defer teardown()

	providerService := NewProviderService(gormDB)

	// Set up mock expectations for a database error
	mock.ExpectQuery(`^SELECT \* FROM "payment_providers" WHERE name = \$1 ORDER BY "payment_providers"."id" LIMIT \$2$`).
		WithArgs("TestProvider", 1).
		WillReturnError(fmt.Errorf("database error"))

	// Call the service method
	result, err := providerService.FindProviderByName(context.TODO(), "TestProvider")

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFindProviderConfig(t *testing.T) {
	gormDB, mock, teardown := setupTest(t)
	defer teardown()

	providerService := NewProviderService(gormDB)

	// Set up mock expectations for a successful query
	sqlRows := sqlmock.NewRows([]string{"currency_id", "country_id", "provider_id", "provider_name"}).
		AddRow(1, 1, 1, "TestProvider")
	mock.ExpectQuery(`^SELECT provider_configurations\.\*, payment_providers\.name as provider_name FROM "provider_configurations"`).
		WillReturnRows(sqlRows)

	// Call the service method
	result, err := providerService.FindProviderConfig(context.TODO(), "USD", "US")

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "TestProvider", result.ProviderName)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFindProviderConfig_NotFound(t *testing.T) {
	gormDB, mock, teardown := setupTest(t)
	defer teardown()

	providerService := NewProviderService(gormDB)

	// Set up mock expectations for a not found case
	mock.ExpectQuery(`^SELECT provider_configurations\.\*, payment_providers\.name as provider_name FROM "provider_configurations"`).
		WillReturnError(gorm.ErrRecordNotFound)

	// Call the service method
	result, err := providerService.FindProviderConfig(context.TODO(), "USD", "US")

	// Assertions
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
	assert.Nil(t, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFindProviderConfig_DBError(t *testing.T) {
	gormDB, mock, teardown := setupTest(t)
	defer teardown()

	providerService := NewProviderService(gormDB)

	// Set up mock expectations for a database error
	mock.ExpectQuery(`^SELECT provider_configurations\.\*, payment_providers\.name as provider_name FROM "provider_configurations"`).
		WillReturnError(fmt.Errorf("database error"))

	// Call the service method
	result, err := providerService.FindProviderConfig(context.TODO(), "USD", "US")

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}
