package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/CatoSystems/rim-pay/pkg/money"
	"github.com/CatoSystems/rim-pay/pkg/phone"
	_ "github.com/CatoSystems/rim-pay/pkg/providers" // register all providers
	"github.com/CatoSystems/rim-pay/pkg/rimpay"
	"github.com/shopspring/decimal"
)

func main() {
	// Create client configuration
	config := &rimpay.Config{
		Environment:     rimpay.EnvironmentSandbox,
		DefaultProvider: "bpay",
		Providers: map[string]rimpay.ProviderConfig{
			"bpay": {
				Enabled: true,
				BaseURL: "https://api.sandbox.bpay.mr",
				Timeout: 30 * time.Second,
				Credentials: map[string]string{
					"username":  "test_user",
					"password":  "test_pass",
					"client_id": "test_client",
				},
			},
		},
	}

	// Create client
	client, err := rimpay.NewClient(config)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Register the B-PAY provider instance on the client.
	if err := client.AddBPayProvider(config.Providers["bpay"]); err != nil {
		log.Fatalf("Failed to add B-PAY provider: %v", err)
	}

	// Create phone number
	phoneNumber, err := phone.NewPhone("22334455")
	if err != nil {
		log.Fatalf("Failed to create phone: %v", err)
	}

	// Create amount
	amount := money.New(decimal.NewFromFloat(100.00), money.MRU)

	// Create payment request
	request := &rimpay.BPayPaymentRequest{
		Amount:      amount,
		PhoneNumber: phoneNumber,
		Reference:   "TEST-" + fmt.Sprintf("%d", time.Now().Unix()),
		Description: "Retry demo payment",
		Passcode:    "1234",
	}

	ctx := context.Background()

	fmt.Printf("🚀 Processing payment with retry functionality...\n")
	fmt.Printf("   Amount: %s\n", amount.String())
	fmt.Printf("   Phone: %s\n", phoneNumber.ForProvider(true))
	fmt.Printf("   Reference: %s\n", request.Reference)
	fmt.Printf("   Provider: bpay (with automatic retry on failures)\n\n")

	// Process payment - this will use the retry functionality
	// If the first attempt fails with a retryable error, it will automatically retry
	resp, err := client.ProcessBPayPayment(ctx, request)
	if err != nil {
		// Even if this fails, the retry mechanism would have attempted up to 3 times
		// for retryable errors before giving up
		fmt.Printf("❌ Payment failed after retries: %v\n", err)

		// Check if it's a payment error with retry information
		if paymentErr, ok := err.(*rimpay.PaymentError); ok {
			fmt.Printf("   Error Code: %s\n", paymentErr.Code)
			fmt.Printf("   Provider: %s\n", paymentErr.Provider)
			fmt.Printf("   Was Retryable: %v\n", paymentErr.IsRetryable())
		}
		return
	}

	// Success
	fmt.Printf("✅ Payment processed successfully!\n")
	fmt.Printf("   Transaction ID: %s\n", resp.TransactionID)
	fmt.Printf("   Status: %s\n", resp.Status)
	fmt.Printf("   Provider: %s\n", resp.Provider)
	fmt.Printf("   Amount: %s\n", resp.Amount.String())

	fmt.Printf("\n🔄 Retry Configuration:\n")
	fmt.Printf("   Max Attempts: 3\n")
	fmt.Printf("   Initial Delay: 1s\n")
	fmt.Printf("   Max Delay: 30s\n")
	fmt.Printf("   Backoff Multiplier: 2.0x\n")
	fmt.Printf("   Jitter: Enabled\n")

	fmt.Printf("\n📝 How retry works:\n")
	fmt.Printf("   • Network errors → Retryable\n")
	fmt.Printf("   • Authentication failures → Retryable\n")
	fmt.Printf("   • Server errors (5xx) → Retryable\n")
	fmt.Printf("   • Validation errors → Not retryable\n")
	fmt.Printf("   • Insufficient funds → Not retryable\n")
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
