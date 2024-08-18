package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
)

// PaymentRequest represents the structure of the payment request
type PaymentRequest struct {
	XMLName     xml.Name `xml:"PaymentRequest"`
	Amount      float64  `xml:"Amount"`
	PaymentType string   `xml:"PaymentType"`
	Currency    string   `xml:"Currency"`
	Country     string   `xml:"Country"`
}

// PaymentResponse represents the structure of the payment response
type PaymentResponse struct {
	XMLName    xml.Name `xml:"PaymentResponse"`
	URL        string   `xml:"URL"`
	ExternalID string   `xml:"ExternalID"`
}

// CallbackRequest represents the structure of the callback request
type CallbackRequest struct {
	XMLName    xml.Name `xml:"CallbackRequest"`
	ExternalID string   `xml:"ExternalID"`
	Status     string   `xml:"Status"`
}

func main() {
	http.HandleFunc("/adcb/payment", handleADCMPayment)
	http.HandleFunc("/adcb/callback", handleADCBCallback)
	log.Println("ADCB Mock Service running on port 8082")
	log.Fatal(http.ListenAndServe(":8082", nil))
}

// handleADCMPayment handles the payment request for ADCB
func handleADCMPayment(w http.ResponseWriter, r *http.Request) {
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
	if err := xml.NewDecoder(r.Body).Decode(&paymentRequest); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Generate a UUID for external_id
	externalID := uuid.New().String()

	// Simulate generating a URL for the payment
	paymentURL := fmt.Sprintf("http://localhost:8082/adcb/callback?external_id=%s", externalID)

	// Respond with the payment URL and external ID
	response := PaymentResponse{
		URL:        paymentURL,
		ExternalID: externalID,
	}

	responseXML, err := xml.MarshalIndent(response, "", "  ")
	if err != nil {
		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/xml")
	w.Write(responseXML)

	log.Printf("Payment request received: User ID: %s, Amount: %.2f, PaymentType: %s, Currency: %s, Country: %s", userID, paymentRequest.Amount, paymentRequest.PaymentType, paymentRequest.Currency, paymentRequest.Country)
	log.Printf("Generated payment URL: %s and External ID: %s", paymentURL, externalID)
}

// handleADCBCallback simulates the callback to the payment service's callback URL
func handleADCBCallback(w http.ResponseWriter, r *http.Request) {
	externalID := r.URL.Query().Get("external_id")

	if externalID == "" {
		http.Error(w, "Missing required query parameter: external_id", http.StatusBadRequest)
		return
	}

	callbackURL := fmt.Sprintf("http://localhost:8080/payment/callbacks/success?id=%s", externalID)

	//Redirect
	http.Redirect(w, r, callbackURL, http.StatusFound)
}
