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

	if strings.TrimSpace(r.Passcode) == "" {
		return fmt.Errorf("passcode is required (the customer's Bankily verification code)")
	}

	if !isFourDigitPasscode(r.Passcode) {
		return fmt.Errorf("passcode must be exactly 4 digits")
	}

	return nil
}

// isFourDigitPasscode reports whether s is exactly four ASCII digits, matching
// the Bankily B-PAY passcode format.
func isFourDigitPasscode(s string) bool {
	if len(s) != 4 {
		return false
	}
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

// ToGenericRequest converts B-PAY request to generic payment request
func (r *BPayPaymentRequest) ToGenericRequest() *PaymentRequest {
	metadata := make(map[string]interface{})
	for k, v := range r.Metadata {
		metadata[k] = v
	}

	return &PaymentRequest{
		PhoneNumber: r.PhoneNumber,
		Amount:      r.Amount,
		Description: r.Description,
		Reference:   r.Reference,
		Passcode:    r.Passcode,
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

// ClickPaymentRequest represents a CLICK (TagPay/BNM) specific payment request.
type ClickPaymentRequest struct {
	PhoneNumber *phone.Phone           `json:"phone_number,omitempty"` // optional
	Amount      money.Money            `json:"amount"`
	Reference   string                 `json:"reference"`
	Description string                 `json:"description,omitempty"`
	SuccessURL  string                 `json:"success_url,omitempty"`
	FailureURL  string                 `json:"failure_url,omitempty"`
	CancelURL   string                 `json:"cancel_url,omitempty"`
	Brand       string                 `json:"brand,omitempty"`
	Language    Language               `json:"language,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// Validate validates the CLICK payment request.
func (r *ClickPaymentRequest) Validate() error {
	if r.Amount.IsZero() {
		return fmt.Errorf("amount must be positive")
	}
	if strings.TrimSpace(r.Reference) == "" {
		return fmt.Errorf("reference cannot be empty")
	}
	if len(r.Reference) > 250 {
		return fmt.Errorf("reference cannot exceed 250 characters")
	}
	if len(r.Description) > 255 {
		return fmt.Errorf("description cannot exceed 255 characters")
	}
	return nil
}

// GetLanguage returns the language with fallback to French.
func (r *ClickPaymentRequest) GetLanguage() Language {
	if r.Language == "" {
		return LanguageFrench
	}
	return r.Language
}

// ToGenericRequest converts CLICK request to the generic payment request.
func (r *ClickPaymentRequest) ToGenericRequest() *PaymentRequest {
	metadata := make(map[string]interface{})
	for k, v := range r.Metadata {
		metadata[k] = v
	}
	if r.Brand != "" {
		metadata["brand"] = r.Brand
	}
	return &PaymentRequest{
		PhoneNumber: r.PhoneNumber,
		Amount:      r.Amount,
		Reference:   r.Reference,
		Description: r.Description,
		Language:    r.GetLanguage(),
		SuccessURL:  r.SuccessURL,
		FailureURL:  r.FailureURL,
		CancelURL:   r.CancelURL,
		Metadata:    metadata,
	}
}

// ClickNotificationData is the public shape of a TagPay server-to-server
// notification (GET query parameters) for the CLICK provider.
type ClickNotificationData struct {
	Status      string `json:"status"` // OK / NOK
	PurchaseRef string `json:"purchaseref"`
	Amount      string `json:"amount"`
	Currency    string `json:"currency"`
	ClientID    string `json:"clientid"`
	ClientName  string `json:"cname"`
	Mobile      string `json:"mobile"`
	PaymentRef  string `json:"paymentref"`
	PayID       string `json:"payid"`
	Timestamp   string `json:"timestamp"`
	IPAddress   string `json:"ipaddr"`
	Error       string `json:"error,omitempty"`
	Reason      string `json:"reason,omitempty"`
}
