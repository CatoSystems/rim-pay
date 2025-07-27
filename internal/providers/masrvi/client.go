package masrvi

import (
	"context"
	"fmt"
	"github.com/CatoSystems/rim-pay/internal/providers/common"
	"github.com/CatoSystems/rim-pay/pkg/rimpay"
)

// Provider implements the MASRVI payment provider
type Provider struct {
	name             string
	config           rimpay.ProviderConfig
	httpClient       common.HTTPClient
	sessionManager   *SessionManager
	paymentProcessor *PaymentProcessor
	retryExecutor    *common.RetryExecutor
	logger           rimpay.Logger
}

// NewProvider creates a new MASRVI provider
func NewProvider(config rimpay.ProviderConfig, logger rimpay.Logger) (*Provider, error) {
	if err := validateConfig(config); err != nil {
		return nil, fmt.Errorf("invalid MASRVI configuration: %w", err)
	}

	// Create HTTP client
	httpClient := common.NewHTTPClient(common.HTTPConfig{
		Timeout:         config.Timeout,
		MaxIdleConns:    10,
		MaxConnsPerHost: 5,
	})

	// Create session manager
	sessionManager := NewSessionManager(config, httpClient, logger)

	// Create payment processor
	paymentProcessor := NewPaymentProcessor(config, httpClient, sessionManager, logger)

	// Create retry executor with default config
	retryExecutor := common.NewRetryExecutor(common.DefaultRetryConfig())

	provider := &Provider{
		name:             "masrvi",
		config:           config,
		httpClient:       httpClient,
		sessionManager:   sessionManager,
		paymentProcessor: paymentProcessor,
		retryExecutor:    retryExecutor,
		logger:           logger,
	}

	return provider, nil
}

// Name returns the provider name
func (p *Provider) Name() string {
	return p.name
}

// IsAvailable checks if the provider is available
func (p *Provider) IsAvailable(ctx context.Context) bool {
	_, err := p.sessionManager.GetSessionID(ctx)
	return err == nil
}

// ProcessPayment processes a payment with retry logic
func (p *Provider) ProcessPayment(ctx context.Context, request *rimpay.PaymentRequest) (*rimpay.PaymentResponse, error) {
	// Wrap the payment processing in a retryable function
	retryablePayment := func() (*rimpay.PaymentResponse, error) {
		return p.paymentProcessor.ProcessPayment(ctx, request)
	}

	// Execute with retry logic
	return p.retryExecutor.ExecutePayment(ctx, retryablePayment)
}

// GetPaymentStatus gets payment status
// Note: MASRVI doesn't have a direct status check API, status comes via webhooks
func (p *Provider) GetPaymentStatus(ctx context.Context, transactionID string) (*rimpay.TransactionStatus, error) {
	return &rimpay.TransactionStatus{
		TransactionID: transactionID,
		Status:        rimpay.PaymentStatusPending,
		Reference:     transactionID,
		Message:       "Status check not supported, use webhook notifications",
		ProviderData: map[string]interface{}{
			"note": "MASRVI uses webhook notifications for status updates",
		},
	}, nil
}

// HandleNotification handles webhook notifications
func (p *Provider) HandleNotification(notification *NotificationData) (*rimpay.TransactionStatus, error) {
	return p.paymentProcessor.HandleNotification(notification)
}

// ValidateConfig validates provider configuration
func (p *Provider) ValidateConfig() error {
	return validateConfig(p.config)
}

// validateConfig validates MASRVI configuration
func validateConfig(config rimpay.ProviderConfig) error {
	if config.Credentials["merchant_id"] == "" {
		return fmt.Errorf("missing required credential: merchant_id")
	}

	if config.BaseURL == "" {
		return fmt.Errorf("base_url is required")
	}

	if config.Timeout <= 0 {
		return fmt.Errorf("timeout must be positive")
	}

	return nil
}
