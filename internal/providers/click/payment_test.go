package click

import (
	"context"
	"net/url"
	"testing"

	"github.com/CatoSystems/rim-pay/pkg/money"
	"github.com/CatoSystems/rim-pay/pkg/phone"
	"github.com/CatoSystems/rim-pay/pkg/rimpay"
)

func TestProcessPaymentBuildsOrderForm(t *testing.T) {
	stub := &stubHTTP{body: "OK:SESSION123"}
	cfg := testConfig()
	sm := NewSessionManager(cfg, stub, nopLogger{})
	pp := NewPaymentProcessor(cfg, stub, sm, nopLogger{})

	p, _ := phone.NewPhone("+22220000000")
	req := &rimpay.PaymentRequest{
		Amount:      money.FromFloat64(15.00, money.MRU),
		PhoneNumber: p,
		Reference:   "PURCHASE0987",
		Description: "Online purchase",
		SuccessURL:  "https://shop.test/ok",
	}

	resp, err := pp.ProcessPayment(context.Background(), req)
	if err != nil {
		t.Fatalf("ProcessPayment: %v", err)
	}
	if resp.Status != rimpay.PaymentStatusPending {
		t.Errorf("status = %v, want pending", resp.Status)
	}
	if resp.PaymentURL != "https://tagpay.test/online/online.php" {
		t.Errorf("payment url = %q", resp.PaymentURL)
	}
	form, ok := resp.Metadata["form_data"].(url.Values)
	if !ok {
		t.Fatalf("form_data missing or wrong type")
	}
	if form.Get("sessionid") != "SESSION123" {
		t.Errorf("sessionid = %q", form.Get("sessionid"))
	}
	if form.Get("merchantid") != "0896353536734538" {
		t.Errorf("merchantid = %q", form.Get("merchantid"))
	}
	if form.Get("amount") != "1500" { // 15.00 MRU in cents
		t.Errorf("amount = %q, want 1500", form.Get("amount"))
	}
	if form.Get("currency") != "929" {
		t.Errorf("currency = %q, want 929", form.Get("currency"))
	}
	if form.Get("purchaseref") != "PURCHASE0987" {
		t.Errorf("purchaseref = %q", form.Get("purchaseref"))
	}
	if form.Get("accepturl") != "https://shop.test/ok" {
		t.Errorf("accepturl = %q", form.Get("accepturl"))
	}
}

func TestHandleNotificationSuccess(t *testing.T) {
	pp := NewPaymentProcessor(testConfig(), &stubHTTP{}, nil, nopLogger{})
	ts, err := pp.HandleNotification(&NotificationData{
		Status:      "OK",
		PurchaseRef: "PURCHASE0987",
		PaymentRef:  "O63186141-164065699",
		PayID:       "56092",
	})
	if err != nil {
		t.Fatalf("HandleNotification: %v", err)
	}
	if ts.Status != rimpay.PaymentStatusSuccess {
		t.Errorf("status = %v", ts.Status)
	}
	if ts.TransactionID != "56092" || ts.ProviderReference != "O63186141-164065699" {
		t.Errorf("ids wrong: %+v", ts)
	}
}
