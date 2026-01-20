package middleware

import (
	"net/http"

	"github.com/entorno35/backend/internal/core/ports"
	"github.com/gin-gonic/gin"
)

// SubscriptionMiddleware creates middleware that checks for active subscription
func SubscriptionMiddleware(paymentService ports.PaymentService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get auth context (must be set by AuthMiddleware first)
		authCtx, exists := GetAuthContext(c)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
			c.Abort()
			return
		}

		// Check if company has active subscription
		isSubscribed, err := paymentService.IsCompanySubscribed(c.Request.Context(), authCtx.CompanyID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check subscription status"})
			c.Abort()
			return
		}

		if !isSubscribed {
			c.JSON(http.StatusPaymentRequired, gin.H{
				"error": "active subscription required",
				"message": "Please subscribe to access this feature",
				"action": "subscribe",
			})
			c.Abort()
			return
		}

		// Set subscription status in context for later use
		c.Set("is_subscribed", true)
		c.Next()
	}
}

// OptionalSubscriptionMiddleware allows access but marks subscription status
func OptionalSubscriptionMiddleware(paymentService ports.PaymentService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get auth context
		authCtx, exists := GetAuthContext(c)
		if !exists {
			// No auth context, allow access but mark as not subscribed
			c.Set("is_subscribed", false)
			c.Next()
			return
		}

		// Check subscription status
		isSubscribed, err := paymentService.IsCompanySubscribed(c.Request.Context(), authCtx.CompanyID)
		if err != nil {
			// On error, default to not subscribed for security
			c.Set("is_subscribed", false)
			c.Next()
			return
		}

		c.Set("is_subscribed", isSubscribed)
		c.Next()
	}
}

// GetSubscriptionStatus extracts subscription status from context
func GetSubscriptionStatus(c *gin.Context) bool {
	status, exists := c.Get("is_subscribed")
	if !exists {
		return false
	}

	isSubscribed, ok := status.(bool)
	if !ok {
		return false
	}

	return isSubscribed
}