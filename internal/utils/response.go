package utils

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

// APIResponse is the structure for all API responses
type APIResponse struct {
	Status  string              `json:"status"`
	Message string              `json:"message,omitempty"`
	Errors  map[string][]string `json:"errors,omitempty"`
	Data    interface{}         `json:"data,omitempty"`
}

// ErrorResponse sends a JSON error response with a specific status code and aborts the request
func ErrorResponse(c *gin.Context, statusCode int, message string, errors map[string][]string) {
	requestID := c.GetString("RequestID")

	// Prepare error details as JSON string for logging
	errorDetails, _ := json.Marshal(errors)

	logMessage := fmt.Sprintf("Request ID: %s - Sending error response: %s, Errors: %s", requestID, message, string(errorDetails))
	log.Println(logMessage)

	c.JSON(statusCode, APIResponse{
		Status:  "error",
		Message: message,
		Errors:  errors,
	})
	c.Abort()
}

// SuccessResponse sends a JSON success response with a specific status code
func SuccessResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	requestID := c.GetString("RequestID")

	// Prepare data details as JSON string for logging
	dataDetails, _ := json.Marshal(data)

	logMessage := fmt.Sprintf("Request ID: %s - Sending success response: %s, Data: %s", requestID, message, string(dataDetails))
	log.Println(logMessage)

	c.JSON(statusCode, APIResponse{
		Status:  "success",
		Message: message,
		Data:    data,
	})
}
