package main

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/CatoSystems/rim-pay/pkg/money"
	"github.com/CatoSystems/rim-pay/pkg/phone"
	_ "github.com/CatoSystems/rim-pay/pkg/providers" // register all providers
	"github.com/CatoSystems/rim-pay/pkg/rimpay"
	"github.com/shopspring/decimal"
)

func main() {
	fmt.Println("🏦 RimPay Library - CLICK (BNM) Provider Example")
	fmt.Println("================================================")

	// CLICK is BNM's TagPay online.php hosted-page integration.
	clickConfig := rimpay.ProviderConfig{
		Enabled: true,
		BaseURL: "https://tagpay.example", // TagPay online base URL (from BNM)
		Timeout: 30 * time.Second,
		Credentials: map[string]string{
			"merchant_id": "0896353536734538", // 16-digit merchant ID from BNM
		},
		Options: map[string]interface{}{
			"brand_name": "My Shop",
		},
	}

	// The provider must be present in the config map (so validation passes) and
	// then registered on the client so an instance is created.
	config := rimpay.DefaultConfig()
	config.DefaultProvider = "click"
	config.Providers["click"] = clickConfig

	client, err := rimpay.NewClient(config)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	if err := client.AddClickProvider(clickConfig); err != nil {
		log.Fatalf("Failed to add CLICK provider: %v", err)
	}

	ctx := context.Background()
	phoneNum, _ := phone.NewPhone("+22220000000")

	// Step 1: create the payment. The library requests a TagPay session and
	// returns the order form the customer's browser must POST to PaymentURL.
	resp, err := client.ProcessClickPayment(ctx, &rimpay.ClickPaymentRequest{
		PhoneNumber: phoneNum,
		Amount:      money.New(decimal.NewFromFloat(15.00), money.MRU),
		Reference:   fmt.Sprintf("PURCHASE-%d", time.Now().UnixNano()),
		Description: "Online purchase",
		SuccessURL:  "https://shop.example/ok",
		FailureURL:  "https://shop.example/fail",
		CancelURL:   "https://shop.example/cancel",
	})
	if err != nil {
		log.Fatalf("ProcessClickPayment failed: %v", err)
	}

	fmt.Printf("➡️  Redirect the customer's browser to: %s\n", resp.PaymentURL)
	if form, ok := resp.Metadata["form_data"].(url.Values); ok {
		fmt.Println("   POST these hidden form fields:")
		for key := range form {
			fmt.Printf("     %s = %s\n", key, form.Get(key))
		}
	}

	// Step 2: when TagPay calls your Notification URL (server-to-server GET),
	// parse the query parameters and hand them to HandleClickNotification.
	// Example of what your notification HTTP handler would do:
	//
	//   q := r.URL.Query()
	//   status, _ := client.HandleClickNotification(&rimpay.ClickNotificationData{
	//       Status:      q.Get("status"),
	//       PurchaseRef: q.Get("purchaseref"),
	//       Amount:      q.Get("amount"),
	//       Currency:    q.Get("currency"),
	//       ClientID:    q.Get("clientid"),
	//       ClientName:  q.Get("cname"),
	//       Mobile:      q.Get("mobile"),
	//       PaymentRef:  q.Get("paymentref"),
	//       PayID:       q.Get("payid"),
	//       Timestamp:   q.Get("timestamp"),
	//       IPAddress:   q.Get("ipaddr"),
	//       Error:       q.Get("error"),
	//       Reason:      q.Get("reason"),
	//   })
	//
	// IMPORTANT: verify the amount/reference against your own order and restrict
	// the notification endpoint to TagPay's IP address.
	fmt.Println("\n✅ Redirect the customer, then handle the TagPay notification callback.")
}
