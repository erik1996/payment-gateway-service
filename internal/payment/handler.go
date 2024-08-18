package payment

import (
	"fmt"
	"net/http"
	"payment-gateway-service/config"
	"payment-gateway-service/internal/provider"
	"payment-gateway-service/internal/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// PaymentHandler handles payment-related requests
type PaymentHandler struct {
	service PaymentServiceInterface
}

// NewPaymentHandler initializes a new PaymentHandler
func NewPaymentHandler(db *gorm.DB) *PaymentHandler {
	providerSvc := provider.NewProviderService(db)
	adapterFactory := provider.NewAdapterFactory(providerSvc)
	service := NewPaymentService(db, providerSvc, adapterFactory)
	return &PaymentHandler{service: service}
}

// Deposit handles deposit requests
// @Summary Handles deposit requests
// @Description Processes a deposit request and returns a URL for payment.
// @Tags payment
// @Accept json
// @Produce json
// @Param X-AUTH-TOKEN header string true "Authorization token"
// @Param validatedBody body PaymentRequest true "Validated Payment Request"
// @Success 200 {object} map[string]interface{} "url"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 500 {object} map[string]interface{} "Failed to process request"
// @Router /payment/deposit [post]
// @Param exampleRequest body PaymentRequest true "Example request" Example({"amount": 40, "country_code": "US", "currency_code": "USD" "user_id": 1})
func (h *PaymentHandler) Deposit(c *gin.Context) {

	// Proceed with processing the deposit
	h.processPayment(c, utils.PaymentTypeDeposit)
}

// Withdrawal handles withdrawal requests
// @Summary Handles withdrawal requests
// @Description Processes a withdrawal request and returns a URL for payment.
// @Tags payment
// @Accept json
// @Produce json
// @Param X-AUTH-TOKEN header string true "Authorization token"
// @Param validatedBody body PaymentRequest true "Validated Payment Request"
// @Success 200 {object} map[string]interface{} "url"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 500 {object} map[string]interface{} "Failed to process request"
// @Router /payment/withdrawal [post]
// @Param exampleRequest body PaymentRequest true "Example request" Example({"amount": 40, "country_code": "US", "currency_code": "USD", "user_id": 1})
func (h *PaymentHandler) Withdrawal(c *gin.Context) {

	// Proceed with processing the deposit
	h.processPayment(c, utils.PaymentTypeWithdrawal)
}

// processPayment handles the common logic for deposit and withdrawal
func (h *PaymentHandler) processPayment(c *gin.Context, paymentType utils.PaymentType) {
	utils.LogWithRequestID(c, "Processing payment request")

	req, exists := c.Get("validatedBody")
	if !exists {
		utils.LogWithRequestID(c, "Invalid request: no validated body found")
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request", nil)
		return
	}

	utils.LogWithRequestID(c, "Validated body exists")
	paymentRequest, ok := req.(*PaymentRequest)
	if !ok {
		utils.LogWithRequestID(c, "Failed to process request: type assertion failed")
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to process request", nil)
		return
	}

	// Log the details of the payment request
	utils.LogWithRequestID(c, fmt.Sprintf("Received payment request: UserID=%d, Amount=%.2f, CurrencyCode=%s, CountryCode=%s, PaymentType=%s",
		paymentRequest.UserID, paymentRequest.Amount, paymentRequest.CurrencyCode, paymentRequest.CountryCode, paymentType))

	// Create the payment using the service and get the URL
	url, err := h.service.CreatePayment(c, paymentRequest, paymentType)
	if err != nil {
		utils.LogWithRequestID(c, fmt.Sprintf("Failed to create payment: %v", err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create payment", nil)
		return
	}

	utils.LogWithRequestID(c, fmt.Sprintf("%s created with URL: %s", paymentType, url))
	utils.SuccessResponse(c, http.StatusOK, fmt.Sprintf("%s successful", paymentType), gin.H{"url": url})
}

// HandleSuccessCallback handles successful payment provider callbacks (with and without external_id)
// @Summary Handles successful payment provider callbacks
// @Description Processes a successful payment callback and redirects to a status URL.
// @Tags payment
// @Produce json
// @Param external_id path string true "External ID"
// @Success 302 {string} string "Redirects to status URL"
// @Failure 400 {object} map[string]interface{} "Error extracting external ID"
// @Failure 500 {object} map[string]interface{} "Failed to handle callback"
// @Router /payment/callback/success [get]
func (h *PaymentHandler) HandleSuccessCallback(c *gin.Context) {
	h.handleCallback(c, utils.PaymentStatusSuccess, "successful")
}

// HandleFailedCallback handles failed payment provider callbacks (with and without external_id)
// @Summary Handles failed payment provider callbacks
// @Description Processes a failed payment callback and redirects to a status URL.
// @Tags payment
// @Produce json
// @Param external_id path string true "External ID"
// @Success 302 {string} string "Redirects to status URL"
// @Failure 400 {object} map[string]interface{} "Error extracting external ID"
// @Failure 500 {object} map[string]interface{} "Failed to handle callback"
// @Router /payment/callback/failure [get]
func (h *PaymentHandler) HandleFailedCallback(c *gin.Context) {
	h.handleCallback(c, utils.PaymentStatusFailed, "failed")
}

// handleCallback handles the common logic for both success and failed callbacks
func (h *PaymentHandler) handleCallback(c *gin.Context, status utils.PaymentStatus, result string) {
	utils.LogWithRequestID(c, fmt.Sprintf("Handling %s callback", result))

	externalID, err := ExtractExternalID(c)
	if err != nil {
		utils.LogWithRequestID(c, fmt.Sprintf("Error extracting external ID: %v", err))
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	// Handle the callback using the service
	payment, err := h.service.HandleCallback(c, externalID, status)
	if err != nil {
		utils.LogWithRequestID(c, fmt.Sprintf("Failed to handle %s callback: %v", result, err))
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": fmt.Sprintf("Failed to handle %s callback", result)})
		return
	}

	utils.LogWithRequestID(c, fmt.Sprintf("%s callback handled for payment: %+v", result, payment))
	cfg := config.LoadConfig()
	// Redirect user to the desired URL
	redirectURL := fmt.Sprintf("%s/payment?status=%s&id=%v", cfg.AppHost, result, payment.ID)
	c.Redirect(http.StatusFound, redirectURL)
}

func (h *PaymentHandler) PaymentStatus(c *gin.Context) {
	// Extract the status and id query parameters from the URL
	status := c.Query("status")
	id := c.Query("id")

	// For simplicity
	c.JSON(http.StatusOK, gin.H{
		"status": status,
		"id":     id,
	})
}
