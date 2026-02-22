-- Add company admin credentials and Stripe billing tracking fields
ALTER TABLE companies
    ADD COLUMN admin_email VARCHAR(255),
    ADD COLUMN admin_password_hash VARCHAR(255),
    ADD COLUMN stripe_customer_id VARCHAR(255),
    ADD COLUMN stripe_subscription_id VARCHAR(255),
    ADD COLUMN stripe_price_id VARCHAR(255);

CREATE INDEX idx_companies_admin_email ON companies(admin_email);
CREATE INDEX idx_companies_stripe_customer_id ON companies(stripe_customer_id);
CREATE UNIQUE INDEX idx_companies_stripe_subscription_id ON companies(stripe_subscription_id);

COMMENT ON COLUMN companies.admin_email IS 'Primary admin email used for billing registration';
COMMENT ON COLUMN companies.admin_password_hash IS 'Bcrypt hashed admin password for company login';
COMMENT ON COLUMN companies.stripe_customer_id IS 'Stripe customer identifier';
COMMENT ON COLUMN companies.stripe_subscription_id IS 'Stripe subscription identifier';
COMMENT ON COLUMN companies.stripe_price_id IS 'Stripe price identifier for active subscription';
