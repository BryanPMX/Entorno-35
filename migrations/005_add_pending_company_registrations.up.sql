-- Store pre-payment registrations until Stripe confirms checkout
CREATE TABLE pending_company_registrations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    rfc VARCHAR(13) NOT NULL,
    company_name VARCHAR(255) NOT NULL,
    address TEXT,
    admin_email VARCHAR(255) NOT NULL,
    admin_password_hash VARCHAR(255) NOT NULL,
    employee_count INTEGER NOT NULL DEFAULT 0,
    plan VARCHAR(20) NOT NULL,
    stripe_checkout_session_id VARCHAR(255),
    stripe_customer_id VARCHAR(255),
    stripe_subscription_id VARCHAR(255),
    company_id UUID REFERENCES companies(id) ON DELETE SET NULL,
    expires_at TIMESTAMP NOT NULL,
    completed_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_pending_company_registrations_rfc ON pending_company_registrations(rfc);
CREATE INDEX idx_pending_company_registrations_session_id ON pending_company_registrations(stripe_checkout_session_id);
CREATE INDEX idx_pending_company_registrations_company_id ON pending_company_registrations(company_id);
CREATE INDEX idx_pending_company_registrations_expires_at ON pending_company_registrations(expires_at);
CREATE INDEX idx_pending_company_registrations_completed_at ON pending_company_registrations(completed_at);
CREATE INDEX idx_pending_company_registrations_deleted_at ON pending_company_registrations(deleted_at);

-- One open pending registration per RFC (allows historical completed rows)
CREATE UNIQUE INDEX uq_pending_company_registrations_open_rfc
    ON pending_company_registrations (rfc)
    WHERE completed_at IS NULL AND deleted_at IS NULL;

CREATE TRIGGER update_pending_company_registrations_updated_at BEFORE UPDATE ON pending_company_registrations
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
