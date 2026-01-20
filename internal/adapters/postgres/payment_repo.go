package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/entorno35/backend/internal/core/ports"
	"github.com/entorno35/backend/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrPaymentNotFound     = errors.New("payment not found")
	ErrSubscriptionNotFound = errors.New("subscription not found")
)

// PaymentRepository implements the PaymentRepository interface
type PaymentRepository struct {
	db *gorm.DB
}

// NewPaymentRepository creates a new payment repository
func NewPaymentRepository(db *gorm.DB) ports.PaymentRepository {
	return &PaymentRepository{
		db: db,
	}
}

// CreatePayment creates a new payment record
func (r *PaymentRepository) CreatePayment(ctx context.Context, payment *domain.Payment) error {
	return r.db.WithContext(ctx).Create(payment).Error
}

// UpdatePaymentStatus updates the status of an existing payment
func (r *PaymentRepository) UpdatePaymentStatus(ctx context.Context, paymentID uuid.UUID, status domain.PaymentStatus) error {
	result := r.db.WithContext(ctx).Model(&domain.Payment{}).
		Where("id = ?", paymentID).
		Update("status", status)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrPaymentNotFound
	}

	return nil
}

// GetPaymentByStripeID retrieves a payment by Stripe payment ID
func (r *PaymentRepository) GetPaymentByStripeID(ctx context.Context, stripePaymentID string) (*domain.Payment, error) {
	var payment domain.Payment
	err := r.db.WithContext(ctx).
		Preload("Company").
		Where("stripe_payment_id = ?", stripePaymentID).
		First(&payment).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPaymentNotFound
		}
		return nil, err
	}

	return &payment, nil
}

// GetPaymentsByCompany retrieves all payments for a company
func (r *PaymentRepository) GetPaymentsByCompany(ctx context.Context, companyID uuid.UUID) ([]domain.Payment, error) {
	var payments []domain.Payment
	err := r.db.WithContext(ctx).
		Where("company_id = ?", companyID).
		Order("created_at DESC").
		Find(&payments).Error

	return payments, err
}

// GetSubscriptionByCompany retrieves the active subscription for a company
func (r *PaymentRepository) GetSubscriptionByCompany(ctx context.Context, companyID uuid.UUID) (*domain.Subscription, error) {
	var subscription domain.Subscription
	err := r.db.WithContext(ctx).
		Preload("Company").
		Where("company_id = ?", companyID).
		First(&subscription).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrSubscriptionNotFound
		}
		return nil, err
	}

	return &subscription, nil
}

// CreateSubscription creates a new subscription record
func (r *PaymentRepository) CreateSubscription(ctx context.Context, subscription *domain.Subscription) error {
	return r.db.WithContext(ctx).Create(subscription).Error
}

// UpdateSubscription updates an existing subscription
func (r *PaymentRepository) UpdateSubscription(ctx context.Context, subscription *domain.Subscription) error {
	return r.db.WithContext(ctx).
		Where("id = ?", subscription.ID).
		Updates(subscription).Error
}

// CancelSubscription marks a subscription for cancellation
func (r *PaymentRepository) CancelSubscription(ctx context.Context, subscriptionID uuid.UUID, canceledAt *time.Time) error {
	updates := map[string]interface{}{
		"status": domain.SubscriptionStatusInactive,
		"canceled_at": canceledAt,
	}

	result := r.db.WithContext(ctx).Model(&domain.Subscription{}).
		Where("id = ?", subscriptionID).
		Updates(updates)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrSubscriptionNotFound
	}

	return nil
}

// GetSubscriptionByStripeID retrieves a subscription by Stripe subscription ID
func (r *PaymentRepository) GetSubscriptionByStripeID(ctx context.Context, stripeSubscriptionID string) (*domain.Subscription, error) {
	var subscription domain.Subscription
	err := r.db.WithContext(ctx).
		Preload("Company").
		Where("stripe_subscription_id = ?", stripeSubscriptionID).
		First(&subscription).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrSubscriptionNotFound
		}
		return nil, err
	}

	return &subscription, nil
}