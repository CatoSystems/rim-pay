package masrvi

import (
	"context"
	"fmt"
	"time"

	"github.com/CatoSystems/rim-pay/internal/providers/common"
	"github.com/CatoSystems/rim-pay/internal/types"
	"github.com/CatoSystems/rim-pay/pkg/rimpay"
)

// Register the MASRVI provider with the client
func init() {
	rimpay.RegisterMasrviProvider(func(config rimpay.ProviderConfig, logger rimpay.Logger) (rimpay.PaymentProvider, error) {
		return NewMasrviProvider(config, logger)
	})
}

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

// NewProvider creates a new MASRVI provider (deprecated, use NewMasrviProvider)
func NewProvider(config rimpay.ProviderConfig, logger rimpay.Logger) (*Provider, error) {
	return NewMasrviProvider(config, logger)
}

// NewMasrviProvider creates a new MASRVI provider
func NewMasrviProvider(config rimpay.ProviderConfig, logger rimpay.Logger) (*Provider, error) {
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

// ProcessMasrviPayment processes a MASRVI payment using provider-specific request
func (p *Provider) ProcessMasrviPayment(ctx context.Context, request *rimpay.MasrviPaymentRequest) (*types.PaymentResponse, error) {
	if request == nil {
		return nil, types.NewValidationError("request", "payment request cannot be nil")
	}

	if err := request.Validate(); err != nil {
		return nil, err
	}

	// Convert to generic request for internal processing
	genericRequest := request.ToGenericRequest()

	return p.ProcessPayment(ctx, genericRequest)
}

// ProcessPayment processes a payment with retry logic
func (p *Provider) ProcessPayment(ctx context.Context, request *types.PaymentRequest) (*types.PaymentResponse, error) {
	// Wrap the payment processing in a retryable function
	retryablePayment := func() (*types.PaymentResponse, error) {
		return p.paymentProcessor.ProcessPayment(ctx, request)
	}

	// Execute with retry logic
	return p.retryExecutor.ExecutePayment(ctx, retryablePayment)
}

// GetPaymentStatus retrieves payment status for MASRVI
func (p *Provider) GetPaymentStatus(ctx context.Context, transactionID string) (*rimpay.TransactionStatus, error) {
	if transactionID == "" {
		return nil, types.NewValidationError("transactionID", "transaction ID cannot be empty")
	}

	// For now, return a basic status since we don't have the full implementation
	status := &rimpay.TransactionStatus{
		TransactionID: transactionID,
		Status:        rimpay.PaymentStatusPending,
		Reference:     transactionID,
		Message:       "MASRVI payment status check",
		LastUpdated:   time.Now(),
	}

	return status, nil
}

// HandleNotification processes MASRVI webhook notifications
func (p *Provider) HandleNotification(notification *rimpay.MasrviNotificationData) (*rimpay.TransactionStatus, error) {
	// Convert to internal notification format
	internalNotification := &NotificationData{
		Status:      notification.Status,
		Mobile:      notification.PhoneNumber,
		PurchaseRef: notification.Reference,
		PaymentRef:  notification.TransactionID,
		Timestamp:   notification.Timestamp,
	}

	return p.paymentProcessor.HandleNotification(internalNotification)
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
