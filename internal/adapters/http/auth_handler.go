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
	"github.com/google/uuid"
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
	adminEmail     string
	adminPassHash  string
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(
	jwtService jwt.Service,
	authRepo ports.AuthRepository,
	tokenExpiry time.Duration,
	adminEmail string,
	adminPassHash string,
) *AuthHandler {
	return &AuthHandler{
		jwtService:     jwtService,
		authRepo:       authRepo,
		passwordHasher: password.NewDefaultHasher(),
		tokenExpiry:    tokenExpiry,
		adminEmail:     strings.ToLower(strings.TrimSpace(adminEmail)),
		adminPassHash:  strings.TrimSpace(adminPassHash),
	}
}

// AdminLoginRequest represents the admin login request payload.
type AdminLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
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

// AdminLogin authenticates platform billing administrators and returns a JWT token.
// POST /auth/admin/login
func (h *AuthHandler) AdminLogin(c *gin.Context) {
	if h.adminEmail == "" || h.adminPassHash == "" {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "admin authentication is not configured"})
		return
	}

	var req AdminLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	email := strings.ToLower(strings.TrimSpace(req.Email))
	if email != h.adminEmail {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid admin credentials"})
		return
	}

	if verifyErr := h.passwordHasher.Verify(h.adminPassHash, req.Password); verifyErr != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid admin credentials"})
		return
	}

	token, err := h.jwtService.GenerateToken(
		uuid.Nil.String(),
		"",
		email,
		"admin",
		h.tokenExpiry,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate admin token"})
		return
	}

	c.JSON(http.StatusOK, LoginResponse{Token: token})
}
