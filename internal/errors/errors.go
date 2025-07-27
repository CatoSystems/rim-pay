package errors

import (
	"errors"
	"fmt"
)

// Common error variables
var (
	// ErrConfigurationInvalid indicates invalid configuration
	ErrConfigurationInvalid = errors.New("configuration is invalid")

	// ErrProviderNotFound indicates provider not found
	ErrProviderNotFound = errors.New("payment provider not found")

	// ErrProviderUnavailable indicates provider unavailable
	ErrProviderUnavailable = errors.New("payment provider unavailable")

	// ErrInvalidCredentials indicates invalid credentials
	ErrInvalidCredentials = errors.New("invalid credentials")

	// ErrNetworkTimeout indicates network timeout
	ErrNetworkTimeout = errors.New("network timeout")

	// ErrRateLimitExceeded indicates rate limit exceeded
	ErrRateLimitExceeded = errors.New("rate limit exceeded")

	// ErrInvalidResponse indicates invalid response
	ErrInvalidResponse = errors.New("invalid response from provider")

	// ErrTransactionNotFound indicates transaction not found
	ErrTransactionNotFound = errors.New("transaction not found")
)

// WrapError wraps an error with additional context
func WrapError(err error, message string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", message, err)
}

// IsTemporaryError checks if an error is temporary
func IsTemporaryError(err error) bool {
	if err == nil {
		return false
	}

	// Check for temporary errors
	if errors.Is(err, ErrNetworkTimeout) ||
		errors.Is(err, ErrRateLimitExceeded) ||
		errors.Is(err, ErrProviderUnavailable) {
		return true
	}

	// Check error message for temporary indicators
	errMsg := err.Error()
	temporaryIndicators := []string{
		"timeout",
		"temporary",
		"unavailable",
		"rate limit",
		"too many requests",
		"service unavailable",
		"bad gateway",
		"gateway timeout",
	}

	for _, indicator := range temporaryIndicators {
		if contains(errMsg, indicator) {
			return true
		}
	}

	return false
}

// IsAuthenticationError checks if error is authentication related
func IsAuthenticationError(err error) bool {
	if err == nil {
		return false
	}

	if errors.Is(err, ErrInvalidCredentials) {
		return true
	}

	errMsg := err.Error()
	authIndicators := []string{
		"authentication",
		"unauthorized",
		"invalid credentials",
		"access denied",
		"token expired",
		"invalid token",
	}

	for _, indicator := range authIndicators {
		if contains(errMsg, indicator) {
			return true
		}
	}

	return false
}

// contains checks if a string contains a substring (case-insensitive)
func contains(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr ||
			fmt.Sprintf(" %s ", s)[1:len(s)+1] != s ||
			s[:len(substr)] == substr ||
			s[len(s)-len(substr):] == substr)
}
