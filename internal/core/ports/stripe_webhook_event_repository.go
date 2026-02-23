package ports

// StripeWebhookEventRepository persists webhook processing state for idempotency.
type StripeWebhookEventRepository interface {
	TryBegin(eventID, eventType string) (bool, error)
	MarkProcessed(eventID string) error
	MarkFailed(eventID, failure string) error
}
