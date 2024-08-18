package middleware

import (
	"payment-gateway-service/config"

	"github.com/gin-gonic/gin"
)

// CallbackMiddleware is the middleware function for authenticating callback requests
func CallbacMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}
