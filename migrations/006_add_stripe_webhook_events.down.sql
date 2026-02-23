DROP TRIGGER IF EXISTS update_stripe_webhook_events_updated_at ON stripe_webhook_events;
DROP INDEX IF EXISTS idx_stripe_webhook_events_deleted_at;
DROP INDEX IF EXISTS idx_stripe_webhook_events_processed_at;
DROP INDEX IF EXISTS idx_stripe_webhook_events_status;
DROP INDEX IF EXISTS uq_stripe_webhook_events_event_id;
DROP TABLE IF EXISTS stripe_webhook_events;
