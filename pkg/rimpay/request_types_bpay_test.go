package rimpay

import (
	"testing"

	"github.com/CatoSystems/rim-pay/pkg/money"
	"github.com/CatoSystems/rim-pay/pkg/phone"
)

func newValidBPayRequest() *BPayPaymentRequest {
	p, _ := phone.NewPhone("+22220000000")
	return &BPayPaymentRequest{
		PhoneNumber: p,
		Amount:      money.FromFloat64(50.00, money.MRU),
		Description: "Test",
		Reference:   "REF-1",
		Passcode:    "1234",
	}
}

func TestBPayRequestForwardsPasscode(t *testing.T) {
	req := newValidBPayRequest()
	if err := req.Validate(); err != nil {
		t.Fatalf("expected valid request, got %v", err)
	}
	generic := req.ToGenericRequest()
	if generic.Passcode != "1234" {
		t.Errorf("passcode not forwarded: got %q want %q", generic.Passcode, "1234")
	}
}

func TestBPayRequestRejectsEmptyPasscode(t *testing.T) {
	req := newValidBPayRequest()
	req.Passcode = ""
	if err := req.Validate(); err == nil {
		t.Error("expected validation error for empty passcode, got nil")
	}
}

func TestBPayRequestRejectsShortPasscode(t *testing.T) {
	req := newValidBPayRequest()
	req.Passcode = "12"
	if err := req.Validate(); err == nil {
		t.Error("expected validation error for short passcode, got nil")
	}
}

func TestBPayRequestRejectsNonFourDigitPasscode(t *testing.T) {
	for _, bad := range []string{"12345", "abcd", "12a4", " 123"} {
		req := newValidBPayRequest()
		req.Passcode = bad
		if err := req.Validate(); err == nil {
			t.Errorf("expected validation error for passcode %q, got nil", bad)
		}
	}
}
