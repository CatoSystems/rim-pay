package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/CatoSystems/rim-pay/pkg/money"
	"github.com/CatoSystems/rim-pay/pkg/phone"
	"github.com/CatoSystems/rim-pay/pkg/rimpay"
	"github.com/shopspring/decimal"
)

func main() {
	fmt.Println("🏦 RimPay Library - MASRVI Provider Example")
	fmt.Println("==========================================")

	config := &rimpay.Config{
		Environment:     rimpay.EnvironmentSandbox,
		DefaultProvider: "masrvi",
		Providers: map[string]rimpay.ProviderConfig{
			"masrvi": {
				Enabled: true,
				BaseURL: "https://masrviapp.mr/online",
				Timeout: 30 * time.Second,
				Credentials: map[string]string{
					"merchant_id": "TEST_MERCHANT_123", // Your MASRVI merchant ID
				},
			},
		},
	}

	client, err := rimpay.NewClient(config)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	fmt.Println("🔐 MASRVI Session Management")
	fmt.Println("The library handles session management automatically:")
	fmt.Println("1. Requests session ID using merchant ID")
	fmt.Println("2. Session is valid for 5 minutes")
	fmt.Println("3. Automatically creates new sessions when needed\n")

	ctx := context.Background()

	// Example 1: Basic e-commerce payment
	fmt.Println("🛒 Example 1: E-commerce Purchase")
	ecommercePayment := createMasrviPayment(
		"33445566",
		199.99,     // Amount
		"Online store purchase - Premium subscription",
		"ORDER-PREMIUM-001",
	)
	processMasrviPayment(client, ctx, ecommercePayment)

	// Example 2: Service payment with phone number
	fmt.Println("\n💳 Example 2: Service Payment")
	servicePayment := createMasrviPayment(
		"37889900",
		45.00,      // Amount
		"Internet service payment",
		"INTERNET-BILL-789",
	)
	processMasrviPayment(client, ctx, servicePayment)

	// Example 3: Mobile top-up
	fmt.Println("\n📱 Example 3: Mobile Top-up")
	topupPayment := createMasrviPayment(
		"48990011",
		25.50,      // Amount
		"Mobile credit top-up",
		"TOPUP-"+fmt.Sprintf("%d", time.Now().Unix()),
	)
	processMasrviPayment(client, ctx, topupPayment)

	// Start webhook server example
	fmt.Println("\n🔔 Setting up webhook server for notifications...")
	go startWebhookServer()

	fmt.Println("\n💡 MASRVI Features Demonstrated:")
	fmt.Println("✅ Session-based authentication")
	fmt.Println("✅ Payment form generation")
	fmt.Println("✅ Multiple redirect URLs (success/failure/cancel)")
	fmt.Println("✅ Webhook notification handling")
	fmt.Println("✅ Support for all Mauritanian operators")
	fmt.Println("✅ E-commerce integration ready")

	// Keep the webhook server running
	fmt.Println("\n🌐 Webhook server running on http://localhost:8080")
	fmt.Println("   POST /webhook - for payment notifications")
	fmt.Println("   Press Ctrl+C to stop")
	select {} // Keep running
}

func createMasrviPayment(phoneNumber string, amount float64, description, reference string) *rimpay.MasrviPaymentRequest {
	phone, err := phone.NewPhone(phoneNumber)
	if err != nil {
		log.Fatalf("Invalid phone number: %v", err)
	}

	money := money.New(decimal.NewFromFloat(amount), money.MRU)

	return &rimpay.MasrviPaymentRequest{
		Amount:      money,
		PhoneNumber: phone,
		Reference:   reference,
		Description: description,
		// MASRVI specific URLs
		ReturnURL:   "https://yourapp.com/payment/return",
		CallbackURL: "https://yourapp.com/webhook", // For notifications
	}
}

func processMasrviPayment(client *rimpay.Client, ctx context.Context, request *rimpay.MasrviPaymentRequest) {
	fmt.Printf("   Processing: %s → %s\n", 
		request.PhoneNumber.ForProvider(true), 
		request.Amount.String())
	fmt.Printf("   Reference: %s\n", request.Reference)
	fmt.Printf("   Description: %s\n", request.Description)

	// Process payment - this creates the payment form
	response, err := client.ProcessMasrviPayment(ctx, request)
	if err != nil {
		fmt.Printf("   ❌ Payment failed: %v\n", err)
		
		if paymentErr, ok := err.(*rimpay.PaymentError); ok {
			switch paymentErr.Code {
			case rimpay.ErrorCodeProviderError:
				fmt.Printf("   💡 MASRVI service may be temporarily unavailable\n")
			case rimpay.ErrorCodeNetworkError:
				fmt.Printf("   💡 Network issue - payment was retried automatically\n")
			case rimpay.ErrorCodeInvalidRequest:
				fmt.Printf("   💡 Check merchant ID and parameters\n")
			}
		}
		return
	}

	fmt.Printf("   ✅ Payment form created successfully!\n")
	fmt.Printf("   Transaction ID: %s\n", response.TransactionID)
	fmt.Printf("   Status: %s (pending customer action)\n", response.Status)
	
	if paymentURL, exists := response.Metadata["payment_url"]; exists {
		fmt.Printf("   🌐 Payment URL: %s\n", paymentURL)
		fmt.Printf("   💡 Customer should be redirected to this URL to complete payment\n")
	}

	fmt.Printf("   ⏳ Waiting for customer to complete payment...\n")
	fmt.Printf("   📱 Customer will receive SMS with payment instructions\n")
	fmt.Printf("   🔔 Webhook notifications will be sent to your callback URL\n")
}

// Webhook server to handle MASRVI notifications
func startWebhookServer() {
	http.HandleFunc("/webhook", handleWebhook)
	
	fmt.Println("🚀 Starting webhook server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Printf("Webhook server error: %v", err)
	}
}

func handleWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse MASRVI notification parameters
	status := r.URL.Query().Get("status")
	clientID := r.URL.Query().Get("clientid")
	customerName := r.URL.Query().Get("cname")
	mobile := r.URL.Query().Get("mobile")
	purchaseRef := r.URL.Query().Get("purchaseref")
	paymentRef := r.URL.Query().Get("paymentref")
	payID := r.URL.Query().Get("payid")
	timestamp := r.URL.Query().Get("timestamp")
	ipAddr := r.URL.Query().Get("ipaddr")
	errorMsg := r.URL.Query().Get("error")

	fmt.Printf("\n🔔 Webhook Notification Received:\n")
	fmt.Printf("   Status: %s\n", status)
	fmt.Printf("   Client ID: %s\n", clientID)
	fmt.Printf("   Customer: %s\n", customerName)
	fmt.Printf("   Mobile: %s\n", mobile)
	fmt.Printf("   Purchase Ref: %s\n", purchaseRef)
	fmt.Printf("   Payment Ref: %s\n", paymentRef)
	fmt.Printf("   Pay ID: %s\n", payID)
	fmt.Printf("   Timestamp: %s\n", timestamp)
	fmt.Printf("   IP Address: %s\n", ipAddr)
	
	if errorMsg != "" {
		fmt.Printf("   Error: %s\n", errorMsg)
	}

	// Process the notification
	switch status {
	case "Ok":
		fmt.Printf("   ✅ Payment successful!\n")
		// Update your database, send confirmation email, etc.
		handleSuccessfulPayment(purchaseRef, paymentRef, mobile)
		
	case "NOK":
		fmt.Printf("   ❌ Payment failed!\n")
		// Handle failed payment
		handleFailedPayment(purchaseRef, errorMsg)
		
	default:
		fmt.Printf("   ❓ Unknown status: %s\n", status)
	}

	// Respond to MASRVI
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func handleSuccessfulPayment(purchaseRef, paymentRef, mobile string) {
	fmt.Printf("   💰 Processing successful payment:\n")
	fmt.Printf("   📋 Order: %s\n", purchaseRef)
	fmt.Printf("   💳 Payment: %s\n", paymentRef)
	fmt.Printf("   📱 Customer: %s\n", mobile)
	
	// business logic here:
	// - Update order status in database
	// - Send confirmation email/SMS
	// - Trigger fulfillment process
	// - Update analytics
}

func handleFailedPayment(purchaseRef, errorMsg string) {
	fmt.Printf("   💔 Processing failed payment:\n")
	fmt.Printf("   📋 Order: %s\n", purchaseRef)
	fmt.Printf("   ❌ Error: %s\n", errorMsg)
	
	// business logic here:
	// - Update order status to failed
	// - Send failure notification
	// - Offer alternative payment methods
	// - Log for analysis
}