package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"payment-gateway-service/internal/utils"
	"time"
)

type HSBCAdapter struct {
	baseURL    string
	userID     string
	userSecret string
}

func NewHSBCAdapter(baseURL string) *HSBCAdapter {
	userID := os.Getenv("HSBC_USER_ID")
	userSecret := os.Getenv("HSBC_USER_SECRET")

	return &HSBCAdapter{
		baseURL:    baseURL,
		userID:     userID,
		userSecret: userSecret,
	}
}

type HSBCResponse struct {
	URL        string `json:"url"`
	ExternalID string `json:"external_id"`
}

func (a *HSBCAdapter) GetDetails(ctx context.Context, amount float64, paymentType, currencyCode, countryCode string) (string, string, error) {
	startTime := time.Now() // Capture the start time
	utils.LogWithRequestID(ctx, fmt.Sprintf("HSBC Adapter: Starting to generate payment details for Amount: %.2f, Transaction Type: %s, Currency Code: %s, Country Code: %s", amount, paymentType, currencyCode, countryCode))

	// Defer the logging of the elapsed time until the function returns
	defer func() {
		elapsedTime := time.Since(startTime)
		utils.LogWithRequestID(ctx, fmt.Sprintf("HSBC Adapter: Completed generating payment details in %v for Amount: %.2f, Transaction Type: %s, Currency Code: %s, Country Code: %s", elapsedTime, amount, paymentType, currencyCode, countryCode))
	}()

	select {
	case <-ctx.Done():
		utils.LogWithRequestID(ctx, fmt.Sprintf("HSBC Adapter: Request cancelled or timed out for Amount: %.2f, Transaction Type: %s, Currency Code: %s. Error: %v", amount, paymentType, currencyCode, ctx.Err()))
		return "", "", ctx.Err()
	default:
		utils.LogWithRequestID(ctx, fmt.Sprintf("HSBC Adapter: Context is active. Continuing with the payment generation process for Amount: %.2f, Transaction Type: %s, Currency Code: %s", amount, paymentType, currencyCode))
	}

	requestURL := fmt.Sprintf("%s/hsbc/payment", a.baseURL)
	reqBody := map[string]interface{}{
		"amount":       amount,
		"payment_type": paymentType,
		"currency":     currencyCode,
		"country":      countryCode,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		utils.LogWithRequestID(ctx, "HSBC Adapter: Failed to marshal request body to JSON")
		return "", "", err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", requestURL, bytes.NewBuffer(jsonData))
	if err != nil {
		utils.LogWithRequestID(ctx, "HSBC Adapter: Failed to create new HTTP request")
		return "", "", err
	}

	req.Header.Set("user_id", a.userID)
	req.Header.Set("user_secret", a.userSecret)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		utils.LogWithRequestID(ctx, "HSBC Adapter: Failed to perform HTTP request")
		return "", "", err
	}
	defer resp.Body.Close()

	// Read and log the raw response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		utils.LogWithRequestID(ctx, "HSBC Adapter: Failed to read response body")
		return "", "", err
	}
	utils.LogWithRequestID(ctx, fmt.Sprintf("HSBC Adapter: Raw response body: %s", string(body)))

	var hsbcResponse HSBCResponse
	if err := json.Unmarshal(body, &hsbcResponse); err != nil {
		utils.LogWithRequestID(ctx, "HSBC Adapter: Failed to decode response from HSBC service")
		return "", "", err
	}

	utils.LogWithRequestID(ctx, fmt.Sprintf("HSBC Adapter: Successfully received response with URL: %s and ExternalID: %s", hsbcResponse.URL, hsbcResponse.ExternalID))

	return hsbcResponse.URL, hsbcResponse.ExternalID, nil
}
