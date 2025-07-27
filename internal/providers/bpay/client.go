package bpay

import (
	"context"
	"fmt"
	_ "strings"

	"github.com/CatoSystems/rim-pay/internal/providers/common"
	"github.com/CatoSystems/rim-pay/pkg/rimpay"
)

// Provider implements the B-PAY payment provider
type Provider struct {
	name             string
	config           rimpay.ProviderConfig
	httpClient       common.HTTPClient
	authManager      *AuthManager
	paymentProcessor *PaymentProcessor
	retryExecutor    *common.RetryExecutor
	logger           rimpay.Logger
}

// NewProvider creates a new B-PAY provider
func NewProvider(config rimpay.ProviderConfig, logger rimpay.Logger) (*Provider, error) {
	if err := validateConfig(config); err != nil {
		return nil, fmt.Errorf("invalid B-PAY configuration: %w", err)
	}

	// Create HTTP client
	httpClient := common.NewHTTPClient(common.HTTPConfig{
		Timeout:         config.Timeout,
		MaxIdleConns:    10,
		MaxConnsPerHost: 5,
	})

	// Create authentication manager
	authManager := NewAuthManager(config, httpClient, logger)

	// Create payment processor
	paymentProcessor := NewPaymentProcessor(config, httpClient, authManager, logger)

	// Create retry executor with default config
	retryExecutor := common.NewRetryExecutor(common.DefaultRetryConfig())

	provider := &Provider{
		name:             "bpay",
		config:           config,
		httpClient:       httpClient,
		authManager:      authManager,
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
	_, err := p.authManager.GetAccessToken(ctx)
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
func (p *Provider) GetPaymentStatus(ctx context.Context, transactionID string) (*rimpay.TransactionStatus, error) {
	return p.paymentProcessor.CheckPaymentStatus(ctx, transactionID)
}

// ValidateConfig validates provider configuration
func (p *Provider) ValidateConfig() error {
	return validateConfig(p.config)
}

// validateConfig validates B-PAY configuration
func validateConfig(config rimpay.ProviderConfig) error {
	requiredCredentials := []string{"username", "password", "client_id"}

	for _, field := range requiredCredentials {
		if config.Credentials[field] == "" {
			return fmt.Errorf("missing required credential: %s", field)
		}
	}

	if config.BaseURL == "" {
		return fmt.Errorf("base_url is required")
	}

	if config.Timeout <= 0 {
		return fmt.Errorf("timeout must be positive")
	}

	return nil
}
