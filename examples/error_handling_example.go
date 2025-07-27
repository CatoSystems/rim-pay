package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/CatoSystems/rim-pay/pkg/money"
	"github.com/CatoSystems/rim-pay/pkg/phone"
	"github.com/CatoSystems/rim-pay/pkg/rimpay"
	"github.com/shopspring/decimal"

	// Import providers
	_ "github.com/CatoSystems/rim-pay/internal/providers/bpay"
	_ "github.com/CatoSystems/rim-pay/internal/providers/masrvi"
)

func main() {
	fmt.Println("🏦 RimPay Library - Error Handling & Retry Example")
	fmt.Println("=================================================\n")

	config := createTestConfig()
	client, err := rimpay.NewClient(config)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()

	// Example 1: Network error with automatic retry
	fmt.Println("🔄 Example 1: Network Error Handling")
	demonstrateNetworkErrorRetry(client, ctx)

	// Example 2: Authentication error handling
	fmt.Println("\n🔐 Example 2: Authentication Error")
	demonstrateAuthError(client, ctx)

	// Example 3: Validation error (non-retryable)
	fmt.Println("\n✅ Example 3: Validation Error")
	demonstrateValidationError(client, ctx)

	// Example 4: Business logic error handling
	fmt.Println("\n💳 Example 4: Business Logic Errors")
	demonstrateBusinessErrors(client, ctx)

	// Example 5: Context timeout handling
	fmt.Println("\n⏰ Example 5: Timeout Handling")
	demonstrateTimeoutHandling(client, ctx)

	fmt.Println("\n💡 Error Handling Features Demonstrated:")
	fmt.Println("✅ Automatic retry for transient failures")
	fmt.Println("✅ Exponential backoff with jitter")
	fmt.Println("✅ Context timeout handling")
	fmt.Println("✅ Detailed error classification")
	fmt.Println("✅ Provider-specific error mapping")
	fmt.Println("✅ Graceful degradation strategies")
}

func createTestConfig() *rimpay.Config {
	return &rimpay.Config{
		Environment:     rimpay.EnvironmentSandbox,
		DefaultProvider: "bpay",
		Providers: map[string]rimpay.ProviderConfig{
			"bpay": {
				Enabled: true,
				BaseURL: "https://ebankily-tst.appspot.com",
				Timeout: 30 * time.Second,
				Credentials: map[string]string{
					"username":  "test_user",
					"password":  "test_pass",
					"client_id": "test_client",
				},
			},
		},
	}
}

func demonstrateNetworkErrorRetry(client *rimpay.Client, ctx context.Context) {
	phone, _ := phone.NewPhone("22334455")
	amount := money.New(decimal.NewFromFloat(50.00), money.MRU)

	request := &rimpay.PaymentRequest{
		Amount:      amount,
		PhoneNumber: phone,
		Reference:   fmt.Sprintf("RETRY-TEST-%d", time.Now().Unix()),
		Language:    rimpay.LanguageFrench,
		Passcode:    "1234",
		Description: "Network retry test payment",
	}

	fmt.Printf("   Simulating payment with potential network issues...\n")
	fmt.Printf("   The library will automatically retry on network failures\n")
	fmt.Printf("   Retry configuration: 3 attempts, exponential backoff\n\n")

	start := time.Now()
	response, err := client.ProcessPayment(ctx, request)
	duration := time.Since(start)

	if err != nil {
		fmt.Printf("   ❌ Payment failed after retries: %v\n", err)
		fmt.Printf("   Duration: %v\n", duration)
		
		if paymentErr, ok := err.(*rimpay.PaymentError); ok {
			fmt.Printf("   Error details:\n")
			fmt.Printf("     Code: %s\n", paymentErr.Code)
			fmt.Printf("     Provider: %s\n", paymentErr.Provider)
			fmt.Printf("     Retryable: %v\n", paymentErr.IsRetryable())
			
			if paymentErr.IsRetryable() {
				fmt.Printf("   💡 This error was retried automatically\n")
				fmt.Printf("   💡 The library attempted up to 3 times with exponential backoff\n")
			}
		}
	} else {
		fmt.Printf("   ✅ Payment successful: %s\n", response.TransactionID)
		fmt.Printf("   Duration: %v\n", duration)
	}
}

func demonstrateAuthError(client *rimpay.Client, ctx context.Context) {
	// Create a client with invalid credentials to simulate auth error
	badConfig := &rimpay.Config{
		Environment:     rimpay.EnvironmentSandbox,
		DefaultProvider: "bpay",
		Providers: map[string]rimpay.ProviderConfig{
			"bpay": {
				Enabled: true,
				BaseURL: "https://ebankily-tst.appspot.com",
				Timeout: 10 * time.Second,
				Credentials: map[string]string{
					"username":  "invalid_user",
					"password":  "wrong_password",
					"client_id": "bad_client_id",
				},
			},
		},
	}

	badClient, err := rimpay.NewClient(badConfig)
	if err != nil {
		fmt.Printf("   ❌ Failed to create client with bad credentials: %v\n", err)
		return
	}

	phone, _ := phone.NewPhone("22334455")
	amount := money.New(decimal.NewFromFloat(25.00), money.MRU)

	request := &rimpay.PaymentRequest{
		Amount:      amount,
		PhoneNumber: phone,
		Reference:   fmt.Sprintf("AUTH-TEST-%d", time.Now().Unix()),
		Language:    rimpay.LanguageFrench,
		Passcode:    "1234",
		Description: "Authentication error test",
	}

	fmt.Printf("   Testing with invalid credentials...\n")
	
	_, err = badClient.ProcessPayment(ctx, request)
	if err != nil {
		if paymentErr, ok := err.(*rimpay.PaymentError); ok {
			fmt.Printf("   ❌ Authentication failed as expected\n")
			fmt.Printf("   Error code: %s\n", paymentErr.Code)
			fmt.Printf("   Retryable: %v\n", paymentErr.IsRetryable())
			
			if paymentErr.Code == rimpay.ErrorCodeAuthenticationFailed {
				fmt.Printf("   💡 This is an authentication error\n")
				fmt.Printf("   💡 Check your provider credentials\n")
				if paymentErr.IsRetryable() {
					fmt.Printf("   💡 The library will retry auth failures (token might have expired)\n")
				}
			}
		}
	}
}

func demonstrateValidationError(client *rimpay.Client, ctx context.Context) {
	fmt.Printf("   Testing various validation scenarios...\n\n")

	// Test 1: Invalid phone number
	fmt.Printf("   Test 1: Invalid phone number\n")
	testInvalidPhone(client, ctx)

	// Test 2: Zero amount
	fmt.Printf("\n   Test 2: Zero amount\n")
	testZeroAmount(client, ctx)

	// Test 3: Missing reference
	fmt.Printf("\n   Test 3: Missing reference\n")
	testMissingReference(client, ctx)

	// Test 4: Invalid URL
	fmt.Printf("\n   Test 4: Invalid callback URL\n")
	testInvalidURL(client, ctx)
}

func testInvalidPhone(client *rimpay.Client, ctx context.Context) {
	// This will fail at phone creation level
	_, err := phone.NewPhone("invalid-phone")
	if err != nil {
		fmt.Printf("     ❌ Phone validation failed: %v\n", err)
		fmt.Printf("     💡 Phone numbers must be valid Mauritanian numbers\n")
	}
}

func testZeroAmount(client *rimpay.Client, ctx context.Context) {
	phone, _ := phone.NewPhone("22334455")
	amount := money.New(decimal.Zero, money.MRU)

	request := &rimpay.PaymentRequest{
		Amount:      amount,
		PhoneNumber: phone,
		Reference:   "ZERO-AMOUNT-TEST",
		Language:    rimpay.LanguageFrench,
		Passcode:    "1234",
	}

	_, err := client.ProcessPayment(ctx, request)
	if err != nil {
		if paymentErr, ok := err.(*rimpay.PaymentError); ok {
			fmt.Printf("     ❌ Zero amount rejected: %s\n", paymentErr.Message)
			fmt.Printf("     💡 Amount must be greater than zero\n")
		}
	}
}

func testMissingReference(client *rimpay.Client, ctx context.Context) {
	phone, _ := phone.NewPhone("22334455")
	amount := money.New(decimal.NewFromFloat(50.00), money.MRU)

	request := &rimpay.PaymentRequest{
		Amount:      amount,
		PhoneNumber: phone,
		Reference:   "", // Empty reference
		Language:    rimpay.LanguageFrench,
		Passcode:    "1234",
	}

	_, err := client.ProcessPayment(ctx, request)
	if err != nil {
		if paymentErr, ok := err.(*rimpay.PaymentError); ok {
			fmt.Printf("     ❌ Missing reference rejected: %s\n", paymentErr.Message)
			fmt.Printf("     💡 Reference is required for tracking\n")
		}
	}
}

func testInvalidURL(client *rimpay.Client, ctx context.Context) {
	phone, _ := phone.NewPhone("22334455")
	amount := money.New(decimal.NewFromFloat(50.00), money.MRU)

	request := &rimpay.PaymentRequest{
		Amount:      amount,
		PhoneNumber: phone,
		Reference:   "URL-TEST",
		Language:    rimpay.LanguageFrench,
		Passcode:    "1234",
		CallbackURL: "invalid-url", // Invalid URL format
	}

	_, err := client.ProcessPayment(ctx, request)
	if err != nil {
		if paymentErr, ok := err.(*rimpay.PaymentError); ok {
			fmt.Printf("     ❌ Invalid URL rejected: %s\n", paymentErr.Message)
			fmt.Printf("     💡 URLs must be valid format\n")
		}
	}
}

func demonstrateBusinessErrors(client *rimpay.Client, ctx context.Context) {
	fmt.Printf("   Simulating common business logic errors...\n")

	phone, _ := phone.NewPhone("22334455")
	amount := money.New(decimal.NewFromFloat(50.00), money.MRU)

	// Test insufficient funds scenario
	request := &rimpay.PaymentRequest{
		Amount:      amount,
		PhoneNumber: phone,
		Reference:   fmt.Sprintf("BUSINESS-TEST-%d", time.Now().Unix()),
		Language:    rimpay.LanguageFrench,
		Passcode:    "0000", // Wrong passcode to simulate business error
		Description: "Business error test",
	}

	_, err := client.ProcessPayment(ctx, request)
	if err != nil {
		if paymentErr, ok := err.(*rimpay.PaymentError); ok {
			fmt.Printf("   ❌ Business error: %s\n", paymentErr.Message)
			
			switch paymentErr.Code {
			case rimpay.ErrorCodeInsufficientFunds:
				fmt.Printf("   💡 Customer needs to add funds to their mobile money account\n")
			case rimpay.ErrorCodePaymentDeclined:
				fmt.Printf("   💡 Customer entered wrong PIN - ask them to retry\n")
			case rimpay.ErrorCodeProviderError:
				fmt.Printf("   💡 Customer's account may be blocked - contact provider support\n")
			default:
				fmt.Printf("   💡 Check error code %s for specific handling\n", paymentErr.Code)
			}
			
			fmt.Printf("   Retryable: %v\n", paymentErr.IsRetryable())
		}
	}
}

func demonstrateTimeoutHandling(client *rimpay.Client, ctx context.Context) {
	// Create context with short timeout
	shortCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	phone, _ := phone.NewPhone("22334455")
	amount := money.New(decimal.NewFromFloat(50.00), money.MRU)

	request := &rimpay.PaymentRequest{
		Amount:      amount,
		PhoneNumber: phone,
		Reference:   fmt.Sprintf("TIMEOUT-TEST-%d", time.Now().Unix()),
		Language:    rimpay.LanguageFrench,
		Passcode:    "1234",
		Description: "Timeout test payment",
	}

	fmt.Printf("   Testing with 1-second timeout...\n")
	
	start := time.Now()
	_, err := client.ProcessPayment(shortCtx, request)
	duration := time.Since(start)

	if err != nil {
		fmt.Printf("   ❌ Payment failed due to timeout: %v\n", err)
		fmt.Printf("   Duration: %v\n", duration)
		
		if err == context.DeadlineExceeded {
			fmt.Printf("   💡 This was a context timeout\n")
			fmt.Printf("   💡 Consider increasing timeout for production\n")
			fmt.Printf("   💡 Recommended timeout: 30-60 seconds\n")
		}
	} else {
		fmt.Printf("   ✅ Payment completed within timeout\n")
		fmt.Printf("   Duration: %v\n", duration)
	}
}
