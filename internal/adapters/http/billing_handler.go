package http

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/entorno35/backend/internal/adapters/postgres"
	"github.com/entorno35/backend/internal/auth"
	"github.com/entorno35/backend/internal/core/password"
	"github.com/entorno35/backend/internal/core/ports"
	"github.com/entorno35/backend/internal/core/services"
	"github.com/entorno35/backend/internal/domain"
	"github.com/entorno35/backend/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateCheckoutSessionRequest represents registration + billing checkout input.
type CreateCheckoutSessionRequest struct {
	RFC           string `json:"rfc" binding:"required"`
	CompanyName   string `json:"company_name" binding:"required"`
	Address       string `json:"address"`
	AdminEmail    string `json:"admin_email" binding:"required,email"`
	Password      string `json:"password" binding:"required,min=8"`
	Plan          string `json:"plan" binding:"required,oneof=monthly yearly"`
	EmployeeCount int    `json:"employee_count"`
}

// BillingHandler handles Stripe billing flows.
type BillingHandler struct {
	authRepo         ports.AuthRepository
	pendingRepo      ports.PendingRegistrationRepository
	webhookEventRepo ports.StripeWebhookEventRepository
	passwordHasher   password.Hasher
	stripeService    *services.StripeService
	successURL       string
	cancelURL        string
	portalReturnURL  string
}

// NewBillingHandler creates a billing handler.
func NewBillingHandler(authRepo ports.AuthRepository, pendingRepo ports.PendingRegistrationRepository, webhookEventRepo ports.StripeWebhookEventRepository, stripeService *services.StripeService, successURL, cancelURL, portalReturnURL string) *BillingHandler {
	return &BillingHandler{
		authRepo:         authRepo,
		pendingRepo:      pendingRepo,
		webhookEventRepo: webhookEventRepo,
		passwordHasher:   password.NewDefaultHasher(),
		stripeService:    stripeService,
		successURL:       successURL,
		cancelURL:        cancelURL,
		portalReturnURL:  portalReturnURL,
	}
}

const pendingRegistrationTTL = 24 * time.Hour

// CreateExistingCompanyCheckoutSessionRequest represents authenticated company billing checkout input.
type CreateExistingCompanyCheckoutSessionRequest struct {
	Plan string `json:"plan" binding:"required,oneof=monthly yearly"`
}

// CreateCheckoutSession creates a Stripe checkout session and stores a pending registration only.
// POST /billing/checkout-session
func (h *BillingHandler) CreateCheckoutSession(c *gin.Context) {
	if h.stripeService == nil || h.pendingRepo == nil || !h.stripeService.Enabled() {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "billing is not configured"})
		return
	}

	var req CreateCheckoutSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	normalizedRFC := strings.ToUpper(strings.TrimSpace(req.RFC))
	if len(normalizedRFC) < 12 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "RFC must contain at least 12 characters"})
		return
	}

	normalizedEmail := strings.ToLower(strings.TrimSpace(req.AdminEmail))
	hashedPassword, err := h.passwordHasher.Hash(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to secure credentials"})
		return
	}

	priceID, err := h.stripeService.PriceIDForPlan(req.Plan)
	if err != nil {
		if errors.Is(err, services.ErrInvalidStripePlan) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid plan"})
			return
		}
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "billing plan is not configured"})
		return
	}

	existingCompany, err := h.authRepo.GetCompanyByRFC(normalizedRFC)
	if err != nil && !errors.Is(err, postgres.ErrCompanyNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to validate company"})
		return
	}
	if err == nil && existingCompany != nil {
		if existingCompany.SubscriptionStatus == domain.SubscriptionStatusActive {
			c.JSON(http.StatusConflict, gin.H{"error": "company already has an active subscription"})
			return
		}
		c.JSON(http.StatusConflict, gin.H{"error": "company already exists. contact support to reactivate or update billing"})
		return
	}

	employeeCount := req.EmployeeCount
	if employeeCount < 0 {
		employeeCount = 0
	}
	pending := &domain.PendingCompanyRegistration{
		ID:                uuid.New(),
		RFC:               normalizedRFC,
		CompanyName:       strings.TrimSpace(req.CompanyName),
		Address:           strings.TrimSpace(req.Address),
		AdminEmail:        normalizedEmail,
		AdminPasswordHash: hashedPassword,
		EmployeeCount:     employeeCount,
		Plan:              strings.ToLower(req.Plan),
		ExpiresAt:         time.Now().UTC().Add(pendingRegistrationTTL),
	}
	if upsertErr := h.pendingRepo.UpsertPending(pending); upsertErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to store registration draft"})
		return
	}

	session, err := h.stripeService.CreateCheckoutSession(services.CreateCheckoutSessionRequest{
		PriceID:           priceID,
		CustomerEmail:     normalizedEmail,
		SuccessURL:        h.successURL,
		CancelURL:         h.cancelURL,
		ClientReferenceID: pending.ID.String(),
		Metadata: map[string]string{
			"pending_registration_id": pending.ID.String(),
			"company_rfc":             normalizedRFC,
			"plan":                    strings.ToLower(req.Plan),
		},
	})
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": fmt.Sprintf("failed to create checkout session: %v", err)})
		return
	}

	pending.StripeCheckoutSessionID = &session.ID
	if updateErr := h.pendingRepo.Update(pending); updateErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to finalize registration draft"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"session_id":   session.ID,
		"checkout_url": session.URL,
	})
}

// CreateExistingCompanyCheckoutSession creates a Stripe checkout session for an authenticated company
// to reactivate billing or start a managed subscription (existing company path).
// POST /api/v1/billing/checkout-session
func (h *BillingHandler) CreateExistingCompanyCheckoutSession(c *gin.Context) {
	if h.stripeService == nil || !h.stripeService.Enabled() {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "billing is not configured"})
		return
	}

	authCtx, ok := h.requireCompanyAuth(c)
	if !ok {
		return
	}

	var req CreateExistingCompanyCheckoutSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	company, err := h.authRepo.GetCompanyByID(authCtx.CompanyID)
	if err != nil {
		if errors.Is(err, postgres.ErrCompanyNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "company not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load company"})
		return
	}

	if company.SubscriptionStatus == domain.SubscriptionStatusActive && company.StripeSubscriptionID != nil && *company.StripeSubscriptionID != "" {
		c.JSON(http.StatusConflict, gin.H{"error": "company already has an active subscription. use billing portal to manage plan changes"})
		return
	}

	priceID, err := h.stripeService.PriceIDForPlan(req.Plan)
	if err != nil {
		if errors.Is(err, services.ErrInvalidStripePlan) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid plan"})
			return
		}
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "billing plan is not configured"})
		return
	}

	stripeReq := services.CreateCheckoutSessionRequest{
		PriceID:           priceID,
		SuccessURL:        h.successURL,
		CancelURL:         h.cancelURL,
		ClientReferenceID: company.ID.String(),
		Metadata: map[string]string{
			"company_id":     company.ID.String(),
			"company_rfc":    company.RFC,
			"billing_flow":   "existing_company",
			"requested_plan": strings.ToLower(req.Plan),
		},
	}
	if company.StripeCustomerID != nil && *company.StripeCustomerID != "" {
		stripeReq.CustomerID = *company.StripeCustomerID
	} else if company.AdminEmail != nil && *company.AdminEmail != "" {
		stripeReq.CustomerEmail = strings.ToLower(strings.TrimSpace(*company.AdminEmail))
	}

	session, err := h.stripeService.CreateCheckoutSession(stripeReq)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": fmt.Sprintf("failed to create checkout session: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"session_id":   session.ID,
		"checkout_url": session.URL,
	})
}

// CreateCustomerPortalSession creates a Stripe Billing Portal session for an authenticated company.
// POST /api/v1/billing/customer-portal
func (h *BillingHandler) CreateCustomerPortalSession(c *gin.Context) {
	if h.stripeService == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "billing is not configured"})
		return
	}

	authCtx, ok := h.requireCompanyAuth(c)
	if !ok {
		return
	}

	company, err := h.authRepo.GetCompanyByID(authCtx.CompanyID)
	if err != nil {
		if errors.Is(err, postgres.ErrCompanyNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "company not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load company"})
		return
	}
	if company.StripeCustomerID == nil || *company.StripeCustomerID == "" {
		c.JSON(http.StatusConflict, gin.H{"error": "billing customer is not configured for this company"})
		return
	}

	returnURL := h.portalReturnURL
	if returnURL == "" {
		returnURL = h.successURL
	}

	portalSession, err := h.stripeService.CreateBillingPortalSession(services.CreateBillingPortalSessionRequest{
		CustomerID: *company.StripeCustomerID,
		ReturnURL:  returnURL,
	})
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": fmt.Sprintf("failed to create billing portal session: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"session_id": portalSession.ID,
		"url":        portalSession.URL,
	})
}

// HandleWebhook processes Stripe webhook events.
// POST /billing/webhook
func (h *BillingHandler) HandleWebhook(c *gin.Context) {
	if h.stripeService == nil || h.pendingRepo == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "billing is not configured"})
		return
	}

	payload, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read webhook payload"})
		return
	}

	event, err := h.stripeService.ParseWebhookEvent(payload, c.GetHeader("Stripe-Signature"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if h.webhookEventRepo != nil {
		shouldProcess, beginErr := h.webhookEventRepo.TryBegin(event.ID, event.Type)
		if beginErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to initialize webhook processing"})
			return
		}
		if !shouldProcess {
			c.JSON(http.StatusOK, gin.H{"received": true, "duplicate": true})
			return
		}
	}

	var handleErr error
	responseStatus := http.StatusInternalServerError
	switch event.Type {
	case "checkout.session.completed":
		session, parseErr := services.ParseCheckoutSessionObject(event.Data.Object)
		if parseErr != nil {
			handleErr = parseErr
			responseStatus = http.StatusBadRequest
			break
		}
		handleErr = h.activateCompanyFromCheckoutSession(session)
	case "customer.subscription.created", "customer.subscription.updated", "customer.subscription.deleted":
		subscription, parseErr := services.ParseSubscriptionObject(event.Data.Object)
		if parseErr != nil {
			handleErr = parseErr
			responseStatus = http.StatusBadRequest
			break
		}
		handleErr = h.syncCompanyFromSubscription(subscription)
	default:
		handleErr = nil
	}

	if handleErr != nil {
		if h.webhookEventRepo != nil {
			_ = h.webhookEventRepo.MarkFailed(event.ID, handleErr.Error())
		}
		c.JSON(responseStatus, gin.H{"error": handleErr.Error()})
		return
	}

	if h.webhookEventRepo != nil {
		if err := h.webhookEventRepo.MarkProcessed(event.ID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to finalize webhook processing"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"received": true})
}

// VerifyCheckoutSession checks Stripe session state and activates subscription if payment succeeded.
// GET /billing/checkout-session/:id/verify
func (h *BillingHandler) VerifyCheckoutSession(c *gin.Context) {
	if h.stripeService == nil || h.pendingRepo == nil || !h.stripeService.Enabled() {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "billing is not configured"})
		return
	}

	sessionID := strings.TrimSpace(c.Param("id"))
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "session ID is required"})
		return
	}

	session, err := h.stripeService.GetCheckoutSession(sessionID)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": fmt.Sprintf("failed to verify session: %v", err)})
		return
	}

	if session.Mode == "subscription" && session.SubscriptionID != "" && (session.Status == "complete" || session.PaymentStatus == "paid") {
		if err := h.activateCompanyFromCheckoutSession(session); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	company, err := h.resolveCompanyForSession(session)
	if err != nil {
		if errors.Is(err, postgres.ErrCompanyNotFound) {
			pendingID, idErr := pendingRegistrationIDFromSession(session)
			if idErr != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "company not found"})
				return
			}
			pending, pendingErr := h.pendingRepo.GetByID(pendingID)
			if pendingErr != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "company not found"})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"session_id":          session.ID,
				"session_status":      session.Status,
				"payment_status":      session.PaymentStatus,
				"subscription_status": domain.SubscriptionStatusInactive,
				"active":              false,
				"company_name":        pending.CompanyName,
				"login_identifier":    pending.RFC,
				"admin_email":         pending.AdminEmail,
			})
			return
		}
		if errors.Is(err, postgres.ErrPendingRegistrationNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "company not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load company"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"session_id":          session.ID,
		"session_status":      session.Status,
		"payment_status":      session.PaymentStatus,
		"subscription_status": company.SubscriptionStatus,
		"active":              company.SubscriptionStatus == domain.SubscriptionStatusActive,
		"company_name":        company.Name,
		"login_identifier":    company.RFC,
		"admin_email":         company.AdminEmail,
	})
}

func (h *BillingHandler) activateCompanyFromCheckoutSession(session *services.CheckoutSession) error {
	if session == nil {
		return fmt.Errorf("checkout session is required")
	}
	if session.SubscriptionID == "" {
		return nil
	}

	if companyID, err := companyIDFromSession(session); err == nil {
		return h.activateExistingCompanyFromCheckoutSession(companyID, session)
	}

	pendingID, err := pendingRegistrationIDFromSession(session)
	if err != nil {
		return err
	}

	pending, err := h.pendingRepo.GetByID(pendingID)
	if err != nil {
		return fmt.Errorf("failed to fetch pending registration for activation: %w", err)
	}

	subscription, err := h.stripeService.GetSubscription(session.SubscriptionID)
	if err != nil {
		return fmt.Errorf("failed to fetch stripe subscription: %w", err)
	}

	if !isStripeSubscriptionActive(subscription.Status) && !strings.EqualFold(session.PaymentStatus, "paid") {
		pending.StripeSubscriptionID = &subscription.ID
		if session.CustomerID != "" {
			customerID := session.CustomerID
			pending.StripeCustomerID = &customerID
		}
		if pending.StripeCheckoutSessionID == nil || *pending.StripeCheckoutSessionID == "" {
			sessionID := session.ID
			pending.StripeCheckoutSessionID = &sessionID
		}
		if err := h.pendingRepo.Update(pending); err != nil {
			return fmt.Errorf("failed to persist pending subscription linkage: %w", err)
		}
		return nil
	}

	return h.activateCompanyFromPendingSubscription(pending, subscription, session.ID)
}

func (h *BillingHandler) activateExistingCompanyFromCheckoutSession(companyID uuid.UUID, session *services.CheckoutSession) error {
	if session == nil {
		return fmt.Errorf("checkout session is required")
	}
	if companyID == uuid.Nil {
		return fmt.Errorf("company ID is required")
	}
	if session.SubscriptionID == "" {
		return nil
	}

	company, err := h.authRepo.GetCompanyByID(companyID)
	if err != nil {
		return fmt.Errorf("failed to fetch company for existing-company checkout activation: %w", err)
	}

	subscription, err := h.stripeService.GetSubscription(session.SubscriptionID)
	if err != nil {
		return fmt.Errorf("failed to fetch stripe subscription: %w", err)
	}

	applySubscriptionToCompany(company, subscription)
	if session.CustomerID != "" {
		customerID := session.CustomerID
		company.StripeCustomerID = &customerID
	}
	subscriptionID := session.SubscriptionID
	company.StripeSubscriptionID = &subscriptionID

	if err := h.authRepo.UpdateCompany(company); err != nil {
		return fmt.Errorf("failed to persist company subscription: %w", err)
	}

	return nil
}

func (h *BillingHandler) syncCompanyFromSubscription(subscription *services.Subscription) error {
	if subscription == nil || subscription.ID == "" {
		return nil
	}

	company, err := h.authRepo.GetCompanyByStripeSubscriptionID(subscription.ID)
	if err != nil {
		if errors.Is(err, postgres.ErrCompanyNotFound) {
			return h.syncPendingRegistrationFromSubscription(subscription)
		}
		return fmt.Errorf("failed to fetch company by subscription ID: %w", err)
	}

	applySubscriptionToCompany(company, subscription)
	if err := h.authRepo.UpdateCompany(company); err != nil {
		return fmt.Errorf("failed to sync company subscription: %w", err)
	}

	return nil
}

func (h *BillingHandler) syncPendingRegistrationFromSubscription(subscription *services.Subscription) error {
	if h.pendingRepo == nil || subscription == nil || subscription.ID == "" {
		return nil
	}

	pending, err := h.pendingRepo.GetByStripeSubscriptionID(subscription.ID)
	if err != nil {
		if errors.Is(err, postgres.ErrPendingRegistrationNotFound) {
			return nil
		}
		return fmt.Errorf("failed to fetch pending registration by subscription ID: %w", err)
	}

	if !isStripeSubscriptionActive(subscription.Status) {
		if subscription.CustomerID != "" {
			customerID := subscription.CustomerID
			pending.StripeCustomerID = &customerID
		}
		pending.StripeSubscriptionID = &subscription.ID
		if err := h.pendingRepo.Update(pending); err != nil {
			return fmt.Errorf("failed to sync pending registration subscription linkage: %w", err)
		}
		return nil
	}

	return h.activateCompanyFromPendingSubscription(pending, subscription, "")
}

func (h *BillingHandler) activateCompanyFromPendingSubscription(pending *domain.PendingCompanyRegistration, subscription *services.Subscription, checkoutSessionID string) error {
	if pending == nil {
		return fmt.Errorf("pending registration is required")
	}
	if subscription == nil {
		return fmt.Errorf("subscription is required")
	}

	var err error
	var company *domain.Company
	if pending.CompanyID != nil {
		company, err = h.authRepo.GetCompanyByID(*pending.CompanyID)
		if err != nil && !errors.Is(err, postgres.ErrCompanyNotFound) {
			return fmt.Errorf("failed to fetch completed company from pending registration: %w", err)
		}
	}

	if company == nil {
		company, err = h.authRepo.GetCompanyByRFC(pending.RFC)
		if err != nil && !errors.Is(err, postgres.ErrCompanyNotFound) {
			return fmt.Errorf("failed to look up company by RFC during activation: %w", err)
		}
	}

	if company == nil {
		adminEmail := pending.AdminEmail
		adminPasswordHash := pending.AdminPasswordHash
		company = &domain.Company{
			ID:                 uuid.New(),
			RFC:                pending.RFC,
			Name:               pending.CompanyName,
			Address:            pending.Address,
			AdminEmail:         &adminEmail,
			AdminPasswordHash:  &adminPasswordHash,
			SubscriptionStatus: domain.SubscriptionStatusInactive,
			EmployeeCount:      pending.EmployeeCount,
		}
		if err := h.authRepo.CreateCompany(company); err != nil {
			return fmt.Errorf("failed to create company after payment: %w", err)
		}
	}

	applySubscriptionToCompany(company, subscription)
	if err := h.authRepo.UpdateCompany(company); err != nil {
		return fmt.Errorf("failed to persist company subscription: %w", err)
	}

	now := time.Now().UTC()
	pending.CompletedAt = &now
	pending.ExpiresAt = now.Add(pendingRegistrationTTL)
	pendingCompanyID := company.ID
	pending.CompanyID = &pendingCompanyID
	subscriptionID := subscription.ID
	pending.StripeSubscriptionID = &subscriptionID
	if subscription.CustomerID != "" {
		customerID := subscription.CustomerID
		pending.StripeCustomerID = &customerID
	}
	if checkoutSessionID != "" {
		sessionID := checkoutSessionID
		pending.StripeCheckoutSessionID = &sessionID
	}
	if err := h.pendingRepo.Update(pending); err != nil {
		return fmt.Errorf("failed to complete pending registration: %w", err)
	}

	return nil
}

func applySubscriptionToCompany(company *domain.Company, subscription *services.Subscription) {
	if company == nil || subscription == nil {
		return
	}

	if isStripeSubscriptionActive(subscription.Status) {
		company.SubscriptionStatus = domain.SubscriptionStatusActive
	} else {
		company.SubscriptionStatus = domain.SubscriptionStatusInactive
	}

	company.SubscriptionStartDate = subscription.CurrentPeriodStart
	company.SubscriptionEndDate = subscription.CurrentPeriodEnd

	if subscription.CustomerID != "" {
		customerID := subscription.CustomerID
		company.StripeCustomerID = &customerID
	}
	if subscription.ID != "" {
		subscriptionID := subscription.ID
		company.StripeSubscriptionID = &subscriptionID
	}
	if subscription.PriceID != "" {
		priceID := subscription.PriceID
		company.StripePriceID = &priceID
	}
}

func isStripeSubscriptionActive(status string) bool {
	switch strings.ToLower(strings.TrimSpace(status)) {
	case "active", "trialing":
		return true
	default:
		return false
	}
}

func (h *BillingHandler) requireCompanyAuth(c *gin.Context) (*auth.Context, bool) {
	authCtx, ok := middleware.RequireAuth(c)
	if !ok {
		return nil, false
	}
	if authCtx.IsStaff() || !strings.EqualFold(authCtx.Role, "company") {
		c.JSON(http.StatusForbidden, gin.H{"error": "company authentication required"})
		c.Abort()
		return nil, false
	}
	return authCtx, true
}

func companyIDFromSession(session *services.CheckoutSession) (uuid.UUID, error) {
	if session == nil {
		return uuid.Nil, fmt.Errorf("session is required")
	}

	if companyIDRaw, ok := session.Metadata["company_id"]; ok && companyIDRaw != "" {
		companyID, err := uuid.Parse(companyIDRaw)
		if err == nil {
			return companyID, nil
		}
	}

	return uuid.Nil, fmt.Errorf("company ID not found in checkout session")
}

func pendingRegistrationIDFromSession(session *services.CheckoutSession) (uuid.UUID, error) {
	if session == nil {
		return uuid.Nil, fmt.Errorf("session is required")
	}

	if pendingIDRaw, ok := session.Metadata["pending_registration_id"]; ok && pendingIDRaw != "" {
		pendingID, err := uuid.Parse(pendingIDRaw)
		if err == nil {
			return pendingID, nil
		}
	}

	if session.ClientReferenceID != "" {
		pendingID, err := uuid.Parse(session.ClientReferenceID)
		if err == nil {
			return pendingID, nil
		}
	}

	return uuid.Nil, fmt.Errorf("pending registration ID not found in checkout session")
}

func (h *BillingHandler) resolveCompanyForSession(session *services.CheckoutSession) (*domain.Company, error) {
	pendingID, err := pendingRegistrationIDFromSession(session)
	if err != nil {
		return nil, err
	}

	pending, err := h.pendingRepo.GetByID(pendingID)
	if err != nil {
		return nil, err
	}

	if pending.CompanyID != nil {
		company, companyErr := h.authRepo.GetCompanyByID(*pending.CompanyID)
		if companyErr == nil {
			return company, nil
		}
		if !errors.Is(companyErr, postgres.ErrCompanyNotFound) {
			return nil, companyErr
		}
	}

	return h.authRepo.GetCompanyByRFC(pending.RFC)
}
