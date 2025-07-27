package common

import (
	"context"
	"testing"
	"time"

	"github.com/CatoSystems/rim-pay/internal/types"
	"github.com/CatoSystems/rim-pay/pkg/money"
	"github.com/shopspring/decimal"
)

func TestRetryExecutor_ExecutePayment(t *testing.T) {
	tests := []struct {
		name           string
		maxAttempts    int
		shouldSucceed  bool
		expectAttempts int
		errorRetryable bool
	}{
		{
			name:           "Success on first attempt",
			maxAttempts:    3,
			shouldSucceed:  true,
			expectAttempts: 1,
			errorRetryable: false,
		},
		{
			name:           "Success after retries",
			maxAttempts:    3,
			shouldSucceed:  false, // Will succeed on 3rd attempt in our mock
			expectAttempts: 3,
			errorRetryable: true,
		},
		{
			name:           "Non-retryable error fails immediately",
			maxAttempts:    3,
			shouldSucceed:  false,
			expectAttempts: 1,
			errorRetryable: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := RetryConfig{
				MaxAttempts:  tt.maxAttempts,
				InitialDelay: 10 * time.Millisecond, // Short delay for testing
				MaxDelay:     100 * time.Millisecond,
				Multiplier:   2.0,
				EnableJitter: false, // Disable jitter for predictable testing
			}

			executor := NewRetryExecutor(config)
			attemptCount := 0

			mockFunc := func() (*types.PaymentResponse, error) {
				attemptCount++
				
				if tt.shouldSucceed || (tt.name == "Success after retries" && attemptCount == 3) {
					// Success case
					amount := money.New(decimal.NewFromInt(1000), "MRU")
					return &types.PaymentResponse{
						TransactionID: "test-123",
						Status:        types.PaymentStatusSuccess,
						Amount:        amount,
						Reference:     "ref-123",
						Provider:      "test",
						CreatedAt:     time.Now(),
						UpdatedAt:     time.Now(),
					}, nil
				}

				// Error case
				return nil, types.NewPaymentError(
					types.ErrorCodeNetworkError,
					"network error",
					"test",
					tt.errorRetryable,
				)
			}

			ctx := context.Background()
			resp, err := executor.ExecutePayment(ctx, mockFunc)

			// Check attempt count
			if attemptCount != tt.expectAttempts {
				t.Errorf("Expected %d attempts, got %d", tt.expectAttempts, attemptCount)
			}

			// Check result based on test case
			if tt.shouldSucceed || (tt.name == "Success after retries") {
				if err != nil {
					t.Errorf("Expected success, got error: %v", err)
				}
				if resp == nil {
					t.Error("Expected response, got nil")
				}
			} else {
				if err == nil {
					t.Error("Expected error, got success")
				}
				// For non-retryable errors, should fail immediately
				if !tt.errorRetryable && attemptCount != 1 {
					t.Errorf("Non-retryable error should fail immediately, but got %d attempts", attemptCount)
				}
			}
		})
	}
}

func TestRetryExecutor_ContextCancellation(t *testing.T) {
	config := DefaultRetryConfig()
	config.InitialDelay = 100 * time.Millisecond // Longer delay to test cancellation
	
	executor := NewRetryExecutor(config)
	
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	
	mockFunc := func() (*types.PaymentResponse, error) {
		return nil, types.NewPaymentError(
			types.ErrorCodeNetworkError,
			"network error",
			"test",
			true, // retryable
		)
	}
	
	start := time.Now()
	_, err := executor.ExecutePayment(ctx, mockFunc)
	duration := time.Since(start)
	
	if err != context.DeadlineExceeded {
		t.Errorf("Expected context.DeadlineExceeded, got %v", err)
	}
	
	// Should cancel quickly, not wait for full retry delay
	if duration > 200*time.Millisecond {
		t.Errorf("Expected quick cancellation, but took %v", duration)
	}
}

func TestDefaultRetryConfig(t *testing.T) {
	config := DefaultRetryConfig()
	
	if config.MaxAttempts != 3 {
		t.Errorf("Expected MaxAttempts=3, got %d", config.MaxAttempts)
	}
	if config.InitialDelay != 1*time.Second {
		t.Errorf("Expected InitialDelay=1s, got %v", config.InitialDelay)
	}
	if config.MaxDelay != 30*time.Second {
		t.Errorf("Expected MaxDelay=30s, got %v", config.MaxDelay)
	}
	if config.Multiplier != 2.0 {
		t.Errorf("Expected Multiplier=2.0, got %f", config.Multiplier)
	}
	if !config.EnableJitter {
		t.Error("Expected EnableJitter=true")
	}
}
