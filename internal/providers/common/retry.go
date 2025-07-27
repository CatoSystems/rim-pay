package common

import (
	"context"
	"github.com/CatoSystems/rim-pay/internal/types"
	"math"
	"math/rand"
	"time"
)

type RetryConfig struct {
	MaxAttempts  int           `json:"max_attempts"`
	InitialDelay time.Duration `json:"initial_delay"`
	MaxDelay     time.Duration `json:"max_delay"`
	Multiplier   float64       `json:"multiplier"`
	EnableJitter bool          `json:"enable_jitter"`
}

// DefaultRetryConfig returns default retry configuration
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxAttempts:  3,
		InitialDelay: 1 * time.Second,
		MaxDelay:     30 * time.Second,
		Multiplier:   2.0,
		EnableJitter: true,
	}
}

// RetryablePaymentFunc represents a payment function that can be retried
type RetryablePaymentFunc func() (*types.PaymentResponse, error)

// RetryExecutor handles retry logic
type RetryExecutor struct {
	config RetryConfig
}

// NewRetryExecutor creates a new retry executor
func NewRetryExecutor(config RetryConfig) *RetryExecutor {
	return &RetryExecutor{
		config: config,
	}
}

// ExecutePayment executes a payment function with retry logic
func (re *RetryExecutor) ExecutePayment(ctx context.Context, fn RetryablePaymentFunc) (*types.PaymentResponse, error) {
	var lastErr error
	var lastResp *types.PaymentResponse

	for attempt := 1; attempt <= re.config.MaxAttempts; attempt++ {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		resp, err := fn()
		if err == nil {
			return resp, nil
		}

		lastErr = err
		lastResp = resp

		// Check if error is retryable
		if paymentErr, ok := err.(*types.PaymentError); ok {
			if !paymentErr.IsRetryable() {
				return lastResp, err
			}
		}

		// Don't sleep after last attempt
		if attempt == re.config.MaxAttempts {
			break
		}

		delay := re.calculateDelay(attempt)
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(delay):
		}
	}

	return lastResp, lastErr
}

// calculateDelay calculates the delay for the next retry attempt
func (re *RetryExecutor) calculateDelay(attempt int) time.Duration {
	// Calculate exponential backoff
	delay := time.Duration(float64(re.config.InitialDelay) * math.Pow(re.config.Multiplier, float64(attempt-1)))

	// Apply maximum delay limit
	if delay > re.config.MaxDelay {
		delay = re.config.MaxDelay
	}

	// Apply jitter if enabled
	if re.config.EnableJitter && delay > 0 {
		jitter := time.Duration(rand.Int63n(int64(delay / 2)))
		delay = delay/2 + jitter
	}

	return delay
}
