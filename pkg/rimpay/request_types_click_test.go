package rimpay

import (
	"testing"

	"github.com/CatoSystems/rim-pay/pkg/money"
	"github.com/CatoSystems/rim-pay/pkg/phone"
)

func newValidClickRequest() *ClickPaymentRequest {
	p, _ := phone.NewPhone("+22220000000")
	return &ClickPaymentRequest{
		PhoneNumber: p,
		Amount:      money.FromFloat64(15.00, money.MRU),
		Reference:   "PURCHASE0987",
		Description: "Online purchase",
		SuccessURL:  "https://shop.test/ok",
		FailureURL:  "https://shop.test/fail",
		CancelURL:   "https://shop.test/cancel",
	}
}

func TestClickRequestValid(t *testing.T) {
	if err := newValidClickRequest().Validate(); err != nil {
		t.Fatalf("expected valid, got %v", err)
	}
}

func TestClickRequestRequiresAmountAndReference(t *testing.T) {
	req := newValidClickRequest()
	req.Reference = ""
	if err := req.Validate(); err == nil {
		t.Error("expected error for empty reference")
	}
}

func TestClickRequestToGeneric(t *testing.T) {
	req := newValidClickRequest()
	g := req.ToGenericRequest()
	if g.Reference != "PURCHASE0987" || g.SuccessURL != "https://shop.test/ok" {
		t.Errorf("generic mapping wrong: %+v", g)
	}
}
