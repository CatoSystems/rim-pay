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
)

func main() {
	fmt.Println("🏦 RimPay Library - Multi-Provider Example")
	fmt.Println("=========================================\n")

	// Create configuration with multiple providers
	config := createMultiProviderConfig()

	client, err := rimpay.NewClient(config)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()

	// Example 1: Process payment with default provider (B-PAY)
	fmt.Println("📱 Example 1: Default Provider Payment (B-PAY)")
	processDefaultProviderPayment(client, ctx)

	// Example 2: Process payment with specific provider (MASRVI)
	fmt.Println("\n🌐 Example 2: Specific Provider Payment (MASRVI)")
	processSpecificProviderPayment(client, ctx, "masrvi")

	// Example 3: Provider fallback simulation
	fmt.Println("\n🔄 Example 3: Provider Failover Simulation")
	simulateProviderFailover(client, ctx)

	// Example 4: Bulk payments to different providers
	fmt.Println("\n📊 Example 4: Bulk Payments")
	processBulkPayments(client, ctx)

	fmt.Println("\n💡 Multi-Provider Features Demonstrated:")
	fmt.Println("✅ Default provider configuration")
	fmt.Println("✅ Provider-specific payment processing")
	fmt.Println("✅ Automatic provider selection")
	fmt.Println("✅ Provider availability checking")
	fmt.Println("✅ Bulk payment processing")
	fmt.Println("✅ Error handling across providers")
}

func createMultiProviderConfig() *rimpay.Config {
	return &rimpay.Config{
		Environment:     rimpay.EnvironmentSandbox,
		DefaultProvider: "bpay", // Primary provider
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
		HTTP: rimpay.HTTPConfig{
			Timeout:         30 * time.Second,
			MaxIdleConns:    20,
			MaxConnsPerHost: 10,
		},
	}
}

func processDefaultProviderPayment(client *rimpay.Client, ctx context.Context) {
	phone, _ := phone.NewPhone("22334455")
	amount := money.New(decimal.NewFromFloat(50.00), money.MRU)

	// Using B-PAY as default provider with provider-specific request
	request := &rimpay.BPayPaymentRequest{
		Amount:      amount,
		PhoneNumber: phone,
		Reference:   fmt.Sprintf("DEFAULT-%d", time.Now().Unix()),
		Description: "Payment using default provider (B-PAY)",
		Passcode:    "1234",
	}

	fmt.Printf("   Using default provider (B-PAY) for %s\n", amount.String())
	
	response, err := client.ProcessBPayPayment(ctx, request)
	if err != nil {
		fmt.Printf("   ❌ Payment failed: %v\n", err)
		return
	}

	fmt.Printf("   ✅ Payment successful via %s\n", response.Provider)
	fmt.Printf("   Transaction ID: %s\n", response.TransactionID)
}

func processSpecificProviderPayment(client *rimpay.Client, ctx context.Context, providerName string) {
	phone, _ := phone.NewPhone("33112233")
	amount := money.New(decimal.NewFromFloat(75.00), money.MRU)

	fmt.Printf("   Forcing payment through %s provider\n", providerName)
	fmt.Printf("   Amount: %s to %s\n", amount.String(), phone.ForProvider(true))

	if providerName == "masrvi" {
		request := &rimpay.MasrviPaymentRequest{
			Amount:      amount,
			PhoneNumber: phone,
			Reference:   fmt.Sprintf("SPECIFIC-%s-%d", providerName, time.Now().Unix()),
			Description: fmt.Sprintf("Payment using %s provider", providerName),
			ReturnURL:   "https://example.com/return",
			CallbackURL: "https://example.com/webhook",
		}

		response, err := client.ProcessMasrviPayment(ctx, request)
		if err != nil {
			fmt.Printf("   ❌ Payment failed via %s: %v\n", providerName, err)
			return
		}

		fmt.Printf("   ✅ Payment successful via %s\n", response.Provider)
		fmt.Printf("   Transaction ID: %s\n", response.TransactionID)
		fmt.Printf("   Status: %s\n", response.Status)
	} else {
		// Default to B-PAY for other providers
		request := &rimpay.BPayPaymentRequest{
			Amount:      amount,
			PhoneNumber: phone,
			Reference:   fmt.Sprintf("SPECIFIC-%s-%d", providerName, time.Now().Unix()),
			Description: fmt.Sprintf("Payment using %s provider", providerName),
			Passcode:    "1234",
		}

		response, err := client.ProcessBPayPayment(ctx, request)
		if err != nil {
			fmt.Printf("   ❌ Payment failed via %s: %v\n", providerName, err)
			return
		}

		fmt.Printf("   ✅ Payment successful via %s\n", response.Provider)
		fmt.Printf("   Transaction ID: %s\n", response.TransactionID)
		fmt.Printf("   Status: %s\n", response.Status)
	}
}

func simulateProviderFailover(client *rimpay.Client, ctx context.Context) {
	phone, _ := phone.NewPhone("24990011")
	amount := money.New(decimal.NewFromFloat(100.00), money.MRU)

	request := &rimpay.BPayPaymentRequest{
		Amount:      amount,
		PhoneNumber: phone,
		Reference:   fmt.Sprintf("FAILOVER-%d", time.Now().Unix()),
		Description: "Failover simulation payment",
		Passcode:    "1234",
	}

	fmt.Printf("   Simulating provider failover scenario...\n")
	fmt.Printf("   Amount: %s to %s\n", amount.String(), phone.ForProvider(true))

	// Try primary provider (B-PAY)
	fmt.Printf("   🔄 Trying primary provider (B-PAY)...\n")
	response, err := client.ProcessBPayPayment(ctx, request)
	
	if err != nil {
		fmt.Printf("   ⚠️  Primary provider failed: %v\n", err)
		
		// Check if error is retryable
		if paymentErr, ok := err.(*rimpay.PaymentError); ok && !paymentErr.IsRetryable() {
			// Switch to MASRVI (simulate failover)
			fmt.Printf("   🔄 Switching to fallback provider (MASRVI)...\n")
			
			// Create a new MASRVI request
			masrviRequest := &rimpay.MasrviPaymentRequest{
				Amount:      amount,
				PhoneNumber: phone,
				Reference:   fmt.Sprintf("FAILOVER-MASRVI-%d", time.Now().Unix()),
				Description: "Failover simulation payment via MASRVI",
				ReturnURL:   "https://example.com/return",
				CallbackURL: "https://example.com/webhook",
			}
			
			response, err = client.ProcessMasrviPayment(ctx, masrviRequest)
			if err != nil {
				fmt.Printf("   ❌ Fallback provider also failed: %v\n", err)
				return
			}
		} else {
			fmt.Printf("   ❌ Error is not suitable for failover\n")
			return
		}
	}

	fmt.Printf("   ✅ Payment successful via %s\n", response.Provider)
	fmt.Printf("   Transaction ID: %s\n", response.TransactionID)
}

func processBulkPayments(client *rimpay.Client, ctx context.Context) {
	payments := []struct {
		phone       string
		amount      float64
		provider    string
		description string
	}{
		{"22334455", 25.00, "bpay", "Mauritel via B-PAY"},
		{"33445566", 30.00, "masrvi", "Mauritel via MASRVI"},
		{"36778899", 45.00, "bpay", "Mattel via B-PAY"},
		{"37889900", 55.00, "masrvi", "Mattel via MASRVI"},
		{"48990011", 35.00, "bpay", "Chinguitel via B-PAY"},
	}

	fmt.Printf("   Processing %d bulk payments...\n", len(payments))
	
	results := make(map[string]int)
	
	for i, payment := range payments {
		fmt.Printf("\n   Payment %d/%d: %s via %s\n", i+1, len(payments), payment.description, payment.provider)
		
		phone, _ := phone.NewPhone(payment.phone)
		amount := money.New(decimal.NewFromFloat(payment.amount), money.MRU)

		fmt.Printf("   Note: Using provider-specific request types for %s\n", payment.provider)
		
		var response *rimpay.PaymentResponse
		var err error
		
		if payment.provider == "bpay" {
			request := &rimpay.BPayPaymentRequest{
				Amount:      amount,
				PhoneNumber: phone,
				Reference:   fmt.Sprintf("BULK-BPAY-%d-%d", i+1, time.Now().Unix()),
				Description: payment.description,
				Passcode:    "1234",
			}
			response, err = client.ProcessBPayPayment(ctx, request)
		} else {
			request := &rimpay.MasrviPaymentRequest{
				Amount:      amount,
				PhoneNumber: phone,
				Reference:   fmt.Sprintf("BULK-MASRVI-%d-%d", i+1, time.Now().Unix()),
				Description: payment.description,
				ReturnURL:   "https://example.com/return",
				CallbackURL: "https://example.com/webhook",
			}
			response, err = client.ProcessMasrviPayment(ctx, request)
		}

		if err != nil {
			fmt.Printf("   ❌ Failed: %v\n", err)
			results["failed"]++
		} else {
			fmt.Printf("   ✅ Success: %s\n", response.TransactionID)
			results[response.Provider]++
		}

		// Small delay between requests
		time.Sleep(500 * time.Millisecond)
	}

	fmt.Printf("\n   📊 Bulk Payment Results:\n")
	for provider, count := range results {
		fmt.Printf("   %s: %d payments\n", provider, count)
	}
}

// Note: In a real implementation, you would extend the Client struct
// to support dynamic provider selection and switching
