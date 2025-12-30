package middleware

import (
	"net/http"
	"strings"

	"github.com/entorno35/backend/internal/auth"
	"github.com/entorno35/backend/internal/core/jwt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	// AuthHeaderKey is the header key for the authorization token
	AuthHeaderKey = "Authorization"
	// AuthHeaderPrefix is the prefix for bearer tokens
	AuthHeaderPrefix = "Bearer "
)

// AuthMiddleware creates a middleware that validates JWT tokens (high cohesion - auth validation)
func AuthMiddleware(jwtService jwt.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract token from Authorization header
		authHeader := c.GetHeader(AuthHeaderKey)
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header required"})
			c.Abort()
			return
		}

		// Check Bearer prefix
		if !strings.HasPrefix(authHeader, AuthHeaderPrefix) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			c.Abort()
			return
		}

		// Extract token
		tokenString := strings.TrimPrefix(authHeader, AuthHeaderPrefix)
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token required"})
			c.Abort()
			return
		}

		// Validate token
		claims, err := jwtService.ValidateToken(tokenString)
		if err != nil {
			statusCode := http.StatusUnauthorized
			if err == jwt.ErrExpiredToken {
				statusCode = http.StatusUnauthorized
			}
			c.JSON(statusCode, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		// Parse UUIDs
		companyID, err := uuid.Parse(claims.CompanyID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid company ID in token"})
			c.Abort()
			return
		}

		var staffID *uuid.UUID
		if claims.StaffID != "" {
			parsedStaffID, err := uuid.Parse(claims.StaffID)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid staff ID in token"})
				c.Abort()
				return
			}
			staffID = &parsedStaffID
		}

		// Set auth context in gin context (low coupling - uses interface, not concrete type)
		authCtx := &auth.Context{
			CompanyID: companyID,
			StaffID:   staffID,
			Email:     claims.Email,
			Role:      claims.Role,
		}
		c.Set("auth_context", authCtx)

		c.Next()
	}
}

// GetAuthContext extracts the auth context from gin context (helper function)
func GetAuthContext(c *gin.Context) (*auth.Context, bool) {
	ctx, exists := c.Get("auth_context")
	if !exists {
		return nil, false
	}

	authCtx, ok := ctx.(*auth.Context)
	return authCtx, ok
}

// RequireAuth is a helper that aborts if no auth context is found
func RequireAuth(c *gin.Context) (*auth.Context, bool) {
	authCtx, exists := GetAuthContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		c.Abort()
		return nil, false
	}
	return authCtx, true
}

