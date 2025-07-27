
package rimpay

import (
	"context"
	"fmt"
	"sync"
)

// Provider constants
const (
	ProviderBPay   = "bpay"
	ProviderMasrvi = "masrvi"
	
	// Error message constants
	providerNotAvailableMsg = "provider %s not available"
)

// Factory functions for creating providers - these will be set by provider packages
var (
	createBPayProvider   func(ProviderConfig, Logger) (PaymentProvider, error)
	createMasrviProvider func(ProviderConfig, Logger) (PaymentProvider, error)
)

// RegisterBPayProvider registers the B-PAY provider factory
func RegisterBPayProvider(factory func(ProviderConfig, Logger) (PaymentProvider, error)) {
	createBPayProvider = factory
}

// RegisterMasrviProvider registers the MASRVI provider factory
func RegisterMasrviProvider(factory func(ProviderConfig, Logger) (PaymentProvider, error)) {
	createMasrviProvider = factory
}

// Client represents the main payment client
type Client struct {
	providers map[string]PaymentProvider
	config    *Config
	logger    Logger
	mu        sync.RWMutex
}

// NewClient creates a new payment client
func NewClient(config *Config) (*Client, error) {
	if config == nil {
		return nil, ErrInvalidConfig
	}

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	// Create a default logger if none provided
	logger := newDefaultLogger(config.Logging)

	return &Client{
		providers: make(map[string]PaymentProvider),
		config:    config,
		logger:    logger,
	}, nil
}

// newDefaultLogger creates a default logger
func newDefaultLogger(config LoggingConfig) Logger {
	return &simpleLogger{}
}

// Simple logger implementation
type simpleLogger struct{}

func (l *simpleLogger) Debug(msg string, fields ...interface{}) {
	fmt.Printf("[DEBUG] %s %v\n", msg, fields)
}

func (l *simpleLogger) Info(msg string, fields ...interface{}) {
	fmt.Printf("[INFO] %s %v\n", msg, fields)
}

func (l *simpleLogger) Warn(msg string, fields ...interface{}) {
	fmt.Printf("[WARN] %s %v\n", msg, fields)
}

func (l *simpleLogger) Error(msg string, fields ...interface{}) {
	fmt.Printf("[ERROR] %s %v\n", msg, fields)
}

// ProcessBPayPayment processes a payment using B-PAY provider
func (c *Client) ProcessBPayPayment(ctx context.Context, request *BPayPaymentRequest) (*PaymentResponse, error) {
	if request == nil {
		return nil, ErrInvalidRequest
	}

	provider, ok := c.providers[ProviderBPay]
	if !ok {
		return nil, fmt.Errorf(providerNotAvailableMsg, ProviderBPay)
	}

	bpayProvider, ok := provider.(BPayProvider)
	if !ok {
		return nil, fmt.Errorf("provider %s does not implement BPayProvider interface", ProviderBPay)
	}

	return bpayProvider.ProcessBPayPayment(ctx, request)
}

// ProcessMasrviPayment processes a payment using MASRVI provider
func (c *Client) ProcessMasrviPayment(ctx context.Context, request *MasrviPaymentRequest) (*PaymentResponse, error) {
	if request == nil {
		return nil, ErrInvalidRequest
	}

	provider, ok := c.providers[ProviderMasrvi]
	if !ok {
		return nil, fmt.Errorf(providerNotAvailableMsg, ProviderMasrvi)
	}

	masrviProvider, ok := provider.(MasrviProvider)
	if !ok {
		return nil, fmt.Errorf("provider %s does not implement MasrviProvider interface", ProviderMasrvi)
	}

	return masrviProvider.ProcessMasrviPayment(ctx, request)
}

// HandleMasrviNotification handles MASRVI webhook notifications
func (c *Client) HandleMasrviNotification(notification *MasrviNotificationData) (*TransactionStatus, error) {
	if notification == nil {
		return nil, ErrInvalidRequest
	}

	provider, ok := c.providers[ProviderMasrvi]
	if !ok {
		return nil, fmt.Errorf(providerNotAvailableMsg, ProviderMasrvi)
	}

	masrviProvider, ok := provider.(MasrviProvider)
	if !ok {
		return nil, fmt.Errorf("provider %s does not implement MasrviProvider interface", ProviderMasrvi)
	}

	return masrviProvider.HandleNotification(notification)
}

// ProcessPayment processes a payment using the generic interface (deprecated)
func (c *Client) ProcessPayment(ctx context.Context, request *PaymentRequest) (*PaymentResponse, error) {
	if request == nil {
		return nil, ErrInvalidRequest
	}

	// For backward compatibility, use the first available provider
	c.mu.RLock()
	var provider PaymentProvider
	for _, p := range c.providers {
		provider = p
		break
	}
	c.mu.RUnlock()

	if provider == nil {
		return nil, ErrProviderNotFound
	}

	// Check provider availability
	if !provider.IsAvailable(ctx) {
		return nil, fmt.Errorf("provider %s is not available", provider.Name())
	}

	// Process payment
	return provider.ProcessPayment(ctx, request)
}

// GetPaymentStatus retrieves payment status from the first available provider
func (c *Client) GetPaymentStatus(ctx context.Context, transactionID string) (*TransactionStatus, error) {
	if transactionID == "" {
		return nil, ErrInvalidRequest
	}

	c.mu.RLock()
	var provider PaymentProvider
	for _, p := range c.providers {
		provider = p
		break
	}
	c.mu.RUnlock()

	if provider == nil {
		return nil, ErrProviderNotFound
	}

	return provider.GetPaymentStatus(ctx, transactionID)
}

// AddProvider adds a payment provider to the client
func (c *Client) AddProvider(name string, provider PaymentProvider) error {
	if provider == nil {
		return ErrInvalidProvider
	}

	c.mu.Lock()
	c.providers[name] = provider
	c.mu.Unlock()

	c.logger.Info("Provider added", "name", name, "provider", provider.Name())
	return nil
}

// ListProviders returns list of available providers
func (c *Client) ListProviders() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var providers []string
	for name := range c.providers {
		providers = append(providers, name)
	}
	return providers
}
