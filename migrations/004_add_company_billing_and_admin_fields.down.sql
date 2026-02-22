DROP INDEX IF EXISTS idx_companies_stripe_subscription_id;
DROP INDEX IF EXISTS idx_companies_stripe_customer_id;
DROP INDEX IF EXISTS idx_companies_admin_email;

ALTER TABLE companies
    DROP COLUMN IF EXISTS stripe_price_id,
    DROP COLUMN IF EXISTS stripe_subscription_id,
    DROP COLUMN IF EXISTS stripe_customer_id,
    DROP COLUMN IF EXISTS admin_password_hash,
    DROP COLUMN IF EXISTS admin_email;
