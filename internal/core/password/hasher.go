package password

import (
	"golang.org/x/crypto/bcrypt"
)

const (
	// DefaultCost is the default bcrypt cost
	DefaultCost = bcrypt.DefaultCost
)

// Hasher defines the interface for password hashing operations (low coupling)
type Hasher interface {
	Hash(password string) (string, error)
	Verify(hashedPassword, password string) error
}

// bcryptHasher implements the Hasher interface (high cohesion - single responsibility)
type bcryptHasher struct {
	cost int
}

// NewHasher creates a new password hasher with the specified cost
func NewHasher(cost int) Hasher {
	if cost < bcrypt.MinCost || cost > bcrypt.MaxCost {
		cost = DefaultCost
	}
	return &bcryptHasher{cost: cost}
}

// NewDefaultHasher creates a new password hasher with default cost
func NewDefaultHasher() Hasher {
	return &bcryptHasher{cost: DefaultCost}
}

// Hash hashes a plain text password using bcrypt
func (h *bcryptHasher) Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// Verify compares a hashed password with a plain text password
func (h *bcryptHasher) Verify(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

