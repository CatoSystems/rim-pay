package rimpay

import (
	"fmt"
	"time"
)

type Environment string

const (
	// EnvironmentSandbox for testing
	EnvironmentSandbox Environment = "sandbox"
	// EnvironmentProduction for live transactions
	EnvironmentProduction Environment = "production"
)

// Config represents main configuration
type Config struct {
	Environment     Environment               `json:"environment"`
	DefaultProvider string                    `json:"default_provider"`
	Providers       map[string]ProviderConfig `json:"providers"`
	HTTP            HTTPConfig                `json:"http"`
	Logging         LoggingConfig             `json:"logging"`
	Security        SecurityConfig            `json:"security"`
}

// ProviderConfig represents provider configuration
type ProviderConfig struct {
	Enabled     bool                   `json:"enabled"`
	BaseURL     string                 `json:"base_url"`
	Credentials map[string]string      `json:"credentials"`
	Timeout     time.Duration          `json:"timeout"`
	Options     map[string]interface{} `json:"options"`
}

// HTTPConfig represents HTTP configuration
type HTTPConfig struct {
	Timeout         time.Duration `json:"timeout"`
	MaxIdleConns    int           `json:"max_idle_conns"`
	MaxConnsPerHost int           `json:"max_conns_per_host"`
	UserAgent       string        `json:"user_agent"`
}

// LoggingConfig represents logging configuration
type LoggingConfig struct {
	Level  string `json:"level"`
	Format string `json:"format"`
	Output string `json:"output"`
}

// SecurityConfig represents security configuration
type SecurityConfig struct {
	EncryptionKey string        `json:"encryption_key"`
	SigningKey    string        `json:"signing_key"`
	TokenTTL      time.Duration `json:"token_ttl"`
}

// DefaultConfig returns default configuration
func DefaultConfig() *Config {
	return &Config{
		Environment:     EnvironmentSandbox,
		DefaultProvider: "bpay",
		Providers:       make(map[string]ProviderConfig),
		HTTP: HTTPConfig{
			Timeout:         30 * time.Second,
			MaxIdleConns:    100,
			MaxConnsPerHost: 10,
			UserAgent:       "RimPay/1.0",
		},
		Logging: LoggingConfig{
			Level:  "info",
			Format: "json",
			Output: "stdout",
		},
		Security: SecurityConfig{
			TokenTTL: time.Hour,
		},
	}
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.Environment != EnvironmentSandbox && c.Environment != EnvironmentProduction {
		return fmt.Errorf("invalid environment: %s", c.Environment)
	}

	if c.DefaultProvider == "" {
		return fmt.Errorf("default provider must be specified")
	}

	if _, exists := c.Providers[c.DefaultProvider]; !exists {
		return fmt.Errorf("default provider '%s' not found in providers", c.DefaultProvider)
	}

	for name, provider := range c.Providers {
		if err := c.validateProviderConfig(name, provider); err != nil {
			return fmt.Errorf("invalid config for provider '%s': %w", name, err)
		}
	}

	return nil
}

// validateProviderConfig validates provider configuration
func (c *Config) validateProviderConfig(name string, config ProviderConfig) error {
	if !config.Enabled {
		return nil
	}

	if config.BaseURL == "" {
		return fmt.Errorf("base_url is required")
	}

	if config.Timeout <= 0 {
		return fmt.Errorf("timeout must be positive")
	}

	return nil
}

// GetProviderConfig returns provider configuration
func (c *Config) GetProviderConfig(name string) (ProviderConfig, bool) {
	config, exists := c.Providers[name]
	return config, exists
}

// IsProduction returns true if production environment
func (c *Config) IsProduction() bool {
	return c.Environment == EnvironmentProduction
}
