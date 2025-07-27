package validation

import (
	"github.com/CatoSystems/rim-pay/internal/types"
	"github.com/CatoSystems/rim-pay/pkg/money"
	"github.com/CatoSystems/rim-pay/pkg/phone"
	"regexp"
)

type Validator struct {
	emailRegex *regexp.Regexp
	urlRegex   *regexp.Regexp
}

// NewValidator creates a new validator
func NewValidator() *Validator {
	return &Validator{
		emailRegex: regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`),
		urlRegex:   regexp.MustCompile(`^[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`),
	}
}

// ValidatePaymentRequest validates a payment request
func (v *Validator) ValidatePaymentRequest(request *types.PaymentRequest) error {
	if request == nil {
		return types.NewValidationError("request", "cannot be nil")
	}

	// Validate amount
	if err := v.ValidateAmount(request.Amount); err != nil {
		return err
	}

	// Validate phone number
	if request.PhoneNumber == nil {
		return types.NewValidationError("phone_number", "is required")
	}

	// Validate reference
	if err := v.validateReference(request.Reference); err != nil {
		return err
	}

	// Validate URLs if provided
	if request.SuccessURL != "" && !v.isValidURL(request.SuccessURL) {
		return types.NewValidationError("success_url", "invalid URL format")
	}

	if request.FailureURL != "" && !v.isValidURL(request.FailureURL) {
		return types.NewValidationError("failure_url", "invalid URL format")
	}

	if request.CancelURL != "" && !v.isValidURL(request.CancelURL) {
		return types.NewValidationError("cancel_url", "invalid URL format")
	}

	if request.CallbackURL != "" && !v.isValidURL(request.CallbackURL) {
		return types.NewValidationError("callback_url", "invalid URL format")
	}

	// Validate description length
	if len(request.Description) > 255 {
		return types.NewValidationError("description", "too long (max 255 characters)")
	}

	return nil
}

// ValidateAmount validates a monetary amount
func (v *Validator) ValidateAmount(amount money.Money) error {
	if amount.IsZero() {
		return types.NewValidationError("amount", "cannot be zero")
	}

	if amount.IsNegative() {
		return types.NewValidationError("amount", "cannot be negative")
	}

	// Validate currency
	if err := amount.Validate(); err != nil {
		return types.NewValidationError("amount", err.Error())
	}

	// Check reasonable limits (adjust based on business requirements)
	if amount.Float64() > 10000000 { // 10 million
		return types.NewValidationError("amount", "exceeds maximum allowed amount")
	}

	return nil
}

// ValidatePhoneNumber validates a phone number string
func (v *Validator) ValidatePhoneNumber(phoneStr string) error {
	if phoneStr == "" {
		return types.NewValidationError("phone_number", "cannot be empty")
	}

	// Try to parse as Mauritanian phone number
	_, err := phone.NewPhone(phoneStr)
	if err != nil {
		return types.NewValidationError("phone_number", err.Error())
	}

	return nil
}

// validateReference validates payment reference
func (v *Validator) validateReference(reference string) error {
	if reference == "" {
		return types.NewValidationError("reference", "is required")
	}

	if len(reference) > 50 {
		return types.NewValidationError("reference", "too long (max 50 characters)")
	}

	// Check for valid characters (alphanumeric, dashes, underscores)
	validRefRegex := regexp.MustCompile(`^[a-zA-Z0-9_-]+`)
	if !validRefRegex.MatchString(reference) {
		return types.NewValidationError("reference", "invalid reference format")
	}
	return nil
}

// isValidURL validates URL format
func (v *Validator) isValidURL(url string) bool {
	return v.urlRegex.MatchString(url)
}
