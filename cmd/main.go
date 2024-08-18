package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"payment-gateway-service/config"
	_ "payment-gateway-service/docs"
	"payment-gateway-service/internal/database"
	"payment-gateway-service/internal/middleware"
	"payment-gateway-service/internal/routes"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Connect to the database with GORM
	db, err := database.ConnectPostgres(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Ensure the database connection is closed when the program exits
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get database object: %v", err)
	}
	defer sqlDB.Close()

	// Initialize the Gin engine
	router := gin.Default()

	// Apply the RequestIDMiddleware globally
	router.Use(middleware.RequestIDMiddleware())

	// Register routes with the gorm.DB instance
	routes.RegisterRoutes(router, db)

	// Construct the address with port
	address := ":" + cfg.PORT

	// Print the address to the logs
	log.Printf("Starting server on %s", address)

	// Server settings
	srv := &http.Server{
		Addr:              address,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	// Start the server in a goroutine so it doesn’t block
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
