package masrvi

import (
	"github.com/CatoSystems/rim-pay/pkg/rimpay"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
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
				BaseURL: "https://test.masrvi.com",
				Credentials: map[string]string{
					"merchant_id": "test123",
				},
				Timeout: 30 * time.Second,
			},
			wantError: false,
		},
		{
			name: "missing merchant_id",
			config: rimpay.ProviderConfig{
				BaseURL:     "https://test.masrvi.com",
				Credentials: map[string]string{},
				Timeout:     30 * time.Second,
			},
			wantError: true,
		},
		{
			name: "missing base URL",
			config: rimpay.ProviderConfig{
				Credentials: map[string]string{
					"merchant_id": "test123",
				},
				Timeout: 30 * time.Second,
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

func TestNotificationToPaymentStatus(t *testing.T) {
	tests := []struct {
		status   string
		expected rimpay.PaymentStatus
	}{
		{"Ok", rimpay.PaymentStatusSuccess},
		{"NOK", rimpay.PaymentStatusFailed},
		{"UNKNOWN", rimpay.PaymentStatusPending},
	}

	for _, tt := range tests {
		t.Run(tt.status, func(t *testing.T) {
			notification := &NotificationData{Status: tt.status}
			result := notification.ToPaymentStatus()
			assert.Equal(t, tt.expected, result)
		})
	}
}
