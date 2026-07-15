package click

import (
	"context"
	"fmt"
	"time"

	"github.com/CatoSystems/rim-pay/internal/providers/common"
	"github.com/CatoSystems/rim-pay/internal/types"
	"github.com/CatoSystems/rim-pay/pkg/rimpay"
)

// Register the CLICK provider factory with rimpay.
func init() {
	rimpay.RegisterClickProvider(func(config rimpay.ProviderConfig, logger rimpay.Logger) (rimpay.PaymentProvider, error) {
		return NewClickProvider(config, logger)
	})
}

// Provider implements the CLICK payment provider.
type Provider struct {
	name             string
	config           rimpay.ProviderConfig
	httpClient       common.HTTPClient
	sessionManager   *SessionManager
	paymentProcessor *PaymentProcessor
	retryExecutor    *common.RetryExecutor
	logger           rimpay.Logger
}

// NewClickProvider creates a new CLICK provider.
func NewClickProvider(config rimpay.ProviderConfig, logger rimpay.Logger) (*Provider, error) {
	if err := validateConfig(config); err != nil {
		return nil, fmt.Errorf("invalid CLICK configuration: %w", err)
	}

	httpClient := common.NewHTTPClient(common.HTTPConfig{
		Timeout:         config.Timeout,
		MaxIdleConns:    10,
		MaxConnsPerHost: 5,
	})
	sessionManager := NewSessionManager(config, httpClient, logger)
	paymentProcessor := NewPaymentProcessor(config, httpClient, sessionManager, logger)
	retryExecutor := common.NewRetryExecutor(common.DefaultRetryConfig())

	return &Provider{
		name:             "click",
		config:           config,
		httpClient:       httpClient,
		sessionManager:   sessionManager,
		paymentProcessor: paymentProcessor,
		retryExecutor:    retryExecutor,
		logger:           logger,
	}, nil
}

// NewProvider is an alias for registry use.
func NewProvider(config rimpay.ProviderConfig, logger rimpay.Logger) (*Provider, error) {
	return NewClickProvider(config, logger)
}

// Name returns the provider name.
func (p *Provider) Name() string { return p.name }

// IsAvailable checks if the provider can obtain a session.
func (p *Provider) IsAvailable(ctx context.Context) bool {
	_, err := p.sessionManager.GetSessionID(ctx)
	return err == nil
}

// ProcessClickPayment validates and processes a CLICK-specific request.
func (p *Provider) ProcessClickPayment(ctx context.Context, request *rimpay.ClickPaymentRequest) (*types.PaymentResponse, error) {
	if request == nil {
		return nil, types.NewValidationError("request", "payment request cannot be nil")
	}
	if err := request.Validate(); err != nil {
		return nil, err
	}
	return p.ProcessPayment(ctx, request.ToGenericRequest())
}

// ProcessPayment processes a generic request with retry.
func (p *Provider) ProcessPayment(ctx context.Context, request *types.PaymentRequest) (*types.PaymentResponse, error) {
	return p.retryExecutor.ExecutePayment(ctx, func() (*types.PaymentResponse, error) {
		return p.paymentProcessor.ProcessPayment(ctx, request)
	})
}

// GetPaymentStatus is notification-driven for CLICK; returns a pending placeholder.
func (p *Provider) GetPaymentStatus(ctx context.Context, transactionID string) (*rimpay.TransactionStatus, error) {
	if transactionID == "" {
		return nil, types.NewValidationError("transactionID", "transaction ID cannot be empty")
	}
	return &rimpay.TransactionStatus{
		TransactionID: transactionID,
		Status:        rimpay.PaymentStatusPending,
		Reference:     transactionID,
		Message:       "CLICK payment status is delivered via notification",
		LastUpdated:   time.Now(),
	}, nil
}

// HandleNotification converts a public notification into a TransactionStatus.
func (p *Provider) HandleNotification(notification *rimpay.ClickNotificationData) (*rimpay.TransactionStatus, error) {
	if notification == nil {
		return nil, types.NewValidationError("notification", "is required")
	}
	return p.paymentProcessor.HandleNotification(&NotificationData{
		Status:      notification.Status,
		ClientID:    notification.ClientID,
		ClientName:  notification.ClientName,
		Mobile:      notification.Mobile,
		PurchaseRef: notification.PurchaseRef,
		PaymentRef:  notification.PaymentRef,
		PayID:       notification.PayID,
		Timestamp:   notification.Timestamp,
		IPAddress:   notification.IPAddress,
		Error:       notification.Error,
		Reason:      notification.Reason,
		Amount:      notification.Amount,
		Currency:    notification.Currency,
	})
}

// ValidateConfig validates provider configuration.
func (p *Provider) ValidateConfig() error { return validateConfig(p.config) }

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
