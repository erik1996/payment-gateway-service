package middleware

import (
	"net/http"
	"payment-gateway-service/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// ValidationMiddleware validates the incoming request body against the provided struct
func ValidationMiddleware(obj interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		// Create a new instance of the provided struct type
		objInstance := obj

		// Bind the incoming JSON to the struct
		if err := c.ShouldBindJSON(&objInstance); err != nil {
			utils.LogWithRequestID(ctx, "Validation error occurred")

			// Handle validation errors
			if validationErrs, ok := err.(validator.ValidationErrors); ok {
				errors := make(map[string][]string)
				for _, validationErr := range validationErrs {
					field := validationErr.Field()
					tag := validationErr.Tag()

					var errorMessage string
					switch tag {
					case "required":
						errorMessage = "is required"
					case "gt":
						errorMessage = "must be greater than " + validationErr.Param()
					case "len":
						errorMessage = "must be exactly " + validationErr.Param() + " characters"
					case "oneof":
						errorMessage = "must be one of " + validationErr.Param()
					default:
						errorMessage = "is invalid"
					}

					errors[field] = append(errors[field], errorMessage)
				}

				utils.ErrorResponse(c, http.StatusBadRequest, "Validation failed", errors)
				return
			}

			// Handle other binding errors, including missing required fields
			errors := make(map[string][]string)
			errors["validation"] = append(errors["validation"], err.Error())

			utils.ErrorResponse(c, http.StatusBadRequest, "Validation failed", errors)
			return
		}

		// Log successful validation
		utils.LogWithRequestID(ctx, "Validation succeeded")

		// If validation passes, store the validated struct in the context
		c.Set("validatedBody", objInstance)
		c.Next()
	}
}
