-- Store authenticated refund requests for manual billing review.
CREATE TABLE billing_refund_requests (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    stripe_customer_id VARCHAR(255),
    stripe_subscription_id VARCHAR(255),
    requested_by_email VARCHAR(255),
    reason TEXT NOT NULL CHECK (char_length(reason) BETWEEN 20 AND 2000),
    status VARCHAR(20) NOT NULL DEFAULT 'requested' CHECK (status IN ('requested', 'approved', 'rejected', 'refunded')),
    reviewed_by_email VARCHAR(255),
    reviewed_at TIMESTAMP,
    resolution_note TEXT,
    stripe_refund_id VARCHAR(255),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_billing_refund_requests_company_id ON billing_refund_requests(company_id);
CREATE INDEX idx_billing_refund_requests_status ON billing_refund_requests(status);
CREATE INDEX idx_billing_refund_requests_stripe_customer_id ON billing_refund_requests(stripe_customer_id);
CREATE INDEX idx_billing_refund_requests_stripe_subscription_id ON billing_refund_requests(stripe_subscription_id);
CREATE INDEX idx_billing_refund_requests_reviewed_at ON billing_refund_requests(reviewed_at);
CREATE INDEX idx_billing_refund_requests_stripe_refund_id ON billing_refund_requests(stripe_refund_id);
CREATE INDEX idx_billing_refund_requests_deleted_at ON billing_refund_requests(deleted_at);

-- Allow only one open refund request per company.
CREATE UNIQUE INDEX uq_billing_refund_requests_open_company
    ON billing_refund_requests(company_id)
    WHERE status = 'requested' AND deleted_at IS NULL;

CREATE TRIGGER update_billing_refund_requests_updated_at BEFORE UPDATE ON billing_refund_requests
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
