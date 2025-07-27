package rimpay

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/CatoSystems/rim-pay/internal/validation"
)

// Client represents the main RimPay client
type Client struct {
	config    *Config
	providers map[string]PaymentProvider
	active    PaymentProvider
	logger    Logger
	validator Validator
	mu        sync.RWMutex
}

// ClientOption represents client configuration option
type ClientOption func(*Client) error

// NewClient creates a new RimPay client
func NewClient(config *Config, opts ...ClientOption) (*Client, error) {
	if config == nil {
		return nil, fmt.Errorf("configuration is required")
	}

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	client := &Client{
		config:    config,
		providers: make(map[string]PaymentProvider),
		logger:    newDefaultLogger(config.Logging),
		validator: validation.NewValidator(),
	}

	// Apply options
	for _, opt := range opts {
		if err := opt(client); err != nil {
			return nil, fmt.Errorf("failed to apply option: %w", err)
		}
	}

	// Initialize providers
	if err := client.initializeProviders(); err != nil {
		return nil, fmt.Errorf("failed to initialize providers: %w", err)
	}

	// Set default provider
	if err := client.SetProvider(config.DefaultProvider); err != nil {
		return nil, fmt.Errorf("failed to set default provider: %w", err)
	}

	client.logger.Info("RimPay client initialized",
		"environment", config.Environment,
		"default_provider", config.DefaultProvider,
	)

	return client, nil
}

// ProcessPayment processes a payment using active provider
func (c *Client) ProcessPayment(ctx context.Context, request *PaymentRequest) (*PaymentResponse, error) {
	start := time.Now()

	c.mu.RLock()
	provider := c.active
	c.mu.RUnlock()

	if provider == nil {
		return nil, NewPaymentError(ErrorCodeProviderError, "no active payment provider", "", false)
	}

	// Validate request
	if err := c.validator.ValidatePaymentRequest(request); err != nil {
		return nil, err
	}

	// Check if expired
	if request.IsExpired() {
		return nil, NewPaymentError(ErrorCodePaymentExpired, "payment request has expired", "", false)
	}

	// Check provider availability
	if !provider.IsAvailable(ctx) {
		return nil, NewPaymentError(
			ErrorCodeProviderError,
			"provider is not available",
			provider.Name(),
			false,
		)
	}

	c.logger.Info("Processing payment",
		"provider", provider.Name(),
		"reference", request.Reference,
		"amount", request.Amount.String(),
	)

	// Process payment
	response, err := provider.ProcessPayment(ctx, request)

	duration := time.Since(start)

	if err != nil {
		c.logger.Error("Payment processing failed",
			"provider", provider.Name(),
			"reference", request.Reference,
			"error", err.Error(),
			"duration", duration,
		)
		return nil, err
	}

	c.logger.Info("Payment processed",
		"provider", provider.Name(),
		"transaction_id", response.TransactionID,
		"status", response.Status,
		"duration", duration,
	)

	return response, nil
}

// GetPaymentStatus retrieves payment status
func (c *Client) GetPaymentStatus(ctx context.Context, transactionID string) (*TransactionStatus, error) {
	c.mu.RLock()
	provider := c.active
	c.mu.RUnlock()

	if provider == nil {
		return nil, NewPaymentError(ErrorCodeProviderError, "no active payment provider", "", false)
	}

	if transactionID == "" {
		return nil, NewValidationError("transaction_id", "is required")
	}

	return provider.GetPaymentStatus(ctx, transactionID)
}

// SetProvider sets the active payment provider
func (c *Client) SetProvider(providerName string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	provider, exists := c.providers[providerName]
	if !exists {
		return fmt.Errorf("provider not found: %s", providerName)
	}

	if err := provider.ValidateConfig(); err != nil {
		return fmt.Errorf("provider configuration invalid: %w", err)
	}

	c.active = provider
	c.logger.Info("Active provider changed", "provider", providerName)
	return nil
}

// GetActiveProvider returns currently active provider
func (c *Client) GetActiveProvider() PaymentProvider {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.active
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

// initializeProviders initializes payment providers
func (c *Client) initializeProviders() error {
	// TODO: Register providers here
	
	for name, config := range c.config.Providers {
		if !config.Enabled {
			c.logger.Debug("Skipping disabled provider", "provider", name)
			continue
		}

		provider, err := c.createProvider(name, config)
		if err != nil {
			return fmt.Errorf("failed to create provider '%s': %w", name, err)
		}

		c.providers[name] = provider
		c.logger.Info("Provider initialized", "provider", name)
	}

	if len(c.providers) == 0 {
		return fmt.Errorf("no enabled providers found")
	}

	return nil
}

// createProvider creates a payment provider instance
func (c *Client) createProvider(name string, config ProviderConfig) (PaymentProvider, error) {
	return DefaultRegistry.Create(name, config, c.logger)
}

// WithLogger sets custom logger
func WithLogger(logger Logger) ClientOption {
	return func(c *Client) error {
		c.logger = logger
		return nil
	}
}

// WithValidator sets custom validator
func WithValidator(validator Validator) ClientOption {
	return func(c *Client) error {
		c.validator = validator
		return nil
	}
}

// newDefaultLogger creates default logger
func newDefaultLogger(config LoggingConfig) Logger {
	// Implementation would create appropriate logger based on config
	return &defaultLogger{}
}

// defaultLogger implements basic logging
type defaultLogger struct{}

func (l *defaultLogger) Debug(msg string, fields ...interface{}) {}
func (l *defaultLogger) Info(msg string, fields ...interface{})  {}
func (l *defaultLogger) Warn(msg string, fields ...interface{})  {}
func (l *defaultLogger) Error(msg string, fields ...interface{}) {}
