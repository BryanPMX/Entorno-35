package services

import (
	"regexp"
	"strings"
	"unicode"
)

const (
	// BatchSize is the number of records to process before flushing to database
	BatchSize = 1000

	// CURP format: 4 letters + 6 digits + 2 letters + 6 alphanumeric = 18 characters total
	// Example: ABCD123456HIJKLM01
	// Breakdown: ABCD (4 letters) + 123456 (6 digits) + HI (2 letters) + JKLM01 (6 alphanumeric)
	// Pattern: ^[A-Z]{4}[0-9]{6}[A-Z]{2}[A-Z0-9]{6}$
	curpPattern = `^[A-Z]{4}[0-9]{6}[A-Z]{2}[A-Z0-9]{6}$`
)

var (
	curpRegex = regexp.MustCompile(curpPattern)
)

// ValidateCURPFormat validates CURP format beyond just length
// CURP must be:
// - Exactly 18 characters
// - Format: 4 letters + 6 digits + 2 letters + 6 alphanumeric
// - All uppercase (normalized before validation)
func ValidateCURPFormat(curp string) error {
	// Normalize to uppercase and trim
	curp = strings.ToUpper(strings.TrimSpace(curp))

	// Check length
	if len(curp) != 18 {
		return &ValidationError{
			Field:   "CURP",
			Value:   curp,
			Message: "CURP must be exactly 18 characters",
		}
	}

	// Check for invalid characters first (before regex, so we can give a clearer error message)
	for _, r := range curp {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return &ValidationError{
				Field:   "CURP",
				Value:   curp,
				Message: "CURP contains invalid characters (only letters and digits allowed)",
			}
		}
	}

	// Check format: 4 letters + 6 digits + 2 letters + 6 alphanumeric
	// Pattern: ^[A-Z]{4}[0-9]{6}[A-Z]{2}[A-Z0-9]{6}$
	if !curpRegex.MatchString(curp) {
		return &ValidationError{
			Field:   "CURP",
			Value:   curp,
			Message: "CURP format is invalid. Expected format: 4 letters + 6 digits + 2 letters + 6 alphanumeric (e.g., ABCD123456HIJKLM01)",
		}
	}

	return nil
}

// ValidationError represents a validation error with field and value context
type ValidationError struct {
	Field   string
	Value   string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

// ValidateName validates that name is not empty and doesn't contain only whitespace
func ValidateName(name string) error {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		return &ValidationError{
			Field:   "Name",
			Value:   name,
			Message: "Name is required and cannot be empty",
		}
	}

	// Check for suspicious patterns (e.g., all same character)
	if isSuspiciousPattern(trimmed) {
		return &ValidationError{
			Field:   "Name",
			Value:   name,
			Message: "Name appears to be invalid (repeating characters detected)",
		}
	}

	return nil
}

// isSuspiciousPattern checks if a string is suspicious (e.g., all same character)
func isSuspiciousPattern(s string) bool {
	if len(s) < 3 {
		return false
	}

	// Check if all characters are the same
	firstChar := s[0]
	allSame := true
	for i := 1; i < len(s); i++ {
		if s[i] != firstChar {
			allSame = false
			break
		}
	}

	if allSame {
		return true
	}

	// Check for patterns like AAAAAAAA (repeating single character)
	// More than 10 consecutive identical characters is suspicious
	count := 1
	maxCount := 1
	for i := 1; i < len(s); i++ {
		if s[i] == s[i-1] {
			count++
			if count > maxCount {
				maxCount = count
			}
		} else {
			count = 1
		}
	}

	// More than 10 consecutive identical characters is suspicious for names
	return maxCount > 10
}

// ValidateEmailFormat validates email format (basic validation)
func ValidateEmailFormat(email string) error {
	if email == "" {
		return nil // Email is optional
	}

	trimmed := strings.TrimSpace(email)
	if trimmed == "" {
		return nil // Empty email is allowed
	}

	// Basic email validation: must contain @ and .
	if !strings.Contains(trimmed, "@") || !strings.Contains(trimmed, ".") {
		return &ValidationError{
			Field:   "Email",
			Value:   email,
			Message: "Email format is invalid",
		}
	}

	// Check for suspicious patterns
	if isSuspiciousPattern(trimmed) {
		return &ValidationError{
			Field:   "Email",
			Value:   email,
			Message: "Email appears to be invalid",
		}
	}

	return nil
}

