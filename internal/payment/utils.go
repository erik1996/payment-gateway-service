package payment

import (
	"fmt"
	"payment-gateway-service/internal/utils"

	"github.com/gin-gonic/gin"
)

// PaymentRequest represents the request payload for a payment
type PaymentRequest struct {
	UserID       int     `json:"user_id" binding:"required"`
	Amount       float64 `json:"amount" binding:"required,gt=1"`
	CurrencyCode string  `json:"currency_code" binding:"required,len=3"`
	CountryCode  string  `json:"country_code" binding:"required,len=2"`
}

// ExtractExternalID extracts the external ID from path or query parameters
func ExtractExternalID(c *gin.Context) (string, error) {
	// Try to get external_id from path parameter
	if externalID := c.Param("external_id"); externalID != "" {
		utils.LogWithRequestID(c, fmt.Sprintf("Extracted external ID from path: %s", externalID))
		return externalID, nil
	}

	// Try to get external_id from query parameters "id" or "externalId"
	if externalID := c.Query("id"); externalID != "" {
		utils.LogWithRequestID(c, fmt.Sprintf("Extracted external ID from query (id): %s", externalID))
		return externalID, nil
	}

	if externalID := c.Query("externalId"); externalID != "" {
		utils.LogWithRequestID(c, fmt.Sprintf("Extracted external ID from query (externalId): %s", externalID))
		return externalID, nil
	}

	err := fmt.Errorf("external ID is required")
	utils.LogWithRequestID(c, fmt.Sprintf("Error: %v", err))
	return "", err
}
