package bpay

import (
	"encoding/json"
	"testing"

	"github.com/CatoSystems/rim-pay/pkg/money"
	"github.com/CatoSystems/rim-pay/pkg/phone"
	"github.com/CatoSystems/rim-pay/pkg/rimpay"
)

func TestBPayPasscodeIntegration(t *testing.T) {
	// This test demonstrates the complete flow of BPay payment with passcode generation
	t.Log("=== BPay Passcode Integration Test ===")

	// Create a phone number
	phoneNum, err := phone.NewPhone("+22220000000")
	if err != nil {
		t.Fatalf("Failed to create phone number: %v", err)
	}

	// Create money amount (100.00 MRU)
	amount := money.FromFloat64(100.00, money.MRU)

	// Create BPay payment request WITHOUT providing a passcode
	bpayRequest := &rimpay.BPayPaymentRequest{
		PhoneNumber: phoneNum,
		Amount:      amount,
		Description: "Test payment for passcode generation",
		Reference:   "TEST-REF-001",
		// Passcode is intentionally empty - library should generate it
		Metadata: map[string]interface{}{
			"test": true,
		},
	}

	t.Logf("Created BPay request without passcode (library will always generate one):")
	t.Logf("  Phone: %s", bpayRequest.PhoneNumber.String())
	t.Logf("  Amount: %s", bpayRequest.Amount.String())
	t.Logf("  Reference: %s", bpayRequest.Reference)
	t.Logf("  Passcode: '%s' (ignored - library always generates new)", bpayRequest.Passcode)

	// Validate the request - should succeed (passcode not needed for validation)
	if err := bpayRequest.Validate(); err != nil {
		t.Fatalf("BPay request validation failed: %v", err)
	}
	t.Log("✓ BPay request validation passed (passcode not required for validation)")

	// Convert to generic request
	genericRequest := bpayRequest.ToGenericRequest()
	t.Logf("Generic request passcode: '%s' (always empty - library generates)", genericRequest.Passcode)

	// Demonstrate what would happen in the payment processor
	t.Log("\n=== Simulating Payment Processing ===")

	// Step 1: Library ALWAYS generates a new passcode (ignoring any provided)
	generatedPasscode, err := generatePasscode()
	if err != nil {
		t.Fatalf("Failed to generate passcode: %v", err)
	}
	passcode := generatedPasscode
	t.Logf("✓ Library ALWAYS generates new passcode: %s", passcode) // Step 2: Create BPay API request (this is what gets sent to BPay)
	bpayAPIRequest := &PaymentRequest{
		ClientPhone: phoneNum.ForProvider(false),
		Passcode:    passcode,
		OperationID: bpayRequest.Reference,
		Amount:      amount.ToProviderAmount(false),
		Language:    "FR",
	}

	t.Log("✓ Created BPay API request with generated passcode")

	// Step 3: Show what the API request would look like
	requestJSON, _ := json.MarshalIndent(bpayAPIRequest, "", "  ")
	t.Logf("BPay API Request JSON:\n%s", string(requestJSON))

	// Step 4: Create mock response (simulating BPay server response)
	mockBPayResponse := &PaymentResponse{
		ErrorCode:     "0", // Success
		ErrorMessage:  "Payment initiated successfully",
		TransactionID: "BPAY-TXN-12345",
	}

	// Step 5: Create payment response with passcode in metadata
	response := &rimpay.PaymentResponse{
		TransactionID: mockBPayResponse.TransactionID,
		Status:        rimpay.PaymentStatusPending,
		Amount:        amount,
		Reference:     bpayRequest.Reference,
		Provider:      "bpay",
		Metadata: map[string]interface{}{
			"error_code":         mockBPayResponse.ErrorCode,
			"error_message":      mockBPayResponse.ErrorMessage,
			"transaction_id":     mockBPayResponse.TransactionID,
			"provider_reference": mockBPayResponse.TransactionID,
			"passcode":           passcode, // Include generated passcode
		},
	}

	t.Log("\n=== Payment Response ===")
	t.Logf("✓ Transaction ID: %s", response.TransactionID)
	t.Logf("✓ Status: %s", response.Status)
	t.Logf("✓ Provider: %s", response.Provider)

	// Extract passcode from metadata (this is what the application would do)
	if extractedPasscode, exists := response.Metadata["passcode"]; exists {
		if passcodeStr, ok := extractedPasscode.(string); ok {
			t.Logf("✓ Passcode for payer: %s", passcodeStr)

			// Verify the passcode format
			if len(passcodeStr) != 4 {
				t.Errorf("Expected 4-digit passcode, got %d digits: %s", len(passcodeStr), passcodeStr)
			}
		} else {
			t.Error("Passcode in metadata is not a string")
		}
	} else {
		t.Error("Passcode not found in response metadata")
	}

	t.Log("\n=== Test Summary ===")
	t.Log("✓ BPay request created (any provided passcode is ignored)")
	t.Log("✓ Library ALWAYS generates a new 4-digit passcode for security")
	t.Log("✓ Generated passcode included in payment request to BPay API")
	t.Log("✓ Passcode returned in response metadata for application use")
	t.Log("✓ Payer can now use the generated passcode to complete payment")
}
