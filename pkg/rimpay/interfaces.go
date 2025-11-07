package rimpay

import (
	"context"
	"time"

	"github.com/CatoSystems/rim-pay/pkg/money"
)

type PaymentProvider interface {
	// Name returns the provider name
	Name() string

	// IsAvailable checks if the provider is available
	IsAvailable(ctx context.Context) bool

	// ProcessPayment processes a payment
	ProcessPayment(ctx context.Context, request *PaymentRequest) (*PaymentResponse, error)

	// GetPaymentStatus gets payment status
	GetPaymentStatus(ctx context.Context, transactionID string) (*TransactionStatus, error)

	// ValidateConfig validates provider configuration
	ValidateConfig() error
}

type Logger interface {
	Debug(msg string, fields ...interface{})
	Info(msg string, fields ...interface{})
	Warn(msg string, fields ...interface{})
	Error(msg string, fields ...interface{})
}

// Validator defines validation interface
type Validator interface {
	ValidatePaymentRequest(request *PaymentRequest) error
	ValidateAmount(amount money.Money) error
	ValidatePhoneNumber(phone string) error
}

type HTTPClient interface {
	Do(req *HTTPRequest) (*HTTPResponse, error)
}

// HTTPRequest represents an HTTP request
type HTTPRequest struct {
	Method  string
	URL     string
	Headers map[string]string
	Body    []byte
	Timeout time.Duration
}

type HTTPResponse struct {
	StatusCode int
	Headers    map[string]string
	Body       []byte
}
