package bpay

import (
	"github.com/CatoSystems/rim-pay/internal/types"
	"github.com/CatoSystems/rim-pay/pkg/money"
	"github.com/CatoSystems/rim-pay/pkg/phone"
)

// BPayPaymentRequest represents a B-PAY specific payment request
type BPayPaymentRequest struct {
	// Required fields
	Amount      money.Money  `json:"amount"`
	PhoneNumber *phone.Phone `json:"phone_number"`
	Reference   string       `json:"reference"`
	Passcode    string       `json:"passcode"` // B-PAY specific: Mobile money PIN

	// Optional fields
	Language    types.Language `json:"language,omitempty"`
	Description string         `json:"description,omitempty"`
}

// Validate validates the B-PAY payment request
func (r *BPayPaymentRequest) Validate() error {
	if r.Amount.IsZero() {
		return types.NewValidationError("amount", "cannot be zero")
	}

	if r.Amount.IsNegative() {
		return types.NewValidationError("amount", "cannot be negative")
	}

	if r.PhoneNumber == nil {
		return types.NewValidationError("phone_number", "is required")
	}

	if r.Reference == "" {
		return types.NewValidationError("reference", "is required")
	}

	if r.Passcode == "" {
		return types.NewValidationError("passcode", "is required for B-PAY payments")
	}

	if len(r.Passcode) < 4 || len(r.Passcode) > 8 {
		return types.NewValidationError("passcode", "must be between 4 and 8 digits")
	}

	if len(r.Description) > 255 {
		return types.NewValidationError("description", "too long (max 255 characters)")
	}

	return nil
}

// GetLanguage returns the language with fallback to French
func (r *BPayPaymentRequest) GetLanguage() types.Language {
	if r.Language == "" {
		return types.LanguageFrench
	}
	return r.Language
}

// ToGenericRequest converts to the internal generic payment request
func (r *BPayPaymentRequest) ToGenericRequest() *types.PaymentRequest {
	return &types.PaymentRequest{
		Amount:      r.Amount,
		PhoneNumber: r.PhoneNumber,
		Reference:   r.Reference,
		Language:    r.GetLanguage(),
		Description: r.Description,
		Passcode:    r.Passcode,
	}
}
