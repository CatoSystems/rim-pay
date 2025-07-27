package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/CatoSystems/rim-pay/pkg/money"
	"github.com/CatoSystems/rim-pay/pkg/phone"
	"github.com/CatoSystems/rim-pay/pkg/rimpay"

	// Import providers to register them
	_ "github.com/CatoSystems/rim-pay/internal/providers/bpay"
	_ "github.com/CatoSystems/rim-pay/internal/providers/masrvi"
)

// Provider-specific examples demonstrating type-safe payment requests

const phoneCreationError = "failed to create phone number: %w"

func main() {
	fmt.Println("=== RimPay Provider-Specific API Examples ===\n")

	// Example 1: B-PAY specific payment
	if err := runBPayExample(); err != nil {
		log.Printf("B-PAY example failed: %v\n", err)
	}

	// Example 2: MASRVI specific payment
	if err := runMasrviExample(); err != nil {
		log.Printf("MASRVI example failed: %v\n", err)
	}

	// Example 3: Multi-provider with type safety
	if err := runMultiProviderExample(); err != nil {
		log.Printf("Multi-provider example failed: %v\n", err)
	}
}

func runBPayExample() error {
	fmt.Println("--- B-PAY Specific Payment Example ---")

	// Create client
	client, err := createTestClient()
	if err != nil {
		return err
	}

	// Add B-PAY provider
	bpayConfig := rimpay.ProviderConfig{
		Enabled: true,
		BaseURL: "https://api.bpay.mr",
		Credentials: map[string]string{
			"username":  "test_username",
			"password":  "test_password",
			"client_id": "test_client_id",
		},
		Timeout: 30 * time.Second,
	}

	if err := client.AddBPayProvider(bpayConfig); err != nil {
		return fmt.Errorf("failed to add B-PAY provider: %w", err)
	}

	// Create B-PAY specific payment request
	phoneNum, err := phone.NewPhone("22123456")
	if err != nil {
		return fmt.Errorf(phoneCreationError, err)
	}
	amount := money.NewMRU(1000) // 10.00 MRU

	bpayRequest := &rimpay.BPayPaymentRequest{
		PhoneNumber: phoneNum,
		Amount:      amount,
		Description: "Test B-PAY payment",
		Reference:   "BPAY-TEST-001",
		Passcode:    "1234", // B-PAY specific field
		Metadata: map[string]interface{}{
			"customer_id": "12345",
			"order_id":    "ORD-001",
		},
	}

	// Process B-PAY payment with type safety
	ctx := context.Background()
	response, err := client.ProcessBPayPayment(ctx, bpayRequest)
	if err != nil {
		return fmt.Errorf("B-PAY payment failed: %w", err)
	}

	fmt.Printf("B-PAY Payment Response:\n")
	fmt.Printf("  Transaction ID: %s\n", response.TransactionID)
	fmt.Printf("  Status: %s\n", response.Status)
	fmt.Printf("  Reference: %s\n", response.Reference)
	fmt.Printf("  Provider: %s\n", response.Provider)
	fmt.Println()

	return nil
}

func runMasrviExample() error {
	fmt.Println("--- MASRVI Specific Payment Example ---")

	// Create client
	client, err := createTestClient()
	if err != nil {
		return err
	}

	

	// Add MASRVI provider
	masrviConfig := rimpay.ProviderConfig{
		Enabled: true,
		BaseURL: "https://api.masrvi.mr",
		Credentials: map[string]string{
			"merchant_id": "test_merchant",
		},
		Timeout: 30 * time.Second,
	}

	if err := client.AddMasrviProvider(masrviConfig); err != nil {
		return fmt.Errorf("failed to add MASRVI provider: %w", err)
	}

	// Create MASRVI specific payment request
	phoneNum, err := phone.NewPhone("33987654")
	if err != nil {
		return fmt.Errorf(phoneCreationError, err)
	}
	amount := money.NewMRU(500) // 5.00 MRU

	masrviRequest := &rimpay.MasrviPaymentRequest{
		PhoneNumber: phoneNum,
		Amount:      amount,
		Description: "Test MASRVI payment",
		Reference:   "MASRVI-TEST-001",
		CallbackURL: "https://myapp.com/webhook/masrvi", // MASRVI specific
		ReturnURL:   "https://myapp.com/return",         // MASRVI specific
		Metadata: map[string]interface{}{
			"customer_id": "67890",
			"session_id":  "sess_abc123",
		},
	}

	// Process MASRVI payment with type safety
	ctx := context.Background()
	response, err := client.ProcessMasrviPayment(ctx, masrviRequest)
	if err != nil {
		return fmt.Errorf("MASRVI payment failed: %w", err)
	}

	fmt.Printf("MASRVI Payment Response:\n")
	fmt.Printf("  Transaction ID: %s\n", response.TransactionID)
	fmt.Printf("  Status: %s\n", response.Status)
	fmt.Printf("  Reference: %s\n", response.Reference)
	fmt.Printf("  Provider: %s\n", response.Provider)
	
	// MASRVI specific response data
	if paymentURL, ok := response.Metadata["payment_url"].(string); ok {
		fmt.Printf("  Payment URL: %s\n", paymentURL)
	}
	fmt.Println()

	return nil
}

func runMultiProviderExample() error {
	fmt.Println("--- Multi-Provider Type-Safe Example ---")

	// Create client with both providers
	client, err := createTestClient()
	if err != nil {
		return err
	}

	// Add both providers
	bpayConfig := rimpay.ProviderConfig{
		Enabled: true,
		BaseURL: "https://api.bpay.mr",
		Credentials: map[string]string{
			"username":  "test_username",
			"password":  "test_password",
			"client_id": "test_client_id",
		},
		Timeout: 30 * time.Second,
	}

	masrviConfig := rimpay.ProviderConfig{
		Enabled: true,
		BaseURL: "https://api.masrvi.mr",
		Credentials: map[string]string{
			"merchant_id": "test_merchant",
		},
		Timeout: 30 * time.Second,
	}

	if err := client.AddBPayProvider(bpayConfig); err != nil {
		return fmt.Errorf("failed to add B-PAY provider: %w", err)
	}

	if err := client.AddMasrviProvider(masrviConfig); err != nil {
		return fmt.Errorf("failed to add MASRVI provider: %w", err)
	}

	// Process payment with B-PAY
	phoneNum1, err := phone.NewPhone("22111111")
	if err != nil {
		return fmt.Errorf(phoneCreationError, err)
	}
	amount1 := money.NewMRU(2000)

	bpayRequest := &rimpay.BPayPaymentRequest{
		PhoneNumber: phoneNum1,
		Amount:      amount1,
		Description: "Multi-provider test - B-PAY",
		Reference:   "MULTI-BPAY-001",
		Passcode:    "9999",
	}

	ctx := context.Background()
	bpayResponse, err := client.ProcessBPayPayment(ctx, bpayRequest)
	if err != nil {
		return fmt.Errorf("multi-provider B-PAY payment failed: %w", err)
	}

	fmt.Printf("Multi-Provider B-PAY Result: %s (Status: %s)\n", 
		bpayResponse.TransactionID, bpayResponse.Status)

	// Process payment with MASRVI
	phoneNum2, err := phone.NewPhone("33222222")
	if err != nil {
		return fmt.Errorf(phoneCreationError, err)
	}
	amount2 := money.NewMRU(1500)

	masrviRequest := &rimpay.MasrviPaymentRequest{
		PhoneNumber: phoneNum2,
		Amount:      amount2,
		Description: "Multi-provider test - MASRVI",
		Reference:   "MULTI-MASRVI-001",
		CallbackURL: "https://myapp.com/webhook",
		ReturnURL:   "https://myapp.com/return",
	}

	masrviResponse, err := client.ProcessMasrviPayment(ctx, masrviRequest)
	if err != nil {
		return fmt.Errorf("multi-provider MASRVI payment failed: %w", err)
	}

	fmt.Printf("Multi-Provider MASRVI Result: %s (Status: %s)\n", 
		masrviResponse.TransactionID, masrviResponse.Status)

	// Demonstrate type safety - these would be compile-time errors:
	// client.ProcessBPayPayment(ctx, masrviRequest)   // ❌ Type mismatch
	// client.ProcessMasrviPayment(ctx, bpayRequest)   // ❌ Type mismatch

	fmt.Println("✅ Multi-provider payments completed with type safety!")
	fmt.Println()

	return nil
}

func createTestClient() (*rimpay.Client, error) {
	config := &rimpay.Config{
		Environment:     rimpay.EnvironmentSandbox,
		DefaultProvider: "bpay",
		Providers: map[string]rimpay.ProviderConfig{
			"bpay": {
				Enabled: true,
				BaseURL: "https://ebankily-tst.appspot.com",
				Timeout: 30 * time.Second,
				Credentials: map[string]string{
					"username":  "test_username",
					"password":  "test_password",
					"client_id": "test_client_id",
				},
			},
			"masrvi": {
				Enabled: true,
				BaseURL: "https://masrviapp.mr/online",
				Timeout: 30 * time.Second,
				Credentials: map[string]string{
					"merchant_id": "test_merchant_id",
				},
			},
		},
	}

	return rimpay.NewClient(config)
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
