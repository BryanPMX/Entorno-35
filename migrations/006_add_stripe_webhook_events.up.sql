-- Persist Stripe webhook processing state for idempotency and retries
CREATE TABLE stripe_webhook_events (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    stripe_event_id VARCHAR(255) NOT NULL,
    event_type VARCHAR(100) NOT NULL,
    status VARCHAR(20) NOT NULL,
    attempt_count INTEGER NOT NULL DEFAULT 1,
    last_error TEXT,
    processed_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE UNIQUE INDEX uq_stripe_webhook_events_event_id ON stripe_webhook_events(stripe_event_id);
CREATE INDEX idx_stripe_webhook_events_status ON stripe_webhook_events(status);
CREATE INDEX idx_stripe_webhook_events_processed_at ON stripe_webhook_events(processed_at);
CREATE INDEX idx_stripe_webhook_events_deleted_at ON stripe_webhook_events(deleted_at);

CREATE TRIGGER update_stripe_webhook_events_updated_at BEFORE UPDATE ON stripe_webhook_events
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
