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
)

func main() {
	fmt.Println("🏦 RimPay Library - Multi-Provider Example")
	fmt.Println("==========================================")

	// Create configuration with multiple providers
	config := createMultiProviderConfig()
	client, err := rimpay.NewClient(config)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()

	// Example 1: B-PAY Payment
	fmt.Println("\n📱 Example 1: B-PAY Payment Processing")
	processBPayPayment(client, ctx)

	// Example 2: MASRVI Payment
	fmt.Println("\n🌐 Example 2: MASRVI Payment Processing")
	processMasrviPayment(client, ctx)

	// Example 3: Provider Selection Logic
	fmt.Println("\n🔄 Example 3: Dynamic Provider Selection")
	demonstrateProviderSelection(client, ctx)

	fmt.Println("\n💡 Multi-Provider Features Demonstrated:")
	fmt.Println("✅ Multiple payment providers in one configuration")
	fmt.Println("✅ Provider-specific request types")
	fmt.Println("✅ Type-safe payment methods")
	fmt.Println("✅ Dynamic provider selection")
	fmt.Println("✅ Unified error handling")
}

func createMultiProviderConfig() *rimpay.Config {
	return &rimpay.Config{
		Environment:     rimpay.EnvironmentSandbox,
		DefaultProvider: "bpay",
		Providers: map[string]rimpay.ProviderConfig{
			"bpay": {
				Enabled: true,
				BaseURL: "https://ebankily-tst.appspot.com",
				Timeout: 30 * time.Second,
				Credentials: map[string]string{
					"username":  "bpay_merchant",
					"password":  "bpay_password",
					"client_id": "bpay_client_123",
				},
			},
			"masrvi": {
				Enabled: true,
				BaseURL: "https://masrviapp.mr/online",
				Timeout: 30 * time.Second,
				Credentials: map[string]string{
					"merchant_id": "masrvi_merchant_456",
				},
			},
		},
	}
}

func processBPayPayment(client *rimpay.Client, ctx context.Context) {
	phone, _ := phone.NewPhone("22334455")
	amount := money.New(decimal.NewFromFloat(75.00), money.MRU)

	request := &rimpay.BPayPaymentRequest{
		Amount:      amount,
		PhoneNumber: phone,
		Reference:   fmt.Sprintf("BPAY-MULTI-%d", time.Now().Unix()),
		Description: "Multi-provider B-PAY payment",
		Passcode:    "1234",
	}

	fmt.Printf("   Processing B-PAY payment...\n")
	fmt.Printf("   Amount: %s\n", request.Amount.String())
	fmt.Printf("   Phone: %s\n", request.PhoneNumber.String())

	response, err := client.ProcessBPayPayment(ctx, request)
	if err != nil {
		fmt.Printf("   ❌ B-PAY payment failed: %v\n", err)
		return
	}

	fmt.Printf("   ✅ B-PAY payment successful!\n")
	fmt.Printf("   Transaction ID: %s\n", response.TransactionID)
	fmt.Printf("   Provider: %s\n", response.Provider)
}

func processMasrviPayment(client *rimpay.Client, ctx context.Context) {
	phone, _ := phone.NewPhone("33445566")
	amount := money.New(decimal.NewFromFloat(125.50), money.MRU)

	request := &rimpay.MasrviPaymentRequest{
		Amount:      amount,
		PhoneNumber: phone,
		Reference:   fmt.Sprintf("MASRVI-MULTI-%d", time.Now().Unix()),
		Description: "Multi-provider MASRVI payment",
		CallbackURL: "https://yourapp.com/webhook/masrvi",
		ReturnURL:   "https://yourapp.com/return",
	}

	fmt.Printf("   Processing MASRVI payment...\n")
	fmt.Printf("   Amount: %s\n", request.Amount.String())
	fmt.Printf("   Phone: %s\n", request.PhoneNumber.String())

	response, err := client.ProcessMasrviPayment(ctx, request)
	if err != nil {
		fmt.Printf("   ❌ MASRVI payment failed: %v\n", err)
		return
	}

	fmt.Printf("   ✅ MASRVI payment form created!\n")
	fmt.Printf("   Transaction ID: %s\n", response.TransactionID)
	fmt.Printf("   Provider: %s\n", response.Provider)
	if paymentURL, exists := response.Metadata["payment_url"]; exists {
		fmt.Printf("   Payment URL: %s\n", paymentURL)
	}
}

func demonstrateProviderSelection(client *rimpay.Client, ctx context.Context) {
	// Example: Choose provider based on amount
	testAmounts := []float64{10.00, 100.00, 500.00}

	for _, amount := range testAmounts {
		provider := selectProvider(amount)
		fmt.Printf("   Amount: %.2f MRU → Recommended provider: %s\n", amount, provider)

		phone, _ := phone.NewPhone("22334455")
		money := money.New(decimal.NewFromFloat(amount), money.MRU)

		switch provider {
		case "bpay":
			request := &rimpay.BPayPaymentRequest{
				Amount:      money,
				PhoneNumber: phone,
				Reference:   fmt.Sprintf("DYNAMIC-BPAY-%d", time.Now().Unix()),
				Description: "Dynamic provider selection - B-PAY",
				Passcode:    "1234",
			}

			fmt.Printf("   → Processing with B-PAY...\n")
			_, err := client.ProcessBPayPayment(ctx, request)
			if err != nil {
				fmt.Printf("   ❌ Failed: %v\n", err)
			} else {
				fmt.Printf("   ✅ B-PAY payment initiated\n")
			}

		case "masrvi":
			request := &rimpay.MasrviPaymentRequest{
				Amount:      money,
				PhoneNumber: phone,
				Reference:   fmt.Sprintf("DYNAMIC-MASRVI-%d", time.Now().Unix()),
				Description: "Dynamic provider selection - MASRVI",
				CallbackURL: "https://yourapp.com/webhook",
				ReturnURL:   "https://yourapp.com/return",
			}

			fmt.Printf("   → Processing with MASRVI...\n")
			_, err := client.ProcessMasrviPayment(ctx, request)
			if err != nil {
				fmt.Printf("   ❌ Failed: %v\n", err)
			} else {
				fmt.Printf("   ✅ MASRVI payment form created\n")
			}
		}

		fmt.Println()
	}
}

func selectProvider(amount float64) string {
	// Business logic for provider selection
	switch {
	case amount < 50.00:
		// Small amounts - use B-PAY for instant processing
		return "bpay"
	case amount < 200.00:
		// Medium amounts - either provider works
		return "bpay"
	default:
		// Large amounts - use MASRVI for web-based confirmation
		return "masrvi"
	}
}
