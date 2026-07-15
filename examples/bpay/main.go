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
	fmt.Println("🏦 RimPay Library - B-PAY Provider Example")
	fmt.Println("=========================================")

	// Create B-PAY specific configuration
	config := &rimpay.Config{
		Environment:     rimpay.EnvironmentSandbox,
		DefaultProvider: "bpay",
		Providers: map[string]rimpay.ProviderConfig{
			"bpay": {
				Enabled: true,
				BaseURL: "https://ebankily-tst.appspot.com",
				Timeout: 30 * time.Second,
				Credentials: map[string]string{
					"username":  "test_merchant",     // Your B-PAY username
					"password":  "test_password",     // Your B-PAY password
					"client_id": "test_client_12345", // Your B-PAY client ID
				},
			},
		},
	}

	client, err := rimpay.NewClient(config)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Register the B-PAY provider instance on the client.
	if err := client.AddBPayProvider(config.Providers["bpay"]); err != nil {
		log.Fatalf("Failed to add B-PAY provider: %v", err)
	}

	fmt.Println("🔐 B-PAY Authentication Flow")
	fmt.Println("The library handles OAuth 2.0 authentication automatically:")
	fmt.Println("1. Requests access token using username/password/client_id")
	fmt.Println("2. Automatically refreshes tokens when they expire")
	fmt.Println("3. Handles token errors gracefully with retry")

	// Check if B-PAY provider is available
	ctx := context.Background()
	fmt.Printf("🔍 Checking B-PAY availability...\n")

	// This will attempt to authenticate and return true if successful
	if isAvailable := checkProviderAvailability(client, ctx); !isAvailable {
		fmt.Printf("❌ B-PAY provider is not available (check credentials)\n")
		return
	}

	fmt.Printf("✅ B-PAY provider is available\n\n")

	// Example 1: Mauritel payment
	fmt.Println("📱 Example 1: Mauritel Payment")
	mauritelPayment := createBPayPayment("22334455", 50.00, "Mauritel customer payment")
	processBPayPayment(client, ctx, mauritelPayment)

	// Example 2: Mattel payment
	fmt.Println("\n📱 Example 2: Mattel Payment")
	mattelPayment := createBPayPayment("32334455", 75.50, "Mattel customer payment")
	processBPayPayment(client, ctx, mattelPayment)

	// Example 3: Chinguitel payment
	fmt.Println("\n📱 Example 3: Chinguitel Payment")
	chinguitelPayment := createBPayPayment("44990011", 125.25, "Chinguitel customer payment")
	processBPayPayment(client, ctx, chinguitelPayment)

	fmt.Println("\n💡 B-PAY Features Demonstrated:")
	fmt.Println("✅ OAuth 2.0 authentication with automatic token refresh")
	fmt.Println("✅ Payment processing for all Mauritanian operators")
	fmt.Println("✅ Transaction status checking")
	fmt.Println("✅ Automatic retry on network/auth failures")
	fmt.Println("✅ Proper error handling and classification")
}

func checkProviderAvailability(client *rimpay.Client, ctx context.Context) bool {
	// Try to create a small test request to check authentication
	testPhone, _ := phone.NewPhone("22334455")
	testAmount := money.New(decimal.NewFromFloat(1.00), money.MRU)

	testRequest := &rimpay.BPayPaymentRequest{
		Amount:      testAmount,
		PhoneNumber: testPhone,
		Reference:   "AVAILABILITY-CHECK",
		Description: "Availability test",
		Passcode:    "0000", // This will fail, but we're testing auth
	}

	// This will fail at the payment step but succeed at authentication
	_, err := client.ProcessBPayPayment(ctx, testRequest)
	if err != nil {
		if paymentErr, ok := err.(*rimpay.PaymentError); ok {
			// If it's not an auth error, then auth worked
			return paymentErr.Code != rimpay.ErrorCodeAuthenticationFailed
		}
	}
	return true
}

func createBPayPayment(phoneNumber string, amount float64, description string) *rimpay.BPayPaymentRequest {
	phone, err := phone.NewPhone(phoneNumber)
	if err != nil {
		log.Fatalf("Invalid phone number: %v", err)
	}

	money := money.New(decimal.NewFromFloat(amount), money.MRU)

	return &rimpay.BPayPaymentRequest{
		Amount:      money,
		PhoneNumber: phone,
		Reference:   fmt.Sprintf("BPAY-%d", time.Now().UnixNano()),
		Description: description,
		// The 4-digit verification code Bankily gives the customer AFTER they
		// push the funds to the merchant account via the Bankily app. It is
		// collected from the customer and passed in here — the library never
		// generates it.
		Passcode: "1234",
	}
}

func processBPayPayment(client *rimpay.Client, ctx context.Context, request *rimpay.BPayPaymentRequest) {
	fmt.Printf("   Processing: %s → %s\n",
		request.PhoneNumber.ForProvider(true),
		request.Amount.String())
	fmt.Printf("   Reference: %s\n", request.Reference)

	// Process payment with automatic retry
	response, err := client.ProcessBPayPayment(ctx, request)
	if err != nil {
		fmt.Printf("   ❌ Payment failed: %v\n", err)

		if paymentErr, ok := err.(*rimpay.PaymentError); ok {
			switch paymentErr.Code {
			case rimpay.ErrorCodeInsufficientFunds:
				fmt.Printf("   💡 Customer needs to add funds\n")
			case rimpay.ErrorCodePaymentDeclined:
				fmt.Printf("   💡 Wrong or expired passcode\n")
			case rimpay.ErrorCodeNetworkError:
				fmt.Printf("   💡 Network issue - payment was retried automatically\n")
			}
		}
		return
	}

	fmt.Printf("   ✅ Payment successful!\n")
	fmt.Printf("   Transaction ID: %s\n", response.TransactionID)
	fmt.Printf("   Status: %s\n", response.Status)

	// Check final status
	time.Sleep(2 * time.Second) // Wait a bit before checking status
	checkPaymentStatus(client, ctx, response.TransactionID)
}

func checkPaymentStatus(client *rimpay.Client, ctx context.Context, transactionID string) {
	fmt.Printf("   🔍 Checking final status...\n")

	status, err := client.GetPaymentStatus(ctx, transactionID)
	if err != nil {
		fmt.Printf("   ❌ Status check failed: %v\n", err)
		return
	}

	fmt.Printf("   📊 Final Status: %s\n", status.Status)
	if status.Message != "" {
		fmt.Printf("   📝 Message: %s\n", status.Message)
	}

	// Show status meaning
	switch status.Status {
	case rimpay.PaymentStatusSuccess:
		fmt.Printf("   🎉 Payment completed successfully\n")
	case rimpay.PaymentStatusFailed:
		fmt.Printf("   💔 Payment failed permanently\n")
	case rimpay.PaymentStatusPending:
		fmt.Printf("   ⏳ Payment is still being processed\n")
	}
}
