package utils

import (
	"context"
	"log"
)

// LogWithRequestID logs a message with the RequestID extracted from the context.
func LogWithRequestID(ctx context.Context, message string) {
	requestID, _ := ctx.Value("RequestID").(string)
	if requestID != "" {
		log.Printf("RequestID: %s - %s", requestID, message)
	} else {
		log.Printf("%s", message)
	}
}
