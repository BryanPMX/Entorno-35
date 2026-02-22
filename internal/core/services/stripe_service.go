package services

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var (
	ErrStripeNotConfigured     = errors.New("stripe is not configured")
	ErrInvalidStripePlan       = errors.New("invalid stripe plan")
	ErrInvalidWebhookSignature = errors.New("invalid stripe webhook signature")
	ErrMissingWebhookSignature = errors.New("missing stripe-signature header")
)

const stripeWebhookTolerance = 5 * time.Minute

// StripeServiceConfig contains required Stripe configuration.
type StripeServiceConfig struct {
	SecretKey      string
	WebhookSecret  string
	MonthlyPriceID string
	YearlyPriceID  string
	SuccessURL     string
	CancelURL      string
}

// CreateCheckoutSessionRequest represents a Stripe checkout creation request.
type CreateCheckoutSessionRequest struct {
	PriceID           string
	CustomerEmail     string
	SuccessURL        string
	CancelURL         string
	ClientReferenceID string
	Metadata          map[string]string
}

// CheckoutSession represents a normalized Stripe checkout session.
type CheckoutSession struct {
	ID                string
	URL               string
	CustomerID        string
	SubscriptionID    string
	PaymentStatus     string
	Status            string
	Mode              string
	ClientReferenceID string
	Metadata          map[string]string
}

// Subscription represents a normalized Stripe subscription.
type Subscription struct {
	ID                 string
	Status             string
	CustomerID         string
	PriceID            string
	CurrentPeriodStart *time.Time
	CurrentPeriodEnd   *time.Time
}

// WebhookEvent is a minimal Stripe webhook envelope.
type WebhookEvent struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	Data struct {
		Object json.RawMessage `json:"object"`
	} `json:"data"`
}

// ParseCheckoutSessionObject parses a Stripe webhook object into CheckoutSession.
func ParseCheckoutSessionObject(raw json.RawMessage) (*CheckoutSession, error) {
	var session stripeCheckoutSessionResponse
	if err := json.Unmarshal(raw, &session); err != nil {
		return nil, fmt.Errorf("failed to decode checkout session object: %w", err)
	}
	return session.normalize(), nil
}

// ParseSubscriptionObject parses a Stripe webhook object into Subscription.
func ParseSubscriptionObject(raw json.RawMessage) (*Subscription, error) {
	var subscription stripeSubscriptionResponse
	if err := json.Unmarshal(raw, &subscription); err != nil {
		return nil, fmt.Errorf("failed to decode subscription object: %w", err)
	}
	return subscription.normalize(), nil
}

type StripeService struct {
	httpClient *http.Client
	cfg        StripeServiceConfig
	baseURL    string
}

// NewStripeService creates a new Stripe service.
func NewStripeService(cfg StripeServiceConfig) *StripeService {
	return &StripeService{
		httpClient: &http.Client{Timeout: 20 * time.Second},
		cfg:        cfg,
		baseURL:    "https://api.stripe.com",
	}
}

// Enabled returns true when Stripe checkout can be used.
func (s *StripeService) Enabled() bool {
	return s.cfg.SecretKey != "" && s.cfg.MonthlyPriceID != "" && s.cfg.YearlyPriceID != ""
}

// PriceIDForPlan resolves Stripe price ID for known plans.
func (s *StripeService) PriceIDForPlan(plan string) (string, error) {
	switch strings.ToLower(strings.TrimSpace(plan)) {
	case "monthly":
		if s.cfg.MonthlyPriceID == "" {
			return "", ErrStripeNotConfigured
		}
		return s.cfg.MonthlyPriceID, nil
	case "yearly":
		if s.cfg.YearlyPriceID == "" {
			return "", ErrStripeNotConfigured
		}
		return s.cfg.YearlyPriceID, nil
	default:
		return "", ErrInvalidStripePlan
	}
}

// CreateCheckoutSession creates a Stripe Checkout session in subscription mode.
func (s *StripeService) CreateCheckoutSession(req CreateCheckoutSessionRequest) (*CheckoutSession, error) {
	if !s.Enabled() {
		return nil, ErrStripeNotConfigured
	}

	values := url.Values{}
	values.Set("mode", "subscription")
	values.Set("line_items[0][price]", req.PriceID)
	values.Set("line_items[0][quantity]", "1")
	values.Set("success_url", req.SuccessURL)
	values.Set("cancel_url", req.CancelURL)

	if req.CustomerEmail != "" {
		values.Set("customer_email", req.CustomerEmail)
	}
	if req.ClientReferenceID != "" {
		values.Set("client_reference_id", req.ClientReferenceID)
	}
	for key, value := range req.Metadata {
		values.Set(fmt.Sprintf("metadata[%s]", key), value)
	}

	body, err := s.doFormRequest(http.MethodPost, "/v1/checkout/sessions", values)
	if err != nil {
		return nil, err
	}

	var stripeSession stripeCheckoutSessionResponse
	if err := json.Unmarshal(body, &stripeSession); err != nil {
		return nil, fmt.Errorf("failed to decode stripe checkout response: %w", err)
	}

	return stripeSession.normalize(), nil
}

// GetCheckoutSession fetches a Stripe Checkout session by ID.
func (s *StripeService) GetCheckoutSession(sessionID string) (*CheckoutSession, error) {
	if !s.Enabled() {
		return nil, ErrStripeNotConfigured
	}
	if strings.TrimSpace(sessionID) == "" {
		return nil, fmt.Errorf("session ID is required")
	}

	body, err := s.doRequest(http.MethodGet, "/v1/checkout/sessions/"+url.PathEscape(sessionID), nil, "")
	if err != nil {
		return nil, err
	}

	var stripeSession stripeCheckoutSessionResponse
	if err := json.Unmarshal(body, &stripeSession); err != nil {
		return nil, fmt.Errorf("failed to decode stripe checkout session: %w", err)
	}

	return stripeSession.normalize(), nil
}

// GetSubscription fetches a Stripe subscription by ID.
func (s *StripeService) GetSubscription(subscriptionID string) (*Subscription, error) {
	if !s.Enabled() {
		return nil, ErrStripeNotConfigured
	}
	if strings.TrimSpace(subscriptionID) == "" {
		return nil, fmt.Errorf("subscription ID is required")
	}

	body, err := s.doRequest(http.MethodGet, "/v1/subscriptions/"+url.PathEscape(subscriptionID), nil, "")
	if err != nil {
		return nil, err
	}

	var stripeSubscription stripeSubscriptionResponse
	if err := json.Unmarshal(body, &stripeSubscription); err != nil {
		return nil, fmt.Errorf("failed to decode stripe subscription: %w", err)
	}

	return stripeSubscription.normalize(), nil
}

// ParseWebhookEvent validates the webhook signature and decodes the event payload.
func (s *StripeService) ParseWebhookEvent(payload []byte, signatureHeader string) (*WebhookEvent, error) {
	if s.cfg.WebhookSecret == "" {
		return nil, ErrStripeNotConfigured
	}
	if err := verifyWebhookSignature(payload, signatureHeader, s.cfg.WebhookSecret); err != nil {
		return nil, err
	}

	var event WebhookEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return nil, fmt.Errorf("failed to decode webhook event: %w", err)
	}

	return &event, nil
}

func (s *StripeService) doFormRequest(method, path string, values url.Values) ([]byte, error) {
	return s.doRequest(method, path, strings.NewReader(values.Encode()), "application/x-www-form-urlencoded")
}

func (s *StripeService) doRequest(method, path string, body io.Reader, contentType string) ([]byte, error) {
	request, err := http.NewRequest(method, s.baseURL+path, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create stripe request: %w", err)
	}

	request.Header.Set("Authorization", "Bearer "+s.cfg.SecretKey)
	if contentType != "" {
		request.Header.Set("Content-Type", contentType)
	}

	response, err := s.httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("stripe request failed: %w", err)
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read stripe response: %w", err)
	}

	if response.StatusCode >= http.StatusBadRequest {
		message := extractStripeErrorMessage(responseBody)
		if message == "" {
			message = string(bytes.TrimSpace(responseBody))
		}
		return nil, fmt.Errorf("stripe API error (%d): %s", response.StatusCode, message)
	}

	return responseBody, nil
}

func verifyWebhookSignature(payload []byte, signatureHeader, secret string) error {
	if signatureHeader == "" {
		return ErrMissingWebhookSignature
	}

	var timestamp string
	var signatures []string

	parts := strings.Split(signatureHeader, ",")
	for _, part := range parts {
		kv := strings.SplitN(strings.TrimSpace(part), "=", 2)
		if len(kv) != 2 {
			continue
		}
		switch kv[0] {
		case "t":
			timestamp = kv[1]
		case "v1":
			signatures = append(signatures, kv[1])
		}
	}

	if timestamp == "" || len(signatures) == 0 {
		return ErrInvalidWebhookSignature
	}

	ts, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return ErrInvalidWebhookSignature
	}

	webhookTime := time.Unix(ts, 0)
	if time.Since(webhookTime) > stripeWebhookTolerance || time.Until(webhookTime) > stripeWebhookTolerance {
		return ErrInvalidWebhookSignature
	}

	signedPayload := timestamp + "." + string(payload)
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(signedPayload))
	expected := hex.EncodeToString(mac.Sum(nil))

	for _, signature := range signatures {
		if hmac.Equal([]byte(signature), []byte(expected)) {
			return nil
		}
	}

	return ErrInvalidWebhookSignature
}

func extractStripeErrorMessage(body []byte) string {
	var errResp struct {
		Error struct {
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.Unmarshal(body, &errResp); err != nil {
		return ""
	}
	return errResp.Error.Message
}

type stripeCheckoutSessionResponse struct {
	ID                string            `json:"id"`
	URL               string            `json:"url"`
	Customer          json.RawMessage   `json:"customer"`
	Subscription      json.RawMessage   `json:"subscription"`
	PaymentStatus     string            `json:"payment_status"`
	Status            string            `json:"status"`
	Mode              string            `json:"mode"`
	ClientReferenceID string            `json:"client_reference_id"`
	Metadata          map[string]string `json:"metadata"`
}

func (s stripeCheckoutSessionResponse) normalize() *CheckoutSession {
	return &CheckoutSession{
		ID:                s.ID,
		URL:               s.URL,
		CustomerID:        parseExpandableID(s.Customer),
		SubscriptionID:    parseExpandableID(s.Subscription),
		PaymentStatus:     s.PaymentStatus,
		Status:            s.Status,
		Mode:              s.Mode,
		ClientReferenceID: s.ClientReferenceID,
		Metadata:          s.Metadata,
	}
}

type stripeSubscriptionResponse struct {
	ID                 string          `json:"id"`
	Status             string          `json:"status"`
	Customer           json.RawMessage `json:"customer"`
	CurrentPeriodStart int64           `json:"current_period_start"`
	CurrentPeriodEnd   int64           `json:"current_period_end"`
	Items              struct {
		Data []struct {
			Price struct {
				ID string `json:"id"`
			} `json:"price"`
		} `json:"data"`
	} `json:"items"`
}

func (s stripeSubscriptionResponse) normalize() *Subscription {
	var priceID string
	if len(s.Items.Data) > 0 {
		priceID = s.Items.Data[0].Price.ID
	}

	return &Subscription{
		ID:                 s.ID,
		Status:             s.Status,
		CustomerID:         parseExpandableID(s.Customer),
		PriceID:            priceID,
		CurrentPeriodStart: unixToTimePointer(s.CurrentPeriodStart),
		CurrentPeriodEnd:   unixToTimePointer(s.CurrentPeriodEnd),
	}
}

func parseExpandableID(raw json.RawMessage) string {
	if len(raw) == 0 || string(raw) == "null" {
		return ""
	}

	var asString string
	if err := json.Unmarshal(raw, &asString); err == nil {
		return asString
	}

	var asObject struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(raw, &asObject); err == nil {
		return asObject.ID
	}

	return ""
}

func unixToTimePointer(unix int64) *time.Time {
	if unix <= 0 {
		return nil
	}
	value := time.Unix(unix, 0).UTC()
	return &value
}
