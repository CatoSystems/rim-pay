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

	// Import providers to register them
	_ "github.com/CatoSystems/rim-pay/internal/providers/bpay"
	_ "github.com/CatoSystems/rim-pay/internal/providers/masrvi"
)

func main() {
	fmt.Println("🏦 RimPay Library - Basic Usage Example")
	fmt.Println("=====================================\n")

	// Step 1: Create configuration
	config := createConfig()

	// Step 2: Initialize client
	client, err := rimpay.NewClient(config)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Step 3: Create payment request
	request, err := createPaymentRequest()
	if err != nil {
		log.Fatalf("Failed to create payment request: %v", err)
	}

	// Step 4: Process payment
	fmt.Printf("📱 Processing payment...\n")
	fmt.Printf("   Amount: %s\n", request.Amount.String())
	fmt.Printf("   Phone: %s\n", request.PhoneNumber.ForProvider(true))
	fmt.Printf("   Reference: %s\n", request.Reference)
	fmt.Printf("   Provider: %s\n\n", config.DefaultProvider)

	ctx := context.Background()
	response, err := client.ProcessPayment(ctx, request)
	
	if err != nil {
		handlePaymentError(err)
		return
	}

	// Step 5: Handle successful payment
	fmt.Printf("✅ Payment successful!\n")
	fmt.Printf("   Transaction ID: %s\n", response.TransactionID)
	fmt.Printf("   Status: %s\n", response.Status)
	fmt.Printf("   Provider: %s\n", response.Provider)
	fmt.Printf("   Created: %s\n\n", response.CreatedAt.Format(time.RFC3339))

	// Step 6: Check payment status (for B-PAY)
	if config.DefaultProvider == "bpay" {
		fmt.Printf("🔍 Checking payment status...\n")
		status, err := client.GetPaymentStatus(ctx, response.TransactionID)
		if err != nil {
			fmt.Printf("❌ Failed to get status: %v\n", err)
		} else {
			fmt.Printf("   Status: %s\n", status.Status)
			fmt.Printf("   Message: %s\n\n", status.Message)
		}
	}

	fmt.Println("🎉 Example completed successfully!")
}

func createConfig() *rimpay.Config {
	return &rimpay.Config{
		Environment:     rimpay.EnvironmentSandbox,
		DefaultProvider: "bpay", // Change to "masrvi" to test MASRVI
		Providers: map[string]rimpay.ProviderConfig{
			"bpay": {
				Enabled: true,
				BaseURL: "https://ebankily-tst.appspot.com",
				Timeout: 30 * time.Second,
				Credentials: map[string]string{
					"username":  "your_bpay_username",
					"password":  "your_bpay_password",
					"client_id": "your_bpay_client_id",
				},
			},
			"masrvi": {
				Enabled: true,
				BaseURL: "https://masrviapp.mr/online",
				Timeout: 30 * time.Second,
				Credentials: map[string]string{
					"merchant_id": "your_masrvi_merchant_id",
				},
			},
		},
		HTTP: rimpay.HTTPConfig{
			Timeout:         30 * time.Second,
			MaxIdleConns:    10,
			MaxConnsPerHost: 5,
		},
		Logging: rimpay.LoggingConfig{
			Level:  "info",
			Format: "json",
		},
	}
}

func createPaymentRequest() (*rimpay.PaymentRequest, error) {
	// Create phone number
	phoneNumber, err := phone.NewPhone("22334455") // Mauritel number
	if err != nil {
		return nil, fmt.Errorf("invalid phone number: %w", err)
	}

	// Create amount (100 MRU)
	amount := money.New(decimal.NewFromFloat(100.00), money.MRU)

	// Create payment request
	return &rimpay.PaymentRequest{
		Amount:      amount,
		PhoneNumber: phoneNumber,
		Reference:   fmt.Sprintf("ORDER-%d", time.Now().Unix()),
		Language:    rimpay.LanguageFrench,
		Passcode:    "1234", // For B-PAY
		Description: "Test payment via RimPay",
		// Optional URLs for MASRVI
		SuccessURL:  "https://yourapp.com/success",
		FailureURL:  "https://yourapp.com/failure",
		CancelURL:   "https://yourapp.com/cancel",
		CallbackURL: "https://yourapp.com/webhook",
	}, nil
}

func handlePaymentError(err error) {
	fmt.Printf("❌ Payment failed: %v\n", err)
	
	// Check if it's a payment error for more details
	if paymentErr, ok := err.(*rimpay.PaymentError); ok {
		fmt.Printf("   Error Code: %s\n", paymentErr.Code)
		fmt.Printf("   Provider: %s\n", paymentErr.Provider)
		fmt.Printf("   Retryable: %v\n", paymentErr.IsRetryable())
		
		// Suggest actions based on error type
		switch paymentErr.Code {
		case rimpay.ErrorCodeAuthenticationFailed:
			fmt.Printf("   💡 Check your provider credentials\n")
		case rimpay.ErrorCodeInsufficientFunds:
			fmt.Printf("   💡 Customer needs to add funds to their account\n")
		case rimpay.ErrorCodeInvalidRequest:
			fmt.Printf("   💡 Check phone number and amount format\n")
		case rimpay.ErrorCodeNetworkError:
			fmt.Printf("   💡 This error is retryable - the library will retry automatically\n")
		case rimpay.ErrorCodeProviderError:
			fmt.Printf("   💡 Provider service may be temporarily unavailable\n")
		}
	}
}
