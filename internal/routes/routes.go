package routes

import (
	"payment-gateway-service/internal/middleware"
	"payment-gateway-service/internal/payment"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

func RegisterRoutes(router *gin.Engine, db *gorm.DB) {

	// Initialize handlers with the correct package paths
	paymentHandler := payment.NewPaymentHandler(db)

	// Register payment routes with validation middleware
	paymentRoutes := router.Group("/payment")
	{
		paymentRoutes.POST("/deposit", middleware.AuthMiddleware(), middleware.ValidationMiddleware(&payment.PaymentRequest{}), paymentHandler.Deposit)
		paymentRoutes.POST("/withdrawal", middleware.AuthMiddleware(), middleware.ValidationMiddleware(&payment.PaymentRequest{}), paymentHandler.Withdrawal)
		paymentRoutes.GET("/callbacks/success", paymentHandler.HandleSuccessCallback)
		paymentRoutes.GET("/callbacks/success/:external_id", paymentHandler.HandleSuccessCallback)
		paymentRoutes.GET("/callbacks/failed", paymentHandler.HandleFailedCallback)
		paymentRoutes.GET("/callbacks/failed/:external_id", paymentHandler.HandleFailedCallback)

		paymentRoutes.GET("/", paymentHandler.PaymentStatus)
	}

	// Swagger Route
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
