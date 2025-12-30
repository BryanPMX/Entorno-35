package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateCURPFormat_Valid(t *testing.T) {
	tests := []struct {
		name string
		curp string
	}{
		{"Valid CURP 1", "ABCD123456HIJKLM01"},
		{"Valid CURP 2", "EFGH567890MNOPQR02"},
		{"Valid CURP 3", "IJKL901234RSTUVW03"},
		{"Valid CURP lowercase", "abcd123456hijklm01"}, // Should be normalized
		{"Valid CURP with spaces", "  ABCD123456HIJKLM01  "}, // Should be trimmed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCURPFormat(tt.curp)
			assert.NoError(t, err)
		})
	}
}

func TestValidateCURPFormat_Invalid(t *testing.T) {
	tests := []struct {
		name    string
		curp    string
		wantErr string
	}{
		{"Too short", "ABCD123456HIJKLM0", "exactly 18 characters"},
		{"Too long", "ABCD123456HIJKLM012", "exactly 18 characters"},
		{"Wrong format - missing digits", "ABCDEFGHIJKLMNOPQR", "format is invalid"},
		{"Wrong format - missing letters", "123456789012345678", "format is invalid"},
		{"Wrong format - mixed position", "1234ABCD5678EFGH01", "format is invalid"},
		{"Special characters", "ABCD-23456HIJKLM01", "invalid characters"},
		{"All same character", "AAAAAAAAAAAAAAAAAA", "format is invalid"}, // Caught by regex (must start with 4 letters, but pattern still valid, caught by suspicious pattern check in name validation if used there)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCURPFormat(tt.curp)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErr)
		})
	}
}

func TestValidateName_Valid(t *testing.T) {
	tests := []struct {
		name string
		val  string
	}{
		{"Normal name", "Juan Pérez García"},
		{"Name with accents", "José María"},
		{"Short name", "Ana"},
		{"Long name", "María de los Ángeles González López"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateName(tt.val)
			assert.NoError(t, err)
		})
	}
}

func TestValidateName_Invalid(t *testing.T) {
	tests := []struct {
		name    string
		val     string
		wantErr string
	}{
		{"Empty", "", "required and cannot be empty"},
		{"Whitespace only", "   ", "required and cannot be empty"},
		{"All same character (suspicious)", "AAAAAAAAAAAAAAAAAA", "repeating characters"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateName(tt.val)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErr)
		})
	}
}

func TestValidateEmailFormat_Valid(t *testing.T) {
	tests := []struct {
		name  string
		email string
	}{
		{"Valid email", "user@example.com"},
		{"Valid email with subdomain", "user@mail.example.com"},
		{"Empty (optional)", ""},
		{"Whitespace only (treated as empty)", "   "},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateEmailFormat(tt.email)
			assert.NoError(t, err)
		})
	}
}

func TestValidateEmailFormat_Invalid(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		wantErr string
	}{
		{"No @ symbol", "userexample.com", "format is invalid"},
		{"No . symbol", "user@examplecom", "format is invalid"},
		{"Suspicious pattern", "aaaaaaaaaaaaaaaa@aaaa.aaa", "appears to be invalid"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateEmailFormat(tt.email)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErr)
		})
	}
}

func TestIsSuspiciousPattern(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"All same character", "AAAAAAAAAAAAAAA", true},
		{"Normal text", "Juan Pérez", false},
		{"Short text", "AB", false},
		{"Many same consecutive", "AAAABBBBCCCCDDDD", false}, // Different characters
		{"11+ consecutive same", "AAAAAAAABBBBBBBBBBBB", true}, // More than 10 consecutive
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isSuspiciousPattern(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

