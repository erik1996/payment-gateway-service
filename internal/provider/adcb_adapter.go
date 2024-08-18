package provider

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"payment-gateway-service/internal/utils"
	"time"
)

type ADCBAdapter struct {
	baseURL    string
	userID     string
	userSecret string
}

func NewADCBAdapter(baseURL string) *ADCBAdapter {
	userID := os.Getenv("ADCB_USER_ID")
	userSecret := os.Getenv("ADCB_USER_SECRET")

	return &ADCBAdapter{
		baseURL:    baseURL,
		userID:     userID,
		userSecret: userSecret,
	}
}

// Define the PaymentRequest structure with correct XML tags
type ADCBPaymentRequest struct {
	XMLName     xml.Name `xml:"PaymentRequest"`
	Amount      float64  `xml:"Amount"`
	PaymentType string   `xml:"PaymentType"`
	Currency    string   `xml:"Currency"`
	Country     string   `xml:"Country"`
}

// Define the PaymentResponse structure with correct XML tags
type ADCBPaymentResponse struct {
	XMLName    xml.Name `xml:"PaymentResponse"`
	URL        string   `xml:"URL"`
	ExternalID string   `xml:"ExternalID"`
}

func (a *ADCBAdapter) GetDetails(ctx context.Context, amount float64, paymentType, currencyCode, countryCode string) (string, string, error) {
	startTime := time.Now() // Capture the start time
	utils.LogWithRequestID(ctx, fmt.Sprintf("ADCB Adapter: Starting to generate payment details for Amount: %.2f, Payment Type: %s, currencyCode: %s, countryCode: %s", amount, paymentType, currencyCode, countryCode))

	// Defer the logging of the elapsed time until the function returns
	defer func() {
		utils.LogWithRequestID(ctx, fmt.Sprintf("ADCB Adapter: Completed generating payment details in %v for Amount: %.2f, Payment Type: %s, currencyCode: %s, countryCode: %s", time.Since(startTime), amount, paymentType, currencyCode, countryCode))
	}()
	// Append the specific endpoint to the baseURL
	requestURL := fmt.Sprintf("%s/adcb/payment", a.baseURL)

	// Prepare the request body
	paymentRequest := ADCBPaymentRequest{
		Amount:      amount,
		PaymentType: paymentType,
		Currency:    currencyCode,
		Country:     countryCode,
	}

	// Marshal the request to XML with XML declaration
	xmlHeader := `<?xml version="1.0" encoding="UTF-8"?>`
	soapRequestBody, err := xml.Marshal(paymentRequest)
	if err != nil {
		utils.LogWithRequestID(ctx, "ADCB Adapter: Failed to marshal request")
		return "", "", err
	}

	// Combine XML header and body
	soapRequestBody = []byte(fmt.Sprintf("%s\n%s", xmlHeader, string(soapRequestBody)))

	// Prepare the HTTP request
	request, err := http.NewRequestWithContext(ctx, "POST", requestURL, bytes.NewBuffer(soapRequestBody))
	if err != nil {
		utils.LogWithRequestID(ctx, "ADCB Adapter: Failed to create HTTP request")
		return "", "", err
	}

	// Set the headers for authentication
	request.Header.Set("Content-Type", "application/xml")
	request.Header.Set("user_id", a.userID)
	request.Header.Set("user_secret", a.userSecret)

	// Log the request details
	utils.LogWithRequestID(ctx, fmt.Sprintf("ADCB Adapter: Request URL: %s", requestURL))
	utils.LogWithRequestID(ctx, fmt.Sprintf("ADCB Adapter: Request Headers: %v", request.Header))
	utils.LogWithRequestID(ctx, fmt.Sprintf("ADCB Adapter: Request Body: %s", string(soapRequestBody)))

	// Perform the HTTP request
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		utils.LogWithRequestID(ctx, "ADCB Adapter: Failed to perform HTTP request")
		return "", "", err
	}
	defer resp.Body.Close()

	// Check the HTTP status code
	if resp.StatusCode != http.StatusOK {
		// Read and log the response body
		responseBody, _ := ioutil.ReadAll(resp.Body)
		utils.LogWithRequestID(ctx, fmt.Sprintf("ADCB Adapter: HTTP request failed with status code %d, Response Body: %s", resp.StatusCode, string(responseBody)))
		return "", "", fmt.Errorf("HTTP request failed with status code %d", resp.StatusCode)
	}

	// Read and log the response body
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		utils.LogWithRequestID(ctx, "ADCB Adapter: Failed to read response body")
		return "", "", err
	}
	utils.LogWithRequestID(ctx, fmt.Sprintf("ADCB Adapter: Response body: %s", string(responseBody)))

	// Parse the response
	var paymentResponse ADCBPaymentResponse
	if err := xml.Unmarshal(responseBody, &paymentResponse); err != nil {
		utils.LogWithRequestID(ctx, "ADCB Adapter: Failed to unmarshal response")
		return "", "", err
	}

	// Extract URL and ExternalID
	if paymentResponse.URL == "" || paymentResponse.ExternalID == "" {
		utils.LogWithRequestID(ctx, "ADCB Adapter: Missing URL or ExternalID in response")
		return "", "", fmt.Errorf("failed to get payment details")
	}

	utils.LogWithRequestID(ctx, fmt.Sprintf("ADCB Adapter: Successfully received response with URL: %s and ExternalID: %s", paymentResponse.URL, paymentResponse.ExternalID))

	return paymentResponse.URL, paymentResponse.ExternalID, nil
}
