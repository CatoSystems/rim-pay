package bpay

import (
	"testing"
	"time"

	"github.com/CatoSystems/rim-pay/pkg/rimpay"
	"github.com/stretchr/testify/assert"
)

func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name      string
		config    rimpay.ProviderConfig
		wantError bool
	}{
		{
			name: "valid config",
			config: rimpay.ProviderConfig{
				BaseURL: "https://test.bpay.com",
				Credentials: map[string]string{
					"username":  "test",
					"password":  "test",
					"client_id": "test",
				},
				Timeout: 30 * time.Second,
			},
			wantError: false,
		},
		{
			name: "missing username",
			config: rimpay.ProviderConfig{
				BaseURL: "https://test.bpay.com",
				Credentials: map[string]string{
					"password":  "test",
					"client_id": "test",
				},
				Timeout: 30 * time.Second,
			},
			wantError: true,
		},
		{
			name: "missing base URL",
			config: rimpay.ProviderConfig{
				Credentials: map[string]string{
					"username":  "test",
					"password":  "test",
					"client_id": "test",
				},
				Timeout: 30 * time.Second,
			},
			wantError: true,
		},
		{
			name: "invalid timeout",
			config: rimpay.ProviderConfig{
				BaseURL: "https://test.bpay.com",
				Credentials: map[string]string{
					"username":  "test",
					"password":  "test",
					"client_id": "test",
				},
				Timeout: 0,
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateConfig(tt.config)
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestConvertErrorCodeToStatus(t *testing.T) {
	tests := []struct {
		errorCode string
		expected  rimpay.PaymentStatus
	}{
		{"0", rimpay.PaymentStatusSuccess},
		{"1", rimpay.PaymentStatusFailed},
		{"2", rimpay.PaymentStatusFailed},
		{"4", rimpay.PaymentStatusFailed},
		{"999", rimpay.PaymentStatusPending},
	}

	for _, tt := range tests {
		t.Run(tt.errorCode, func(t *testing.T) {
			result := convertErrorCodeToStatus(tt.errorCode)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestConvertTransactionStatus(t *testing.T) {
	tests := []struct {
		status   string
		expected rimpay.PaymentStatus
	}{
		{"TS", rimpay.PaymentStatusSuccess},
		{"TF", rimpay.PaymentStatusFailed},
		{"TA", rimpay.PaymentStatusPending},
		{"UNKNOWN", rimpay.PaymentStatusPending},
	}

	for _, tt := range tests {
		t.Run(tt.status, func(t *testing.T) {
			result := convertTransactionStatus(tt.status)
			assert.Equal(t, tt.expected, result)
		})
	}
}
