package http

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/entorno35/backend/internal/adapters/postgres"
	"github.com/entorno35/backend/internal/core/jwt"
	"github.com/entorno35/backend/internal/core/password"
	"github.com/entorno35/backend/internal/core/ports"
	"github.com/entorno35/backend/internal/domain"
	"github.com/gin-gonic/gin"
)

// LoginRequest represents the login request payload
type LoginRequest struct {
	Identifier string `json:"identifier" binding:"required"` // RFC for COMPANY
	Type       string `json:"type" binding:"required"`       // Must be "COMPANY"
	Password   string `json:"password" binding:"required"`   // Admin password
}

// LoginResponse represents the login response payload
type LoginResponse struct {
	Token string `json:"token"`
}

// AuthHandler handles authentication HTTP requests (High cohesion - HTTP concerns only)
type AuthHandler struct {
	jwtService     jwt.Service
	authRepo       ports.AuthRepository
	passwordHasher password.Hasher
	tokenExpiry    time.Duration
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(jwtService jwt.Service, authRepo ports.AuthRepository, tokenExpiry time.Duration) *AuthHandler {
	return &AuthHandler{
		jwtService:     jwtService,
		authRepo:       authRepo,
		passwordHasher: password.NewDefaultHasher(),
		tokenExpiry:    tokenExpiry,
	}
}

// Login handles company login and returns a JWT token
// POST /auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate type - only COMPANY is supported
	if req.Type != "COMPANY" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "only COMPANY login is supported"})
		return
	}

	normalizedRFC := strings.ToUpper(strings.TrimSpace(req.Identifier))

	// Find company by RFC
	company, err := h.authRepo.GetCompanyByRFC(normalizedRFC)
	if err != nil {
		if errors.Is(err, postgres.ErrCompanyNotFound) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if company.SubscriptionStatus != domain.SubscriptionStatusActive {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "subscription is inactive. complete payment to activate your account"})
		return
	}

	if company.AdminPasswordHash == nil || *company.AdminPasswordHash == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	if verifyErr := h.passwordHasher.Verify(*company.AdminPasswordHash, req.Password); verifyErr != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// Generate token for company
	token, err := h.jwtService.GenerateToken(
		company.ID.String(),
		"", // No staff ID for company-level authentication
		"", // No email for company
		"company",
		h.tokenExpiry,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, LoginResponse{Token: token})
}
