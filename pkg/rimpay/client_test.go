package rimpay

import (
	"testing"
	"time"

	"github.com/CatoSystems/rim-pay/pkg/money"
	phone "github.com/CatoSystems/rim-pay/pkg/phone"
	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	config := DefaultConfig()
	config.Providers["test"] = ProviderConfig{
		Enabled: true,
		BaseURL: "https://test.example.com",
		Credentials: map[string]string{
			"username": "test",
			"password": "test",
		},
		Timeout: 30 * time.Second,
	}
	config.DefaultProvider = "test"

	// With the new architecture, NewClient should succeed
	// as providers are registered and initialized lazily
	client, err := NewClient(config)
	assert.NoError(t, err)
	assert.NotNil(t, client)
}

func TestConfigValidation(t *testing.T) {
	tests := []struct {
		name      string
		config    *Config
		wantError bool
	}{
		{
			name:      "nil config",
			config:    nil,
			wantError: true,
		},
		{
			name: "valid config",
			config: &Config{
				Environment:     EnvironmentSandbox,
				DefaultProvider: "test",
				Providers: map[string]ProviderConfig{
					"test": {
						Enabled: true,
						BaseURL: "https://test.com",
						Timeout: 30 * time.Second,
					},
				},
			},
			wantError: false,
		},
		{
			name: "invalid environment",
			config: &Config{
				Environment:     "invalid",
				DefaultProvider: "test",
				Providers:       map[string]ProviderConfig{},
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewClient(tt.config)
			if tt.wantError {
				assert.Error(t, err)
			} else {
				// We expect error due to missing provider implementation
				// but config validation should pass
				if tt.config != nil {
					assert.NoError(t, tt.config.Validate())
				}
			}
		})
	}
}

func TestPaymentRequestValidation(t *testing.T) {
	phoneNumber, _ := phone.NewPhone("+22222334455")

	tests := []struct {
		name    string
		request *PaymentRequest
		wantErr bool
	}{
		{
			name: "valid request",
			request: &PaymentRequest{
				Amount:      money.FromFloat64(100, money.MRU),
				PhoneNumber: phoneNumber,
				Reference:   "TEST123",
			},
			wantErr: false,
		},
		{
			name: "zero amount",
			request: &PaymentRequest{
				Amount:      money.FromFloat64(0, money.MRU),
				PhoneNumber: phoneNumber,
				Reference:   "TEST123",
			},
			wantErr: true,
		},
		{
			name: "negative amount",
			request: &PaymentRequest{
				Amount:      money.FromFloat64(-100, money.MRU),
				PhoneNumber: phoneNumber,
				Reference:   "TEST123",
			},
			wantErr: true,
		},
		{
			name: "missing phone",
			request: &PaymentRequest{
				Amount:    money.FromFloat64(100, money.MRU),
				Reference: "TEST123",
			},
			wantErr: true,
		},
		{
			name: "missing reference",
			request: &PaymentRequest{
				Amount:      money.FromFloat64(100, money.MRU),
				PhoneNumber: phoneNumber,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
