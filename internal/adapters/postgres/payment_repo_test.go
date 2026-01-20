package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/entorno35/backend/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPaymentRepository(t *testing.T) {
	// Setup test database
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewPaymentRepository(db)

	ctx := context.Background()
	companyID := uuid.New()

	t.Run("CreatePayment", func(t *testing.T) {
		payment := &domain.Payment{
			CompanyID:       companyID,
			StripePaymentID: "pi_test_123",
			Amount:          50000, // $500.00 in cents
			Currency:        "mxn",
			Status:          domain.PaymentStatusSucceeded,
			Description:     "Test payment",
		}

		err := repo.CreatePayment(ctx, payment)
		require.NoError(t, err)
		assert.NotEqual(t, uuid.Nil, payment.ID)
	})

	t.Run("GetPaymentByStripeID", func(t *testing.T) {
		payment, err := repo.GetPaymentByStripeID(ctx, "pi_test_123")
		require.NoError(t, err)
		assert.Equal(t, companyID, payment.CompanyID)
		assert.Equal(t, int64(50000), payment.Amount)
		assert.Equal(t, domain.PaymentStatusSucceeded, payment.Status)
	})

	t.Run("GetPaymentsByCompany", func(t *testing.T) {
		payments, err := repo.GetPaymentsByCompany(ctx, companyID)
		require.NoError(t, err)
		assert.Len(t, payments, 1)
		assert.Equal(t, "pi_test_123", payments[0].StripePaymentID)
	})

	t.Run("CreateSubscription", func(t *testing.T) {
		now := time.Now()
		subscription := &domain.Subscription{
			CompanyID:            companyID,
			StripeSubscriptionID: "sub_test_123",
			StripePriceID:        "price_monthly",
			Status:               domain.SubscriptionStatusActive,
			Interval:             domain.SubscriptionIntervalMonth,
			CurrentPeriodStart:   now,
			CurrentPeriodEnd:     now.AddDate(0, 1, 0), // 1 month later
			CancelAtPeriodEnd:    false,
		}

		err := repo.CreateSubscription(ctx, subscription)
		require.NoError(t, err)
		assert.NotEqual(t, uuid.Nil, subscription.ID)
	})

	t.Run("GetSubscriptionByCompany", func(t *testing.T) {
		subscription, err := repo.GetSubscriptionByCompany(ctx, companyID)
		require.NoError(t, err)
		assert.Equal(t, companyID, subscription.CompanyID)
		assert.Equal(t, "sub_test_123", subscription.StripeSubscriptionID)
		assert.Equal(t, domain.SubscriptionStatusActive, subscription.Status)
	})

	t.Run("GetSubscriptionByStripeID", func(t *testing.T) {
		subscription, err := repo.GetSubscriptionByStripeID(ctx, "sub_test_123")
		require.NoError(t, err)
		assert.Equal(t, companyID, subscription.CompanyID)
		assert.Equal(t, domain.SubscriptionIntervalMonth, subscription.Interval)
	})

	t.Run("CancelSubscription", func(t *testing.T) {
		subscription, err := repo.GetSubscriptionByCompany(ctx, companyID)
		require.NoError(t, err)

		cancelTime := time.Now()
		err = repo.CancelSubscription(ctx, subscription.ID, &cancelTime)
		require.NoError(t, err)

		// Verify cancellation
		updatedSubscription, err := repo.GetSubscriptionByCompany(ctx, companyID)
		require.NoError(t, err)
		assert.Equal(t, domain.SubscriptionStatusInactive, updatedSubscription.Status)
		assert.Equal(t, cancelTime.Unix(), updatedSubscription.CanceledAt.Unix())
	})

	t.Run("GetNonExistentPayment", func(t *testing.T) {
		_, err := repo.GetPaymentByStripeID(ctx, "nonexistent")
		assert.Equal(t, ErrPaymentNotFound, err)
	})

	t.Run("GetNonExistentSubscription", func(t *testing.T) {
		nonExistentCompanyID := uuid.New()
		_, err := repo.GetSubscriptionByCompany(ctx, nonExistentCompanyID)
		assert.Equal(t, ErrSubscriptionNotFound, err)
	})
}