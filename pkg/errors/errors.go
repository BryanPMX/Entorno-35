package errors

import (
	"errors"
	"net/http"
)

// ErrorCode represents application error codes
type ErrorCode string

const (
	ErrCodeInternal      ErrorCode = "INTERNAL_ERROR"
	ErrCodeNotFound      ErrorCode = "NOT_FOUND"
	ErrCodeBadRequest    ErrorCode = "BAD_REQUEST"
	ErrCodeUnauthorized  ErrorCode = "UNAUTHORIZED"
	ErrCodeForbidden     ErrorCode = "FORBIDDEN"
	ErrCodeConflict      ErrorCode = "CONFLICT"
	ErrCodeValidation    ErrorCode = "VALIDATION_ERROR"
)

// AppError represents an application error
type AppError struct {
	Code       ErrorCode `json:"code"`
	Message    string    `json:"message"`
	StatusCode int       `json:"-"`
	Err        error     `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Message
}

// NewAppError creates a new application error
func NewAppError(code ErrorCode, message string, statusCode int) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		StatusCode: statusCode,
	}
}

// WrapError wraps an existing error
func WrapError(code ErrorCode, message string, statusCode int, err error) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		StatusCode: statusCode,
		Err:        err,
	}
}

// Predefined errors
var (
	ErrNotFound     = NewAppError(ErrCodeNotFound, "Resource not found", http.StatusNotFound)
	ErrBadRequest   = NewAppError(ErrCodeBadRequest, "Bad request", http.StatusBadRequest)
	ErrUnauthorized = NewAppError(ErrCodeUnauthorized, "Unauthorized", http.StatusUnauthorized)
	ErrForbidden    = NewAppError(ErrCodeForbidden, "Forbidden", http.StatusForbidden)
	ErrInternal     = NewAppError(ErrCodeInternal, "Internal server error", http.StatusInternalServerError)
)

// IsAppError checks if error is an AppError
func IsAppError(err error) bool {
	var appErr *AppError
	return errors.As(err, &appErr)
}

// AsAppError converts error to AppError
func AsAppError(err error) (*AppError, bool) {
	var appErr *AppError
	return appErr, errors.As(err, &appErr)
}

