package rimpay

import (
	"github.com/CatoSystems/rim-pay/internal/types"
)

// Re-export types from internal/types for public API
type ErrorCode = types.ErrorCode
type PaymentError = types.PaymentError

// Re-export constants
const (
	ErrorCodeInvalidRequest       = types.ErrorCodeInvalidRequest
	ErrorCodeAuthenticationFailed = types.ErrorCodeAuthenticationFailed
	ErrorCodeInsufficientFunds    = types.ErrorCodeInsufficientFunds
	ErrorCodePaymentDeclined      = types.ErrorCodePaymentDeclined
	ErrorCodeNetworkError         = types.ErrorCodeNetworkError
	ErrorCodeTimeout              = types.ErrorCodeTimeout
	ErrorCodeProviderError        = types.ErrorCodeProviderError
	ErrorCodeValidationError      = types.ErrorCodeValidationError
	ErrorCodePaymentExpired       = types.ErrorCodePaymentExpired
)

// Re-export constructor functions
var (
	NewPaymentError    = types.NewPaymentError
	NewValidationError = types.NewValidationError
	IsRetryableError   = types.IsRetryableError
)
