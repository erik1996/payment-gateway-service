package middleware

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RequestIDMiddleware generates a UUID request ID and adds it to the request context and headers
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Generate a new UUID for the request ID
		requestID := uuid.New().String()

		// Set the request ID in the request context
		c.Set("RequestID", requestID)

		// Add the request ID to the response headers
		c.Writer.Header().Set("X-Request-ID", requestID)

		// Log the request ID
		log.Printf("Generated Request ID: %s", requestID)

		// Proceed to the next middleware/handler
		c.Next()
	}
}
