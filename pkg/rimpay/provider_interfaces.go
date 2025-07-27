package rimpay

import (
	"context"
)

// BPayProvider represents the B-PAY payment provider interface
type BPayProvider interface {
	// Name returns the provider name
	Name() string

	// IsAvailable checks if the provider is available
	IsAvailable(ctx context.Context) bool

	// ProcessBPayPayment processes a B-PAY payment
	ProcessBPayPayment(ctx context.Context, request *BPayPaymentRequest) (*PaymentResponse, error)

	// GetPaymentStatus gets payment status
	GetPaymentStatus(ctx context.Context, transactionID string) (*TransactionStatus, error)

	// ValidateConfig validates provider configuration
	ValidateConfig() error
}

// MasrviProvider represents the MASRVI payment provider interface
type MasrviProvider interface {
	// Name returns the provider name
	Name() string

	// IsAvailable checks if the provider is available
	IsAvailable(ctx context.Context) bool

	// ProcessPayment processes a MASRVI payment
	ProcessMasrviPayment(ctx context.Context, request *MasrviPaymentRequest) (*PaymentResponse, error)

	// GetPaymentStatus gets payment status (note: MASRVI uses webhooks)
	GetPaymentStatus(ctx context.Context, transactionID string) (*TransactionStatus, error)

	// HandleNotification handles MASRVI webhook notifications
	HandleNotification(notification *MasrviNotificationData) (*TransactionStatus, error)

	// ValidateConfig validates provider configuration
	ValidateConfig() error
}
