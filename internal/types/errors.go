package types

import "fmt"

// ErrorCode represents specific error codes
type ErrorCode string

const (
	// ErrorCodeInvalidRequest indicates invalid request
	ErrorCodeInvalidRequest ErrorCode = "INVALID_REQUEST"
	// ErrorCodeAuthenticationFailed indicates authentication failure
	ErrorCodeAuthenticationFailed ErrorCode = "AUTHENTICATION_FAILED"
	// ErrorCodeInsufficientFunds indicates insufficient funds
	ErrorCodeInsufficientFunds ErrorCode = "INSUFFICIENT_FUNDS"
	// ErrorCodePaymentDeclined indicates payment declined
	ErrorCodePaymentDeclined ErrorCode = "PAYMENT_DECLINED"
	// ErrorCodeNetworkError indicates network error
	ErrorCodeNetworkError ErrorCode = "NETWORK_ERROR"
	// ErrorCodeTimeout indicates timeout
	ErrorCodeTimeout ErrorCode = "TIMEOUT"
	// ErrorCodeProviderError indicates provider error
	ErrorCodeProviderError ErrorCode = "PROVIDER_ERROR"
	// ErrorCodeValidationError indicates validation error
	ErrorCodeValidationError ErrorCode = "VALIDATION_ERROR"
	// ErrorCodePaymentExpired indicates payment expired
	ErrorCodePaymentExpired ErrorCode = "PAYMENT_EXPIRED"
)

// PaymentError represents a payment-related error
type PaymentError struct {
	Code      ErrorCode              `json:"code"`
	Message   string                 `json:"message"`
	Provider  string                 `json:"provider,omitempty"`
	Details   map[string]interface{} `json:"details,omitempty"`
	Retryable bool                   `json:"retryable"`
	Cause     error                  `json:"-"`
}

// Error implements the error interface
func (e *PaymentError) Error() string {
	if e.Provider != "" {
		return fmt.Sprintf("[%s] %s: %s", e.Provider, e.Code, e.Message)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap returns the underlying error
func (e *PaymentError) Unwrap() error {
	return e.Cause
}

// IsRetryable returns whether the error is retryable
func (e *PaymentError) IsRetryable() bool {
	return e.Retryable
}

// NewPaymentError creates a new payment error
func NewPaymentError(code ErrorCode, message string, provider string, retryable bool) *PaymentError {
	return &PaymentError{
		Code:      code,
		Message:   message,
		Provider:  provider,
		Retryable: retryable,
		Details:   make(map[string]interface{}),
	}
}

// WithCause adds a cause to the payment error
func (e *PaymentError) WithCause(cause error) *PaymentError {
	e.Cause = cause
	return e
}

// WithDetail adds a detail to the payment error
func (e *PaymentError) WithDetail(key string, value interface{}) *PaymentError {
	if e.Details == nil {
		e.Details = make(map[string]interface{})
	}
	e.Details[key] = value
	return e
}

// NewValidationError creates a validation error
func NewValidationError(field, message string) *PaymentError {
	return &PaymentError{
		Code:    ErrorCodeValidationError,
		Message: fmt.Sprintf("%s: %s", field, message),
		Details: map[string]interface{}{"field": field},
	}
}

// IsRetryableError determines if error code is retryable
func IsRetryableError(code ErrorCode) bool {
	retryableCodes := map[ErrorCode]bool{
		ErrorCodeNetworkError:  true,
		ErrorCodeTimeout:       true,
		ErrorCodeProviderError: true,
	}
	return retryableCodes[code]
}
