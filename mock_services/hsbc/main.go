package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
)

// PaymentRequest represents the structure of the payment request
type PaymentRequest struct {
	Amount      float64 `json:"amount"`
	PaymentType string  `json:"payment_type"`
	Currency    string  `json:"currency"`
	Country     string  `json:"country"`
}

// PaymentResponse represents the structure of the payment response
type PaymentResponse struct {
	URL        string `json:"url"`
	ExternalID string `json:"external_id"`
}

// CallbackRequest represents the structure of the callback request
type CallbackRequest struct {
	ExternalID string `json:"external_id"`
	Status     string `json:"status"`
}

func main() {
	http.HandleFunc("/hsbc/payment", handleHSBCPayment)
	http.HandleFunc("/hsbc/callback", handleHSBCCallback)
	log.Println("HSBC Mock Service running on port 8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}

// handleHSBCPayment handles the payment request
func handleHSBCPayment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	userID := r.Header.Get("user_id")
	userSecret := r.Header.Get("user_secret")

	if userID == "" || userSecret == "" {
		http.Error(w, "Missing required headers: user_id or user_secret", http.StatusBadRequest)
		return
	}

	var paymentRequest PaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&paymentRequest); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Generate a UUID for external_id
	externalID := uuid.New().String()

	// Simulate generating a URL for the payment
	paymentURL := fmt.Sprintf("http://localhost:8081/hsbc/callback?external_id=%s", externalID)

	// Respond with the payment URL and external ID
	response := PaymentResponse{
		URL:        paymentURL,
		ExternalID: externalID,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	log.Printf("Payment request received: User ID: %s, Amount: %.2f, PaymentType: %s, Currency: %s, Country: %s", userID, paymentRequest.Amount, paymentRequest.PaymentType, paymentRequest.Currency, paymentRequest.Country)
	log.Printf("Generated payment URL: %s and External ID: %s", paymentURL, externalID)
}

// handleHSBCCallback simulates the callback to the payment service's callback URL
func handleHSBCCallback(w http.ResponseWriter, r *http.Request) {
	externalID := r.URL.Query().Get("external_id")
	if externalID == "" {
		http.Error(w, "Missing required query parameter: external_id", http.StatusBadRequest)
		return
	}

	// Simulate the callback to the payment service
	callbackURL := fmt.Sprintf("http://localhost:8080/payment/callbacks/success?id=%s", externalID)

	//Redirect
	http.Redirect(w, r, callbackURL, http.StatusFound)
}
