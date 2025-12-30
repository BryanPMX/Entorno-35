package http

import (
	"errors"
	"net/http"
	"time"

	"github.com/entorno35/backend/internal/adapters/postgres"
	"github.com/entorno35/backend/internal/core/jwt"
	"github.com/entorno35/backend/internal/core/ports"
	"github.com/gin-gonic/gin"
)

// Note: For MVP, we validate identifier (RFC/CURP) exists without password verification
// In production, add password field to Staff model and verify using password.Hasher

// LoginRequest represents the login request payload
type LoginRequest struct {
	Identifier string `json:"identifier" binding:"required"` // RFC for COMPANY, CURP for STAFF
	Type       string `json:"type" binding:"required"`       // "COMPANY" or "STAFF"
	CompanyID  string `json:"company_id,omitempty"`          // Required for STAFF type
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

// Login handles user login and returns a JWT token
// POST /auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate type
	if req.Type != "COMPANY" && req.Type != "STAFF" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "type must be COMPANY or STAFF"})
		return
	}

	var token string
	var err error

	if req.Type == "COMPANY" {
		// Find company by RFC
		company, err := h.authRepo.GetCompanyByRFC(req.Identifier)
		if err != nil {
			statusCode := http.StatusInternalServerError
			// Map repository errors to appropriate HTTP status codes
			if errors.Is(err, postgres.ErrCompanyNotFound) || errors.Is(err, postgres.ErrCompanyInactive) {
				statusCode = http.StatusUnauthorized
			}
			c.JSON(statusCode, gin.H{"error": err.Error()})
			return
		}

		// Generate token for company (no staff ID)
		token, err = h.jwtService.GenerateToken(
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
	} else {
		// STAFF type - requires company_id
		if req.CompanyID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "company_id is required for STAFF type"})
			return
		}

		// Find staff by CURP and company ID
		staff, err := h.authRepo.GetStaffByCURP(req.Identifier, req.CompanyID)
		if err != nil {
			statusCode := http.StatusInternalServerError
			if errors.Is(err, postgres.ErrStaffNotFound) || errors.Is(err, postgres.ErrCompanyInactive) || errors.Is(err, postgres.ErrCompanyNotFound) {
				statusCode = http.StatusUnauthorized
			}
			c.JSON(statusCode, gin.H{"error": err.Error()})
			return
		}

		// Generate token for staff
		token, err = h.jwtService.GenerateToken(
			staff.CompanyID.String(),
			staff.ID.String(),
			staff.Email,
			"staff",
			h.tokenExpiry,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
			return
		}
	}

	c.JSON(http.StatusOK, LoginResponse{Token: token})
}

