package middleware

import (
	"net/http"
	config "payment-gateway-service/config"
	"payment-gateway-service/internal/utils"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware is the middleware function for authentication
func AuthMiddleware() gin.HandlerFunc {
	cfg := config.LoadConfig() // Load configuration here

	return func(c *gin.Context) {
		// Check for the X-AUTH-TOKEN header
		authTokenHeader := c.GetHeader("X-AUTH-TOKEN")

		// Validate the header value against the configured token
		if authTokenHeader != cfg.AuthToken {
			// If invalid, respond with 401 Unauthorized
			utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", nil)
			return
		}

		// If valid, proceed with the request
		c.Next()
	}
}
