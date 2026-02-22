package http

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/entorno35/backend/internal/adapters/postgres"
	"github.com/entorno35/backend/internal/core/password"
	"github.com/entorno35/backend/internal/core/ports"
	"github.com/entorno35/backend/internal/core/services"
	"github.com/entorno35/backend/internal/domain"
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
	authRepo       ports.AuthRepository
	passwordHasher password.Hasher
	stripeService  *services.StripeService
	successURL     string
	cancelURL      string
}

// NewBillingHandler creates a billing handler.
func NewBillingHandler(authRepo ports.AuthRepository, stripeService *services.StripeService, successURL, cancelURL string) *BillingHandler {
	return &BillingHandler{
		authRepo:       authRepo,
		passwordHasher: password.NewDefaultHasher(),
		stripeService:  stripeService,
		successURL:     successURL,
		cancelURL:      cancelURL,
	}
}

// CreateCheckoutSession creates a Stripe checkout session and persists inactive registration credentials.
// POST /billing/checkout-session
func (h *BillingHandler) CreateCheckoutSession(c *gin.Context) {
	if h.stripeService == nil || !h.stripeService.Enabled() {
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

	company, err := h.authRepo.GetCompanyByRFC(normalizedRFC)
	if err != nil && !errors.Is(err, postgres.ErrCompanyNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to validate company"})
		return
	}

	adminEmail := normalizedEmail
	adminPasswordHash := hashedPassword

	if errors.Is(err, postgres.ErrCompanyNotFound) {
		employeeCount := req.EmployeeCount
		if employeeCount < 0 {
			employeeCount = 0
		}
		company = &domain.Company{
			ID:                 uuid.New(),
			RFC:                normalizedRFC,
			Name:               strings.TrimSpace(req.CompanyName),
			Address:            strings.TrimSpace(req.Address),
			AdminEmail:         &adminEmail,
			AdminPasswordHash:  &adminPasswordHash,
			SubscriptionStatus: domain.SubscriptionStatusInactive,
			EmployeeCount:      employeeCount,
		}

		if createErr := h.authRepo.CreateCompany(company); createErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to register company"})
			return
		}
	} else {
		if company.SubscriptionStatus == domain.SubscriptionStatusActive {
			c.JSON(http.StatusConflict, gin.H{"error": "company already has an active subscription"})
			return
		}

		company.Name = strings.TrimSpace(req.CompanyName)
		company.Address = strings.TrimSpace(req.Address)
		company.AdminEmail = &adminEmail
		company.AdminPasswordHash = &adminPasswordHash
		company.SubscriptionStatus = domain.SubscriptionStatusInactive
		company.SubscriptionStartDate = nil
		company.SubscriptionEndDate = nil
		company.StripeCustomerID = nil
		company.StripeSubscriptionID = nil
		company.StripePriceID = nil
		if req.EmployeeCount >= 0 {
			company.EmployeeCount = req.EmployeeCount
		}

		if updateErr := h.authRepo.UpdateCompany(company); updateErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update registration"})
			return
		}
	}

	session, err := h.stripeService.CreateCheckoutSession(services.CreateCheckoutSessionRequest{
		PriceID:           priceID,
		CustomerEmail:     normalizedEmail,
		SuccessURL:        h.successURL,
		CancelURL:         h.cancelURL,
		ClientReferenceID: company.ID.String(),
		Metadata: map[string]string{
			"company_id":  company.ID.String(),
			"company_rfc": normalizedRFC,
			"plan":        strings.ToLower(req.Plan),
		},
	})
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": fmt.Sprintf("failed to create checkout session: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"session_id":   session.ID,
		"checkout_url": session.URL,
	})
}

// HandleWebhook processes Stripe webhook events.
// POST /billing/webhook
func (h *BillingHandler) HandleWebhook(c *gin.Context) {
	if h.stripeService == nil {
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

	switch event.Type {
	case "checkout.session.completed":
		session, parseErr := services.ParseCheckoutSessionObject(event.Data.Object)
		if parseErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": parseErr.Error()})
			return
		}
		if handleErr := h.activateCompanyFromCheckoutSession(session); handleErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": handleErr.Error()})
			return
		}
	case "customer.subscription.created", "customer.subscription.updated", "customer.subscription.deleted":
		subscription, parseErr := services.ParseSubscriptionObject(event.Data.Object)
		if parseErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": parseErr.Error()})
			return
		}
		if handleErr := h.syncCompanyFromSubscription(subscription); handleErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": handleErr.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"received": true})
}

// VerifyCheckoutSession checks Stripe session state and activates subscription if payment succeeded.
// GET /billing/checkout-session/:id/verify
func (h *BillingHandler) VerifyCheckoutSession(c *gin.Context) {
	if h.stripeService == nil || !h.stripeService.Enabled() {
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

	companyID, err := companyIDFromSession(session)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	company, err := h.authRepo.GetCompanyByID(companyID)
	if err != nil {
		if errors.Is(err, postgres.ErrCompanyNotFound) {
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

	companyID, err := companyIDFromSession(session)
	if err != nil {
		return err
	}

	company, err := h.authRepo.GetCompanyByID(companyID)
	if err != nil {
		return fmt.Errorf("failed to fetch company for activation: %w", err)
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
			return nil
		}
		return fmt.Errorf("failed to fetch company by subscription ID: %w", err)
	}

	applySubscriptionToCompany(company, subscription)
	if err := h.authRepo.UpdateCompany(company); err != nil {
		return fmt.Errorf("failed to sync company subscription: %w", err)
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

	if session.ClientReferenceID != "" {
		companyID, err := uuid.Parse(session.ClientReferenceID)
		if err == nil {
			return companyID, nil
		}
	}

	return uuid.Nil, fmt.Errorf("company ID not found in checkout session")
}
