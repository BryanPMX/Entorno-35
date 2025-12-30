package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token has expired")
)

// Claims represents JWT claims structure
type Claims struct {
	CompanyID string `json:"company_id"`
	StaffID   string `json:"staff_id,omitempty"`
	Email     string `json:"email,omitempty"`
	Role      string `json:"role,omitempty"`
	jwt.RegisteredClaims
}

// Service defines the interface for JWT operations (low coupling via interface)
type Service interface {
	GenerateToken(companyID, staffID, email, role string, expiresIn time.Duration) (string, error)
	ValidateToken(tokenString string) (*Claims, error)
	RefreshToken(tokenString string, expiresIn time.Duration) (string, error)
}

// jwtService implements the JWT Service interface (high cohesion - single responsibility)
type jwtService struct {
	secretKey []byte
}

// NewService creates a new JWT service instance
func NewService(secretKey string) Service {
	return &jwtService{
		secretKey: []byte(secretKey),
	}
}

// GenerateToken creates a new JWT token with the provided claims
func (s *jwtService) GenerateToken(companyID, staffID, email, role string, expiresIn time.Duration) (string, error) {
	now := time.Now()
	claims := &Claims{
		CompanyID: companyID,
		StaffID:   staffID,
		Email:     email,
		Role:      role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(expiresIn)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "entorno35",
			Subject:   staffID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secretKey)
}

// ValidateToken validates and parses a JWT token string
func (s *jwtService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return s.secretKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

// RefreshToken generates a new token from an existing valid token
func (s *jwtService) RefreshToken(tokenString string, expiresIn time.Duration) (string, error) {
	claims, err := s.ValidateToken(tokenString)
	if err != nil {
		return "", err
	}

	// Generate new token with same claims but new expiration
	return s.GenerateToken(claims.CompanyID, claims.StaffID, claims.Email, claims.Role, expiresIn)
}

