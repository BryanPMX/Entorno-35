package ports

import (
	"context"
	"time"

	"github.com/entorno35/backend/internal/domain"
	"github.com/google/uuid"
)

// PaymentRepository defines the interface for payment data operations
type PaymentRepository interface {
	// CreatePayment creates a new payment record
	CreatePayment(ctx context.Context, payment *domain.Payment) error

	// UpdatePaymentStatus updates the status of an existing payment
	UpdatePaymentStatus(ctx context.Context, paymentID uuid.UUID, status domain.PaymentStatus) error

	// GetPaymentByStripeID retrieves a payment by Stripe payment ID
	GetPaymentByStripeID(ctx context.Context, stripePaymentID string) (*domain.Payment, error)

	// GetPaymentsByCompany retrieves all payments for a company
	GetPaymentsByCompany(ctx context.Context, companyID uuid.UUID) ([]domain.Payment, error)

	// GetSubscriptionByCompany retrieves the active subscription for a company
	GetSubscriptionByCompany(ctx context.Context, companyID uuid.UUID) (*domain.Subscription, error)

	// CreateSubscription creates a new subscription record
	CreateSubscription(ctx context.Context, subscription *domain.Subscription) error

	// UpdateSubscription updates an existing subscription
	UpdateSubscription(ctx context.Context, subscription *domain.Subscription) error

	// CancelSubscription marks a subscription for cancellation
	CancelSubscription(ctx context.Context, subscriptionID uuid.UUID, canceledAt *time.Time) error

	// GetSubscriptionByStripeID retrieves a subscription by Stripe subscription ID
	GetSubscriptionByStripeID(ctx context.Context, stripeSubscriptionID string) (*domain.Subscription, error)
}

// PaymentService defines the interface for payment business logic
type PaymentService interface {
	// CreateSubscription creates a new subscription for a company
	CreateSubscription(ctx context.Context, companyID uuid.UUID, priceID string) (*domain.Subscription, error)

	// CancelSubscription cancels an existing subscription
	CancelSubscription(ctx context.Context, companyID uuid.UUID) error

	// GetSubscription retrieves the current subscription for a company
	GetSubscription(ctx context.Context, companyID uuid.UUID) (*domain.Subscription, error)

	// ProcessWebhook processes incoming Stripe webhooks
	ProcessWebhook(ctx context.Context, payload []byte, signature string) error

	// IsCompanySubscribed checks if a company has an active subscription
	IsCompanySubscribed(ctx context.Context, companyID uuid.UUID) (bool, error)
}