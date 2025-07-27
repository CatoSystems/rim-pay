package types

import (
	"github.com/CatoSystems/rim-pay/pkg/money"
	"github.com/CatoSystems/rim-pay/pkg/phone"
	"time"
)

// PaymentStatus represents the status of a payment
type PaymentStatus string

const (
	// PaymentStatusPending indicates payment is pending
	PaymentStatusPending PaymentStatus = "pending"
	// PaymentStatusSuccess indicates payment was successful
	PaymentStatusSuccess PaymentStatus = "success"
	// PaymentStatusFailed indicates payment failed
	PaymentStatusFailed PaymentStatus = "failed"
	// PaymentStatusCancelled indicates payment was cancelled
	PaymentStatusCancelled PaymentStatus = "cancelled"
	// PaymentStatusExpired indicates payment expired
	PaymentStatusExpired PaymentStatus = "expired"
)

// Language represents supported languages
type Language string

const (
	// LanguageEnglish represents English
	LanguageEnglish Language = "EN"
	// LanguageFrench represents French
	LanguageFrench Language = "FR"
	// LanguageArabic represents Arabic
	LanguageArabic Language = "AR"
)

// PaymentRequest represents a payment request
type PaymentRequest struct {
	Amount      money.Money            `json:"amount"`
	PhoneNumber *phone.Phone           `json:"phone_number"`
	Reference   string                 `json:"reference"`
	Description string                 `json:"description,omitempty"`
	Language    Language               `json:"language,omitempty"`
	Passcode    string                 `json:"passcode,omitempty"`
	ExpiresAt   *time.Time             `json:"expires_at,omitempty"`
	SuccessURL  string                 `json:"success_url,omitempty"`
	FailureURL  string                 `json:"failure_url,omitempty"`
	CancelURL   string                 `json:"cancel_url,omitempty"`
	CallbackURL string                 `json:"callback_url,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// PaymentResponse represents a payment response
type PaymentResponse struct {
	TransactionID string                 `json:"transaction_id"`
	Status        PaymentStatus          `json:"status"`
	Amount        money.Money            `json:"amount"`
	Reference     string                 `json:"reference"`
	Provider      string                 `json:"provider"`
	CreatedAt     time.Time              `json:"created_at"`
	UpdatedAt     time.Time              `json:"updated_at"`
	PaymentURL    string                 `json:"payment_url,omitempty"`
	ExpiresAt     *time.Time             `json:"expires_at,omitempty"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
}

// IsSuccessful returns true if payment was successful
func (ps PaymentStatus) IsSuccessful() bool {
	return ps == PaymentStatusSuccess
}

// IsFailed returns true if payment failed
func (ps PaymentStatus) IsFailed() bool {
	return ps == PaymentStatusFailed
}

// IsCompleted returns true if payment is in a final state
func (ps PaymentStatus) IsCompleted() bool {
	return ps == PaymentStatusSuccess || ps == PaymentStatusFailed || ps == PaymentStatusCancelled || ps == PaymentStatusExpired
}

// String returns string representation
func (ps PaymentStatus) String() string {
	return string(ps)
}

// GetLanguage returns language with default fallback
func (pr *PaymentRequest) GetLanguage() Language {
	if pr.Language == "" {
		return LanguageFrench // Default for Mauritania
	}
	return pr.Language
}

// IsExpired returns true if payment request has expired
func (pr *PaymentRequest) IsExpired() bool {
	if pr.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*pr.ExpiresAt)
}

// IsCompleted returns true if payment is completed
func (pr *PaymentResponse) IsCompleted() bool {
	return pr.Status.IsCompleted()
}

// IsSuccessful returns true if payment was successful
func (pr *PaymentResponse) IsSuccessful() bool {
	return pr.Status.IsSuccessful()
}

// Validate validates payment request
func (pr *PaymentRequest) Validate() error {
	if pr.Amount.IsZero() || pr.Amount.IsNegative() {
		return NewValidationError("amount", "must be positive")
	}

	if err := pr.Amount.Validate(); err != nil {
		return NewValidationError("amount", err.Error())
	}

	if pr.PhoneNumber == nil {
		return NewValidationError("phone_number", "is required")
	}

	if pr.Reference == "" {
		return NewValidationError("reference", "is required")
	}

	if len(pr.Reference) > 50 {
		return NewValidationError("reference", "too long (max 50 characters)")
	}

	return nil
}
