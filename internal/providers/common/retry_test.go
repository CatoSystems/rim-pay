package common

import (
	"context"
	"testing"
	"time"

	"github.com/CatoSystems/rim-pay/internal/types"
	"github.com/CatoSystems/rim-pay/pkg/money"
	"github.com/shopspring/decimal"
)

func TestRetryExecutorExecutePayment(t *testing.T) {
	t.Run("SuccessOnFirstAttempt", testSuccessOnFirstAttempt)
	t.Run("SuccessAfterRetries", testSuccessAfterRetries)
	t.Run("NonRetryableErrorFailsImmediately", testNonRetryableErrorFailsImmediately)
}

const networkErrorMsg = "network error"

func testSuccessOnFirstAttempt(t *testing.T) {
	config := RetryConfig{
		MaxAttempts:  3,
		InitialDelay: 10 * time.Millisecond,
		MaxDelay:     100 * time.Millisecond,
		Multiplier:   2.0,
		EnableJitter: false,
	}
	executor := NewRetryExecutor(config)
	attemptCount := 0

	mockFunc := func() (*types.PaymentResponse, error) {
		attemptCount++
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

	ctx := context.Background()
	resp, err := executor.ExecutePayment(ctx, mockFunc)

	if attemptCount != 1 {
		t.Errorf("Expected 1 attempt, got %d", attemptCount)
	}
	if err != nil {
		t.Errorf("Expected success, got error: %v", err)
	}
	if resp == nil {
		t.Error("Expected response, got nil")
	}
}

func testSuccessAfterRetries(t *testing.T) {
	config := RetryConfig{
		MaxAttempts:  3,
		InitialDelay: 10 * time.Millisecond,
		MaxDelay:     100 * time.Millisecond,
		Multiplier:   2.0,
		EnableJitter: false,
	}
	executor := NewRetryExecutor(config)
	attemptCount := 0

	mockFunc := func() (*types.PaymentResponse, error) {
		attemptCount++
		if attemptCount == 3 {
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
		return nil, types.NewPaymentError(
			types.ErrorCodeNetworkError,
			networkErrorMsg,
			"test",
			true,
		)
	}

	ctx := context.Background()
	resp, err := executor.ExecutePayment(ctx, mockFunc)

	if attemptCount != 3 {
		t.Errorf("Expected 3 attempts, got %d", attemptCount)
	}
	if err != nil {
		t.Errorf("Expected success, got error: %v", err)
	}
	if resp == nil {
		t.Error("Expected response, got nil")
	}
}

func testNonRetryableErrorFailsImmediately(t *testing.T) {
	config := RetryConfig{
		MaxAttempts:  3,
		InitialDelay: 10 * time.Millisecond,
		MaxDelay:     100 * time.Millisecond,
		Multiplier:   2.0,
		EnableJitter: false,
	}
	executor := NewRetryExecutor(config)
	attemptCount := 0

	mockFunc := func() (*types.PaymentResponse, error) {
		attemptCount++
		return nil, types.NewPaymentError(
			types.ErrorCodeNetworkError,
			networkErrorMsg,
			"test",
			false,
		)
	}

	ctx := context.Background()
	resp, err := executor.ExecutePayment(ctx, mockFunc)

	if attemptCount != 1 {
		t.Errorf("Expected 1 attempt, got %d", attemptCount)
	}
	if err == nil {
		t.Error("Expected error, got success")
	}
	if resp != nil {
		t.Error("Expected nil response, got non-nil")
	}
}

func TestRetryExecutorContextCancellation(t *testing.T) {
	config := DefaultRetryConfig()
	config.InitialDelay = 100 * time.Millisecond // Longer delay to test cancellation
	
	executor := NewRetryExecutor(config)
	
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	mockFunc := func() (*types.PaymentResponse, error) {
		return nil, types.NewPaymentError(
			types.ErrorCodeNetworkError,
			networkErrorMsg,
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
