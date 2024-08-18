package payment

import (
	"context"
	"errors"
	"fmt"
	"payment-gateway-service/internal/provider"
	"payment-gateway-service/internal/utils"
	"time"

	"gorm.io/gorm"
)

// PaymentServiceInterface defines the methods that the PaymentService must implement.
type PaymentServiceInterface interface {
	CreatePayment(ctx context.Context, paymentRequest *PaymentRequest, paymentType utils.PaymentType) (string, error)
	HandleCallback(ctx context.Context, externalID string, status utils.PaymentStatus) (*Payment, error)
	UpdatePayment(payment *Payment) error
	FindPaymentByExternalID(externalID string) (*Payment, error)
}

// ProviderServiceInterface defines the methods that the ProviderService must implement.
type ProviderServiceInterface interface {
	FindProviderConfig(ctx context.Context, currencyCode, countryCode string) (*provider.ProviderConfiguration, error)
}

// AdapterFactoryInterface defines the method that the AdapterFactory must implement.
type AdapterFactoryInterface interface {
	GetAdapter(ctx context.Context, currencyCode, countryCode string) (provider.ProviderAdapter, error)
}

// PaymentService handles operations related to payments.
type PaymentService struct {
	db             *gorm.DB
	providerSvc    ProviderServiceInterface
	adapterFactory AdapterFactoryInterface
}

// NewPaymentService initializes a new PaymentService.
func NewPaymentService(db *gorm.DB, providerSvc ProviderServiceInterface, adapterFactory AdapterFactoryInterface) *PaymentService {
	return &PaymentService{
		db:             db,
		providerSvc:    providerSvc,
		adapterFactory: adapterFactory,
	}
}

// CreatePayment creates a new payment in the database and returns the URL for further processing.
func (s *PaymentService) CreatePayment(ctx context.Context, paymentRequest *PaymentRequest, paymentType utils.PaymentType) (string, error) {
	utils.LogWithRequestID(ctx, "PaymentService: Starting payment creation")

	var url string

	err := s.db.Transaction(func(tx *gorm.DB) error {
		// Find the appropriate provider configuration.
		providerConfig, err := s.providerSvc.FindProviderConfig(ctx, paymentRequest.CurrencyCode, paymentRequest.CountryCode)
		if err != nil {
			utils.LogWithRequestID(ctx, "PaymentService: Failed to find provider configuration")
			return errors.New("failed to find provider configuration")
		}

		// Get the right adapter using the factory.
		adapter, err := s.adapterFactory.GetAdapter(ctx, paymentRequest.CurrencyCode, paymentRequest.CountryCode)
		if err != nil {
			utils.LogWithRequestID(ctx, "PaymentService: Failed to get adapter for provider")
			return err
		}

		// Create a new payment record with the initial status.
		payment := &Payment{
			UserID:       paymentRequest.UserID,
			Amount:       paymentRequest.Amount,
			PaymentType:  paymentType,
			Status:       utils.PaymentStatusInitialized,
			CurrencyCode: paymentRequest.CurrencyCode,
			ProviderID:   providerConfig.ProviderID,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		// Save the payment in the database.
		if err := tx.Create(payment).Error; err != nil {
			utils.LogWithRequestID(ctx, "PaymentService: Failed to save payment to the database")
			return err
		}

		// Generate payment details using the adapter.
		var externalID string
		url, externalID, err = adapter.GetDetails(ctx, payment.Amount, string(paymentType), payment.CurrencyCode, paymentRequest.CountryCode)
		if err != nil {
			utils.LogWithRequestID(ctx, "PaymentService: Failed to generate payment details using adapter")
			return err
		}

		// Update the payment record with the external ID and status to "Pending".
		payment.ExternalID = externalID
		payment.Status = utils.PaymentStatusPending // Set status to "Pending"
		if err := tx.Save(payment).Error; err != nil {
			utils.LogWithRequestID(ctx, "PaymentService: Failed to update payment with external ID and pending status")
			return err
		}

		utils.LogWithRequestID(ctx, "PaymentService: Payment created successfully")
		return nil
	})

	if err != nil {
		return "", err
	}

	return url, nil
}

// HandleCallback processes callbacks from payment providers and updates the payment status and user balance.
func (s *PaymentService) HandleCallback(ctx context.Context, externalID string, status utils.PaymentStatus) (*Payment, error) {
	utils.LogWithRequestID(ctx, "PaymentService: Handling callback for ExternalID")

	var payment *Payment
	err := s.db.Transaction(func(tx *gorm.DB) error {
		// Find the payment by the external ID within the transaction.
		if err := tx.Where("external_id = ?", externalID).First(&payment).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.LogWithRequestID(ctx, "PaymentService: Payment not found with ExternalID")
				return errors.New("payment not found")
			}
			utils.LogWithRequestID(ctx, "PaymentService: Failed to find payment with ExternalID")
			return err
		}

		// Check if the current status is "Pending". If not, do not update.
		if payment.Status != utils.PaymentStatusPending {
			utils.LogWithRequestID(ctx, fmt.Sprintf("PaymentService: Payment status is not pending (current status: %s), no update performed", payment.Status))
			return errors.New("payment status is not pending, update skipped")
		}

		// Handle the callback based on the status.
		payment.Status = status
		payment.UpdatedAt = time.Now()
		if err := tx.Save(payment).Error; err != nil {
			utils.LogWithRequestID(ctx, "PaymentService: Failed to update payment status")
			return err
		}

		utils.LogWithRequestID(ctx, "PaymentService: Payment status updated successfully")
		return nil
	})

	if err != nil {
		return nil, err
	}

	return payment, nil
}

// UpdatePayment updates an existing payment in the database.
func (s *PaymentService) UpdatePayment(payment *Payment) error {
	payment.UpdatedAt = time.Now()
	if err := s.db.Save(payment).Error; err != nil {
		return err
	}
	return nil
}

// FindPaymentByExternalID finds a payment by its external ID.
func (s *PaymentService) FindPaymentByExternalID(externalID string) (*Payment, error) {
	var payment Payment
	if err := s.db.Where("external_id = ?", externalID).First(&payment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &payment, nil
}

// Ensure PaymentService implements PaymentServiceInterface.
var _ PaymentServiceInterface = (*PaymentService)(nil)
