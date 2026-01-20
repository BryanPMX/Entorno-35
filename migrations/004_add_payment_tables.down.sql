-- Drop triggers
DROP TRIGGER IF EXISTS update_payments_updated_at ON payments;
DROP TRIGGER IF EXISTS update_subscriptions_updated_at ON subscriptions;

-- Drop RLS policies
DROP POLICY IF EXISTS payments_company_policy ON payments;
DROP POLICY IF EXISTS subscriptions_company_policy ON subscriptions;

-- Disable RLS
ALTER TABLE payments DISABLE ROW LEVEL SECURITY;
ALTER TABLE subscriptions DISABLE ROW LEVEL SECURITY;

-- Drop indexes
DROP INDEX IF EXISTS idx_payments_company_id;
DROP INDEX IF EXISTS idx_payments_stripe_payment_id;
DROP INDEX IF EXISTS idx_payments_status;
DROP INDEX IF EXISTS idx_subscriptions_company_id;
DROP INDEX IF EXISTS idx_subscriptions_stripe_subscription_id;
DROP INDEX IF EXISTS idx_subscriptions_status;

-- Drop tables
DROP TABLE IF EXISTS subscriptions;
DROP TABLE IF EXISTS payments;

-- Drop function
DROP FUNCTION IF EXISTS update_updated_at_column();