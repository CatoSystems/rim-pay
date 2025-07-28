package rimpay

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/CatoSystems/rim-pay/pkg/money"
	"github.com/CatoSystems/rim-pay/pkg/phone"
)

// BPayPaymentRequest represents a B-PAY specific payment request
type BPayPaymentRequest struct {
	PhoneNumber *phone.Phone           `json:"phone_number"`
	Amount      money.Money            `json:"amount"`
	Description string                 `json:"description"`
	Reference   string                 `json:"reference"`
	Passcode    string                 `json:"passcode"` // B-PAY specific: user passcode
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// Validate validates the B-PAY payment request
func (r *BPayPaymentRequest) Validate() error {
	if r.PhoneNumber == nil {
		return fmt.Errorf("phone number is required")
	}

	if r.Amount.IsZero() {
		return fmt.Errorf("amount must be positive")
	}

	if strings.TrimSpace(r.Description) == "" {
		return fmt.Errorf("description cannot be empty")
	}

	if strings.TrimSpace(r.Reference) == "" {
		return fmt.Errorf("reference cannot be empty")
	}

	if len(r.Reference) > 50 {
		return fmt.Errorf("reference cannot exceed 50 characters")
	}

	// Note: Passcode validation is not needed as the library always generates a new passcode

	return nil
}

// ToGenericRequest converts B-PAY request to generic payment request
func (r *BPayPaymentRequest) ToGenericRequest() *PaymentRequest {
	metadata := make(map[string]interface{})
	for k, v := range r.Metadata {
		metadata[k] = v
	}
	// Note: Passcode is not included as the library will always generate a new one

	return &PaymentRequest{
		PhoneNumber: r.PhoneNumber,
		Amount:      r.Amount,
		Description: r.Description,
		Reference:   r.Reference,
		// Passcode is intentionally empty - library will generate a new one
		Metadata:    metadata,
	}
}

// MasrviPaymentRequest represents a MASRVI specific payment request
type MasrviPaymentRequest struct {
	PhoneNumber *phone.Phone           `json:"phone_number"`
	Amount      money.Money            `json:"amount"`
	Description string                 `json:"description"`
	Reference   string                 `json:"reference"`
	CallbackURL string                 `json:"callback_url"` // MASRVI specific: webhook URL
	ReturnURL   string                 `json:"return_url"`   // MASRVI specific: return URL after payment
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// Validate validates the MASRVI payment request
func (r *MasrviPaymentRequest) Validate() error {
	if err := r.validateBasicFields(); err != nil {
		return err
	}

	if err := r.validateStringLengths(); err != nil {
		return err
	}

	if err := r.validateURLs(); err != nil {
		return err
	}

	return nil
}

func (r *MasrviPaymentRequest) validateBasicFields() error {
	if r.PhoneNumber == nil {
		return fmt.Errorf("phone number is required")
	}

	if r.Amount.IsZero() {
		return fmt.Errorf("amount must be positive")
	}

	return nil
}

func (r *MasrviPaymentRequest) validateStringLengths() error {
	const (
		maxDescriptionLength = 200
		maxReferenceLength   = 50
	)

	if strings.TrimSpace(r.Description) == "" {
		return fmt.Errorf("description cannot be empty")
	}

	if len(r.Description) > maxDescriptionLength {
		return fmt.Errorf("description cannot exceed %d characters", maxDescriptionLength)
	}

	if strings.TrimSpace(r.Reference) == "" {
		return fmt.Errorf("reference cannot be empty")
	}

	if len(r.Reference) > maxReferenceLength {
		return fmt.Errorf("reference cannot exceed %d characters", maxReferenceLength)
	}

	return nil
}

func (r *MasrviPaymentRequest) validateURLs() error {
	if strings.TrimSpace(r.CallbackURL) == "" {
		return fmt.Errorf("callback_url cannot be empty")
	}

	if _, err := url.Parse(r.CallbackURL); err != nil {
		return fmt.Errorf("invalid callback_url: %w", err)
	}

	if strings.TrimSpace(r.ReturnURL) == "" {
		return fmt.Errorf("return_url cannot be empty")
	}

	if _, err := url.Parse(r.ReturnURL); err != nil {
		return fmt.Errorf("invalid return_url: %w", err)
	}

	return nil
}

// ToGenericRequest converts MASRVI request to generic payment request
func (r *MasrviPaymentRequest) ToGenericRequest() *PaymentRequest {
	metadata := make(map[string]interface{})
	for k, v := range r.Metadata {
		metadata[k] = v
	}
	metadata["callback_url"] = r.CallbackURL
	metadata["return_url"] = r.ReturnURL

	return &PaymentRequest{
		PhoneNumber: r.PhoneNumber,
		Amount:      r.Amount,
		Description: r.Description,
		Reference:   r.Reference,
		Metadata:    metadata,
	}
}

// MasrviNotificationData represents MASRVI webhook notification
type MasrviNotificationData struct {
	TransactionID string                 `json:"transaction_id"`
	Status        string                 `json:"status"`
	Reference     string                 `json:"reference"`
	Amount        money.Money            `json:"amount"`
	PhoneNumber   string                 `json:"phone_number"`
	Timestamp     string                 `json:"timestamp"`
	Data          map[string]interface{} `json:"data,omitempty"`
}
