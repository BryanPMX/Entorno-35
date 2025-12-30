package middleware

import (
	"net/http"

	"github.com/entorno35/backend/internal/auth"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// TenantMiddleware ensures that the request is scoped to the authenticated company
// and optionally validates that resources belong to that company (high cohesion - tenant isolation)
func TenantMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get auth context (must be set by AuthMiddleware first)
		authCtx, exists := GetAuthContext(c)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
			c.Abort()
			return
		}

		// Extract company ID from URL parameter (if present)
		companyIDParam := c.Param("company_id")
		if companyIDParam != "" {
			companyID, err := uuid.Parse(companyIDParam)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid company ID"})
				c.Abort()
				return
			}

			// Verify company ID matches authenticated user's company
			if companyID != authCtx.CompanyID {
				c.JSON(http.StatusForbidden, gin.H{"error": "access denied: company ID mismatch"})
				c.Abort()
				return
			}
		}

		// Set company ID in context for easy access
		c.Set("company_id", authCtx.CompanyID)
		if authCtx.StaffID != nil {
			c.Set("staff_id", *authCtx.StaffID)
		}

		c.Next()
	}
}

// RequireCompanyID extracts and validates company ID from context
func RequireCompanyID(c *gin.Context) (uuid.UUID, bool) {
	companyID, exists := c.Get("company_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company ID required"})
		c.Abort()
		return uuid.Nil, false
	}

	id, ok := companyID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid company ID type"})
		c.Abort()
		return uuid.Nil, false
	}

	return id, true
}

// RequireStaffID extracts and validates staff ID from context
func RequireStaffID(c *gin.Context) (uuid.UUID, bool) {
	staffID, exists := c.Get("staff_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "staff ID required"})
		c.Abort()
		return uuid.Nil, false
	}

	id, ok := staffID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid staff ID type"})
		c.Abort()
		return uuid.Nil, false
	}

	return id, true
}

// GetAuthContext is re-exported for convenience (low coupling - single import)
// It's defined in auth.go to avoid circular dependencies

