package click

import (
	"context"
	"testing"

	"github.com/CatoSystems/rim-pay/pkg/money"
	"github.com/CatoSystems/rim-pay/pkg/phone"
	"github.com/CatoSystems/rim-pay/pkg/rimpay"
)

func TestProviderImplementsInterface(t *testing.T) {
	var _ rimpay.ClickProvider = (*Provider)(nil)
}

func TestProcessClickPaymentValidates(t *testing.T) {
	p, err := NewClickProvider(testConfig(), nopLogger{})
	if err != nil {
		t.Fatalf("NewClickProvider: %v", err)
	}
	// Missing reference -> validation error, no network call.
	phoneNum, _ := phone.NewPhone("+22220000000")
	_, err = p.ProcessClickPayment(context.Background(), &rimpay.ClickPaymentRequest{
		PhoneNumber: phoneNum,
		Amount:      money.FromFloat64(15.00, money.MRU),
		Reference:   "",
	})
	if err == nil {
		t.Error("expected validation error for empty reference")
	}
}
