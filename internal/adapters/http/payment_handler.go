package http

import (
	"io"
	"net/http"

	"github.com/entorno35/backend/internal/core/ports"
	"github.com/entorno35/backend/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateSubscriptionRequest represents the request to create a subscription
type CreateSubscriptionRequest struct {
	PriceID string `json:"price_id" binding:"required" example:"price_monthly"`
}

// SubscriptionResponse represents a subscription in API responses
type SubscriptionResponse struct {
	ID                   string    `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	CompanyID            string    `json:"company_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	StripeSubscriptionID string    `json:"stripe_subscription_id" example:"sub_1234567890"`
	Status               string    `json:"status" example:"active"`
	Interval             string    `json:"interval" example:"month"`
	CurrentPeriodStart   string    `json:"current_period_start" example:"2024-01-01T00:00:00Z"`
	CurrentPeriodEnd     string    `json:"current_period_end" example:"2024-02-01T00:00:00Z"`
	CancelAtPeriodEnd    bool      `json:"cancel_at_period_end" example:false`
	CreatedAt            string    `json:"created_at" example:"2024-01-01T00:00:00Z"`
}

// PaymentHandler handles payment-related HTTP requests
type PaymentHandler struct {
	paymentService ports.PaymentService
}

// NewPaymentHandler creates a new payment handler
func NewPaymentHandler(paymentService ports.PaymentService) *PaymentHandler {
	return &PaymentHandler{
		paymentService: paymentService,
	}
}

// CreateSubscription creates a new subscription for the authenticated company
// POST /api/v1/payments/subscription
func (h *PaymentHandler) CreateSubscription(c *gin.Context) {
	authCtx, exists := middleware.GetAuthContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	var req CreateSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate that user is a company admin
	if authCtx.StaffID != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "only company administrators can manage subscriptions"})
		return
	}

	subscription, err := h.paymentService.CreateSubscription(c.Request.Context(), authCtx.CompanyID, req.PriceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := SubscriptionResponse{
		ID:                   subscription.ID.String(),
		CompanyID:            subscription.CompanyID.String(),
		StripeSubscriptionID: subscription.StripeSubscriptionID,
		Status:               string(subscription.Status),
		Interval:             string(subscription.Interval),
		CurrentPeriodStart:   subscription.CurrentPeriodStart.Format("2006-01-02T15:04:05Z"),
		CurrentPeriodEnd:     subscription.CurrentPeriodEnd.Format("2006-01-02T15:04:05Z"),
		CancelAtPeriodEnd:    subscription.CancelAtPeriodEnd,
		CreatedAt:            subscription.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}

	c.JSON(http.StatusCreated, response)
}

// GetSubscription retrieves the current subscription for the authenticated company
// GET /api/v1/payments/subscription
func (h *PaymentHandler) GetSubscription(c *gin.Context) {
	authCtx, exists := middleware.GetAuthContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	subscription, err := h.paymentService.GetSubscription(c.Request.Context(), authCtx.CompanyID)
	if err != nil {
		if err.Error() == "subscription not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "no active subscription found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := SubscriptionResponse{
		ID:                   subscription.ID.String(),
		CompanyID:            subscription.CompanyID.String(),
		StripeSubscriptionID: subscription.StripeSubscriptionID,
		Status:               string(subscription.Status),
		Interval:             string(subscription.Interval),
		CurrentPeriodStart:   subscription.CurrentPeriodStart.Format("2006-01-02T15:04:05Z"),
		CurrentPeriodEnd:     subscription.CurrentPeriodEnd.Format("2006-01-02T15:04:05Z"),
		CancelAtPeriodEnd:    subscription.CancelAtPeriodEnd,
		CreatedAt:            subscription.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}

	c.JSON(http.StatusOK, response)
}

// CancelSubscription cancels the current subscription for the authenticated company
// DELETE /api/v1/payments/subscription
func (h *PaymentHandler) CancelSubscription(c *gin.Context) {
	authCtx, exists := middleware.GetAuthContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	// Validate that user is a company admin
	if authCtx.StaffID != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "only company administrators can manage subscriptions"})
		return
	}

	err := h.paymentService.CancelSubscription(c.Request.Context(), authCtx.CompanyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "subscription cancelled successfully"})
}

// HandleWebhook processes Stripe webhooks
// POST /api/v1/payments/webhook
func (h *PaymentHandler) HandleWebhook(c *gin.Context) {
	// Get the raw body for signature verification
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read request body"})
		return
	}

	// Get Stripe signature from headers
	signature := c.GetHeader("Stripe-Signature")
	if signature == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing Stripe signature"})
		return
	}

	// Process webhook
	err = h.paymentService.ProcessWebhook(c.Request.Context(), body, signature)
	if err != nil {
		// Log error but return 200 to prevent Stripe from retrying
		// (since webhook signature verification failed or processing failed)
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "webhook processed successfully"})
}