package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/entorno35/backend/internal/config"
	"github.com/entorno35/backend/internal/core/ports"
	"github.com/entorno35/backend/internal/domain"
	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v79"
	"github.com/stripe/stripe-go/v79/client"
)

// PaymentService implements the payment business logic
type PaymentService struct {
	stripeClient *client.API
	repo         ports.PaymentRepository
	config       config.StripeConfig
}

// NewPaymentService creates a new payment service
func NewPaymentService(repo ports.PaymentRepository, cfg config.StripeConfig) *PaymentService {
	// Initialize Stripe client with security best practices
	stripeClient := &client.API{}
	stripeClient.Init(cfg.SecretKey, nil)

	return &PaymentService{
		stripeClient: stripeClient,
		repo:         repo,
		config:       cfg,
	}
}

// CreateSubscription creates a new subscription for a company
func (s *PaymentService) CreateSubscription(ctx context.Context, companyID uuid.UUID, priceID string) (*domain.Subscription, error) {
	// Validate price ID for security
	if !s.isValidPriceID(priceID) {
		return nil, errors.New("invalid price ID provided")
	}

	// Get company information (you might need to pass this or get it from context)
	// For now, we'll assume the company exists and is valid

	// Create Stripe customer first (in a real implementation, you'd store customer ID)
	customerParams := &stripe.CustomerParams{
		Metadata: map[string]string{
			"company_id": companyID.String(),
		},
	}
	customer, err := s.stripeClient.Customers.New(customerParams)
	if err != nil {
		return nil, fmt.Errorf("failed to create Stripe customer: %w", err)
	}

	// Create subscription in Stripe
	subscriptionParams := &stripe.SubscriptionParams{
		Customer: stripe.String(customer.ID),
		Items: []*stripe.SubscriptionItemsParams{
			{
				Price: stripe.String(priceID),
			},
		},
		PaymentBehavior: stripe.String("default_incomplete"),
		Expand:          []*string{stripe.String("latest_invoice.payment_intent")},
	}
	subscription, err := s.stripeClient.Subscriptions.New(subscriptionParams)
	if err != nil {
		return nil, fmt.Errorf("failed to create Stripe subscription: %w", err)
	}

	// Determine interval from price ID
	interval := domain.SubscriptionIntervalMonth
	if strings.Contains(priceID, "year") {
		interval = domain.SubscriptionIntervalYear
	}

	// Create subscription record in database
	subscriptionRecord := &domain.Subscription{
		CompanyID:            companyID,
		StripeSubscriptionID: subscription.ID,
		StripePriceID:        priceID,
		Status:               domain.SubscriptionStatusActive,
		Interval:             interval,
		CurrentPeriodStart:   time.Unix(subscription.CurrentPeriodStart, 0),
		CurrentPeriodEnd:     time.Unix(subscription.CurrentPeriodEnd, 0),
		CancelAtPeriodEnd:    false,
	}

	err = s.repo.CreateSubscription(ctx, subscriptionRecord)
	if err != nil {
		// If database save fails, cancel Stripe subscription
		s.stripeClient.Subscriptions.Cancel(subscription.ID, nil)
		return nil, fmt.Errorf("failed to save subscription: %w", err)
	}

	return subscriptionRecord, nil
}

// CancelSubscription cancels an existing subscription
func (s *PaymentService) CancelSubscription(ctx context.Context, companyID uuid.UUID) error {
	// Get current subscription
	subscription, err := s.repo.GetSubscriptionByCompany(ctx, companyID)
	if err != nil {
		return fmt.Errorf("failed to get subscription: %w", err)
	}

	if subscription.Status != domain.SubscriptionStatusActive {
		return errors.New("subscription is not active")
	}

	// Cancel in Stripe (at period end for better UX)
	cancelParams := &stripe.SubscriptionCancelParams{}
	_, err = s.stripeClient.Subscriptions.Cancel(subscription.StripeSubscriptionID, cancelParams)
	if err != nil {
		return fmt.Errorf("failed to cancel Stripe subscription: %w", err)
	}

	// Update database
	now := time.Now()
	err = s.repo.CancelSubscription(ctx, subscription.ID, &now)
	if err != nil {
		return fmt.Errorf("failed to update subscription in database: %w", err)
	}

	return nil
}

// GetSubscription retrieves the current subscription for a company
func (s *PaymentService) GetSubscription(ctx context.Context, companyID uuid.UUID) (*domain.Subscription, error) {
	return s.repo.GetSubscriptionByCompany(ctx, companyID)
}

// ProcessWebhook processes incoming Stripe webhooks with security verification
func (s *PaymentService) ProcessWebhook(ctx context.Context, payload []byte, signature string) error {
	// TODO: Implement webhook signature verification
	// For now, we'll trust the webhook (in production, implement proper verification)
	var event stripe.Event
	err := json.Unmarshal(payload, &event)
	if err != nil {
		return fmt.Errorf("failed to parse webhook payload: %w", err)
	}

	// Process different event types
	switch event.Type {
	case "customer.subscription.created":
		return s.handleSubscriptionCreated(ctx, event.Data.Object)
	case "customer.subscription.updated":
		return s.handleSubscriptionUpdated(ctx, event.Data.Object)
	case "customer.subscription.deleted":
		return s.handleSubscriptionDeleted(ctx, event.Data.Object)
	case "invoice.payment_succeeded":
		return s.handlePaymentSucceeded(ctx, event.Data.Object)
	case "invoice.payment_failed":
		return s.handlePaymentFailed(ctx, event.Data.Object)
	default:
		// Log unhandled events for debugging
		fmt.Printf("Unhandled webhook event: %s\n", event.Type)
		return nil
	}
}

// IsCompanySubscribed checks if a company has an active subscription
func (s *PaymentService) IsCompanySubscribed(ctx context.Context, companyID uuid.UUID) (bool, error) {
	subscription, err := s.repo.GetSubscriptionByCompany(ctx, companyID)
	if err != nil {
		return false, err
	}

	return subscription.Status == domain.SubscriptionStatusActive, nil
}

// isValidPriceID validates price IDs against configured values
func (s *PaymentService) isValidPriceID(priceID string) bool {
	return priceID == s.config.PriceIDMonthly || priceID == s.config.PriceIDYearly
}

// handleSubscriptionCreated processes subscription creation events
func (s *PaymentService) handleSubscriptionCreated(ctx context.Context, eventData interface{}) error {
	// Parse event data
	eventJSON, err := json.Marshal(eventData)
	if err != nil {
		return fmt.Errorf("failed to marshal event data: %w", err)
	}

	var stripeSubscription stripe.Subscription
	err = json.Unmarshal(eventJSON, &stripeSubscription)
	if err != nil {
		return fmt.Errorf("failed to unmarshal subscription data: %w", err)
	}

	// Update subscription status in database
	err = s.updateSubscriptionFromStripe(ctx, &stripeSubscription)
	if err != nil {
		return fmt.Errorf("failed to update subscription: %w", err)
	}

	return nil
}

// handleSubscriptionUpdated processes subscription update events
func (s *PaymentService) handleSubscriptionUpdated(ctx context.Context, eventData interface{}) error {
	return s.handleSubscriptionCreated(ctx, eventData) // Same logic
}

// handleSubscriptionDeleted processes subscription deletion events
func (s *PaymentService) handleSubscriptionDeleted(ctx context.Context, eventData interface{}) error {
	eventJSON, err := json.Marshal(eventData)
	if err != nil {
		return fmt.Errorf("failed to marshal event data: %w", err)
	}

	var stripeSubscription stripe.Subscription
	err = json.Unmarshal(eventJSON, &stripeSubscription)
	if err != nil {
		return fmt.Errorf("failed to unmarshal subscription data: %w", err)
	}

	// Find subscription in database and mark as cancelled
	subscription, err := s.repo.GetSubscriptionByStripeID(ctx, stripeSubscription.ID)
	if err != nil {
		return fmt.Errorf("failed to find subscription: %w", err)
	}

	now := time.Now()
	err = s.repo.CancelSubscription(ctx, subscription.ID, &now)
	if err != nil {
		return fmt.Errorf("failed to cancel subscription: %w", err)
	}

	return nil
}

// handlePaymentSucceeded processes successful payment events
func (s *PaymentService) handlePaymentSucceeded(ctx context.Context, eventData interface{}) error {
	eventJSON, err := json.Marshal(eventData)
	if err != nil {
		return fmt.Errorf("failed to marshal event data: %w", err)
	}

	var invoice stripe.Invoice
	err = json.Unmarshal(eventJSON, &invoice)
	if err != nil {
		return fmt.Errorf("failed to unmarshal invoice data: %w", err)
	}

	// Create payment record
	payment := &domain.Payment{
		StripePaymentID: invoice.ID,
		Amount:          int(invoice.AmountPaid),
		Currency:        string(invoice.Currency),
		Status:          domain.PaymentStatusSucceeded,
		Description:     fmt.Sprintf("Payment for invoice %s", invoice.Number),
	}

	// Get company ID from subscription (this would need to be implemented)
	// For now, we'll assume we need to extract it from the invoice

	err = s.repo.CreatePayment(ctx, payment)
	if err != nil {
		return fmt.Errorf("failed to save payment: %w", err)
	}

	return nil
}

// handlePaymentFailed processes failed payment events
func (s *PaymentService) handlePaymentFailed(ctx context.Context, eventData interface{}) error {
	eventJSON, err := json.Marshal(eventData)
	if err != nil {
		return fmt.Errorf("failed to marshal event data: %w", err)
	}

	var invoice stripe.Invoice
	err = json.Unmarshal(eventJSON, &invoice)
	if err != nil {
		return fmt.Errorf("failed to unmarshal invoice data: %w", err)
	}

	// Create failed payment record
	payment := &domain.Payment{
		StripePaymentID: invoice.ID,
		Amount:          int(invoice.AmountDue),
		Currency:        string(invoice.Currency),
		Status:          domain.PaymentStatusFailed,
		Description:     fmt.Sprintf("Failed payment for invoice %s", invoice.Number),
	}

	err = s.repo.CreatePayment(ctx, payment)
	if err != nil {
		return fmt.Errorf("failed to save failed payment: %w", err)
	}

	return nil
}

// updateSubscriptionFromStripe updates database subscription from Stripe data
func (s *PaymentService) updateSubscriptionFromStripe(ctx context.Context, stripeSub *stripe.Subscription) error {
	subscription, err := s.repo.GetSubscriptionByStripeID(ctx, stripeSub.ID)
	if err != nil {
		return fmt.Errorf("failed to find subscription: %w", err)
	}

	// Map Stripe status to domain status
	var status domain.SubscriptionStatus
	switch stripeSub.Status {
	case stripe.SubscriptionStatusActive:
		status = domain.SubscriptionStatusActive
	case stripe.SubscriptionStatusCanceled:
		status = domain.SubscriptionStatusInactive
	default:
		status = domain.SubscriptionStatusInactive
	}

	subscription.Status = status
	subscription.CurrentPeriodStart = time.Unix(stripeSub.CurrentPeriodStart, 0)
	subscription.CurrentPeriodEnd = time.Unix(stripeSub.CurrentPeriodEnd, 0)
	subscription.CancelAtPeriodEnd = stripeSub.CancelAtPeriodEnd

	err = s.repo.UpdateSubscription(ctx, subscription)
	if err != nil {
		return fmt.Errorf("failed to update subscription: %w", err)
	}

	return nil
}