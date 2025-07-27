package main

import (
	"fmt"
	"time"

	"github.com/shopspring/decimal"
	"github.com/CatoSystems/rim-pay/pkg/money"
	"github.com/CatoSystems/rim-pay/pkg/phone"
	"github.com/CatoSystems/rim-pay/pkg/rimpay"
)

func main() {
	// Run examples for provider-specific payment types
	if err := runProviderSpecificExamples(); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
}

func runProviderSpecificExamples() error {
	fmt.Println("=== Provider-Specific Payment Examples ===\n")

	// Create client with proper configuration
	client, err := createTestClient()
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	// Run B-PAY example
	if err := runBPayExample(client); err != nil {
		return err
	}

	// Run MASRVI example
	if err := runMasrviExample(client); err != nil {
		return err
	}

	return nil
}

func runBPayExample(client *rimpay.Client) error {
	fmt.Println("--- B-PAY Provider-Specific Payment ---")

	// Create B-PAY specific payment request
	phone, err := phone.NewPhone("+22233445566")
	if err != nil {
		return fmt.Errorf("invalid phone number: %w", err)
	}

	amount := money.New(decimal.NewFromInt(100000), "MRU") // 1000 MRU

	bpayRequest := &rimpay.BPayPaymentRequest{
		Amount:      amount,
		PhoneNumber: phone,
		Reference:   "BPAY-" + fmt.Sprintf("%d", time.Now().Unix()),
		Description: "Test B-PAY payment with provider-specific request",
		Passcode:    "1234",
	}

	// Validate the request
	if err := bpayRequest.Validate(); err != nil {
		return fmt.Errorf("B-PAY request validation failed: %w", err)
	}

	fmt.Printf("B-PAY Request created successfully:\n")
	fmt.Printf("  Amount: %s\n", amount.String())
	fmt.Printf("  Phone: %s\n", phone.String())
	fmt.Printf("  Reference: %s\n", bpayRequest.Reference)
	fmt.Printf("  Passcode: %s\n", bpayRequest.Passcode)

	// Note: In a real implementation, you would call:
	// response, err := client.ProcessBPayPayment(context.Background(), bpayRequest)
	fmt.Printf("✓ B-PAY request ready for processing\n\n")

	return nil
}

func runMasrviExample(client *rimpay.Client) error {
	fmt.Println("--- MASRVI Provider-Specific Payment ---")

	// Create MASRVI specific payment request
	phone, err := phone.NewPhone("+22233889900")
	if err != nil {
		return fmt.Errorf("invalid phone number: %w", err)
	}

	amount := money.New(decimal.NewFromInt(50000), "MRU") // 500 MRU

	masrviRequest := &rimpay.MasrviPaymentRequest{
		Amount:      amount,
		PhoneNumber: phone,
		Reference:   "MASRVI-" + fmt.Sprintf("%d", time.Now().Unix()),
		Description: "Test MASRVI payment with provider-specific request",
		CallbackURL: "https://webhook.example.com/masrvi",
		ReturnURL:   "https://shop.example.com/return",
	}

	// Validate the request
	if err := masrviRequest.Validate(); err != nil {
		return fmt.Errorf("MASRVI request validation failed: %w", err)
	}

	fmt.Printf("MASRVI Request created successfully:\n")
	fmt.Printf("  Amount: %s\n", amount.String())
	fmt.Printf("  Phone: %s\n", phone.String())
	fmt.Printf("  Reference: %s\n", masrviRequest.Reference)
	fmt.Printf("  Callback URL: %s\n", masrviRequest.CallbackURL)
	fmt.Printf("  Return URL: %s\n", masrviRequest.ReturnURL)

	// Note: In a real implementation, you would call:
	// response, err := client.ProcessMasrviPayment(context.Background(), masrviRequest)
	fmt.Printf("✓ MASRVI request ready for processing\n\n")

	return nil
}

func createTestClient() (*rimpay.Client, error) {
	// Create a minimal configuration
	config := rimpay.DefaultConfig()
	config.Environment = rimpay.EnvironmentSandbox
	config.DefaultProvider = "bpay"

	// Add provider configurations
	config.Providers[rimpay.ProviderBPay] = rimpay.ProviderConfig{
		Enabled: true,
		BaseURL: "https://api.bpay.mr",
		Credentials: map[string]string{
			"client_id":     "test_client_id",
			"client_secret": "test_client_secret",
		},
		Timeout: 30 * time.Second,
	}

	config.Providers[rimpay.ProviderMasrvi] = rimpay.ProviderConfig{
		Enabled: true,
		BaseURL: "https://pay.masrvi.mr",
		Credentials: map[string]string{
			"merchant_id": "test_merchant",
			"api_key":     "test_api_key",
		},
		Timeout: 30 * time.Second,
	}

	return rimpay.NewClient(config)
}

// Demonstration of converting provider-specific requests to generic format
func demonstrateConversion() {
	fmt.Println("--- Demonstrating Provider-Specific to Generic Conversion ---")

	// Create a B-PAY request
	phone, _ := phone.NewPhone("+22233445566")
	amount := money.New(decimal.NewFromInt(100000), "MRU")

	bpayRequest := &rimpay.BPayPaymentRequest{
		Amount:      amount,
		PhoneNumber: phone,
		Reference:   "DEMO-BPAY-123",
		Description: "Demo payment",
		Passcode:    "1234",
	}

	// Convert to generic request
	genericRequest := bpayRequest.ToGenericRequest()

	fmt.Printf("Original B-PAY Request:\n")
	fmt.Printf("  Passcode: %s\n", bpayRequest.Passcode)

	fmt.Printf("Converted Generic Request:\n")
	fmt.Printf("  Amount: %s\n", genericRequest.Amount.String())
	fmt.Printf("  Reference: %s\n", genericRequest.Reference)
	fmt.Printf("  Metadata contains provider data: %v\n", len(genericRequest.Metadata) > 0)

	fmt.Println("✓ Conversion successful\n")
}
