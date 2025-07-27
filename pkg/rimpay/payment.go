package rimpay

import (
	"github.com/CatoSystems/rim-pay/internal/types"
)

// Re-export types from internal/types for public API
type PaymentStatus = types.PaymentStatus
type Language = types.Language
type PaymentRequest = types.PaymentRequest
type PaymentResponse = types.PaymentResponse

// Re-export constants
const (
	PaymentStatusPending   = types.PaymentStatusPending
	PaymentStatusSuccess   = types.PaymentStatusSuccess
	PaymentStatusFailed    = types.PaymentStatusFailed
	PaymentStatusCancelled = types.PaymentStatusCancelled
	PaymentStatusExpired   = types.PaymentStatusExpired

	LanguageEnglish = types.LanguageEnglish
	LanguageFrench  = types.LanguageFrench
	LanguageArabic  = types.LanguageArabic
)
