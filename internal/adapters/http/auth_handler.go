package http

import (
	"errors"
	"net/http"
	"time"

	"github.com/entorno35/backend/internal/adapters/postgres"
	"github.com/entorno35/backend/internal/core/jwt"
	"github.com/entorno35/backend/internal/core/ports"
	"github.com/entorno35/backend/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// LoginRequest represents the login request payload
type LoginRequest struct {
	Identifier string `json:"identifier" binding:"required"` // RFC for COMPANY
	Type       string `json:"type" binding:"required"`       // Must be "COMPANY"
}

// LoginResponse represents the login response payload
type LoginResponse struct {
	Token string `json:"token"`
}

// AuthHandler handles authentication HTTP requests (High cohesion - HTTP concerns only)
type AuthHandler struct {
	jwtService    jwt.Service
	authRepo      ports.AuthRepository
	tokenExpiry   time.Duration
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(jwtService jwt.Service, authRepo ports.AuthRepository, tokenExpiry time.Duration) *AuthHandler {
	return &AuthHandler{
		jwtService:  jwtService,
		authRepo:    authRepo,
		tokenExpiry: tokenExpiry,
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

	// Find company by RFC, or create it if it doesn't exist (for demo purposes)
	company, err := h.authRepo.GetCompanyByRFC(req.Identifier)
	if err != nil {
		// If company not found, auto-create it for demo purposes
		if errors.Is(err, postgres.ErrCompanyNotFound) {
			// Generate UUID in Go to ensure we have it after creation
			// (database default won't populate the struct field)
			company = &domain.Company{
				ID:                 uuid.New(),
				RFC:                req.Identifier,
				Name:               "Demo Company - " + req.Identifier,
				SubscriptionStatus: domain.SubscriptionStatusActive,
				EmployeeCount:      50, // Default for demo
			}
			// Create company in database
			if createErr := h.authRepo.CreateCompany(company); createErr != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create demo company"})
				return
			}
		} else if errors.Is(err, postgres.ErrCompanyInactive) {
			statusCode := http.StatusUnauthorized
			c.JSON(statusCode, gin.H{"error": err.Error()})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
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

