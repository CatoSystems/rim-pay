package masrvi

import (
	"net/url"
	"regexp"

	"github.com/CatoSystems/rim-pay/internal/types"
	"github.com/CatoSystems/rim-pay/pkg/money"
	"github.com/CatoSystems/rim-pay/pkg/phone"
)

// MasrviPaymentRequest represents a MASRVI specific payment request
type MasrviPaymentRequest struct {
	// Required fields
	Amount      money.Money  `json:"amount"`
	PhoneNumber *phone.Phone `json:"phone_number,omitempty"` // Optional for MASRVI
	Reference   string       `json:"reference"`

	// Optional fields
	Language    types.Language `json:"language,omitempty"`
	Description string         `json:"description,omitempty"`
	Brand       string         `json:"brand,omitempty"` // MASRVI specific: Payment brand

	// MASRVI specific URL fields
	SuccessURL  string `json:"success_url,omitempty"`  // Redirect on success
	FailureURL  string `json:"failure_url,omitempty"`  // Redirect on failure
	CancelURL   string `json:"cancel_url,omitempty"`   // Redirect on cancel
	CallbackURL string `json:"callback_url,omitempty"` // Webhook notification URL

	// Additional MASRVI fields
	CustomerName string `json:"customer_name,omitempty"` // Customer display name
	Text         string `json:"text,omitempty"`          // Additional text/instructions
}

// Validate validates the MASRVI payment request
func (r *MasrviPaymentRequest) Validate() error {
	if err := r.validateBasicFields(); err != nil {
		return err
	}

	if err := r.validateStringLengths(); err != nil {
		return err
	}

	return r.validateURLs()
}

// validateBasicFields validates amount and reference
func (r *MasrviPaymentRequest) validateBasicFields() error {
	if r.Amount.IsZero() {
		return types.NewValidationError("amount", "cannot be zero")
	}

	if r.Amount.IsNegative() {
		return types.NewValidationError("amount", "cannot be negative")
	}

	if r.Reference == "" {
		return types.NewValidationError("reference", "is required")
	}

	if len(r.Reference) > 50 {
		return types.NewValidationError("reference", "too long (max 50 characters)")
	}

	// Validate reference format (alphanumeric, dashes, underscores)
	validRefRegex := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	if !validRefRegex.MatchString(r.Reference) {
		return types.NewValidationError("reference", "invalid format (use only letters, numbers, dashes, underscores)")
	}

	return nil
}

// validateStringLengths validates length constraints
func (r *MasrviPaymentRequest) validateStringLengths() error {
	if len(r.Description) > 255 {
		return types.NewValidationError("description", "too long (max 255 characters)")
	}

	if len(r.CustomerName) > 100 {
		return types.NewValidationError("customer_name", "too long (max 100 characters)")
	}

	if len(r.Text) > 500 {
		return types.NewValidationError("text", "too long (max 500 characters)")
	}

	return nil
}

// validateURLs validates URL formats
func (r *MasrviPaymentRequest) validateURLs() error {
	const invalidURLMsg = "invalid URL format"

	urls := map[string]string{
		"success_url":  r.SuccessURL,
		"failure_url":  r.FailureURL,
		"cancel_url":   r.CancelURL,
		"callback_url": r.CallbackURL,
	}

	for field, url := range urls {
		if url != "" && !isValidURL(url) {
			return types.NewValidationError(field, invalidURLMsg)
		}
	}

	return nil
}

// GetLanguage returns the language with fallback to French
func (r *MasrviPaymentRequest) GetLanguage() types.Language {
	if r.Language == "" {
		return types.LanguageFrench
	}
	return r.Language
}

// ToGenericRequest converts to the internal generic payment request
func (r *MasrviPaymentRequest) ToGenericRequest() *types.PaymentRequest {
	return &types.PaymentRequest{
		Amount:      r.Amount,
		PhoneNumber: r.PhoneNumber,
		Reference:   r.Reference,
		Language:    r.GetLanguage(),
		Description: r.Description,
		SuccessURL:  r.SuccessURL,
		FailureURL:  r.FailureURL,
		CancelURL:   r.CancelURL,
		CallbackURL: r.CallbackURL,
	}
}

// isValidURL validates URL format
func isValidURL(urlStr string) bool {
	u, err := url.Parse(urlStr)
	return err == nil && u.Scheme != "" && u.Host != ""
}
