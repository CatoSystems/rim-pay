package click

import (
	"testing"

	"github.com/CatoSystems/rim-pay/pkg/rimpay"
)

func TestNotificationToPaymentStatus(t *testing.T) {
	cases := []struct {
		status  string
		errCode string
		want    rimpay.PaymentStatus
	}{
		{"OK", "", rimpay.PaymentStatusSuccess},
		{"NOK", "CANCEL", rimpay.PaymentStatusCancelled},
		{"NOK", "AUTHENTICATION", rimpay.PaymentStatusFailed},
		{"NOK", "PAYMENT_FAILED", rimpay.PaymentStatusFailed},
		{"NOK", "EXPIRED_SESSION", rimpay.PaymentStatusExpired},
		{"NOK", "", rimpay.PaymentStatusFailed},
		{"weird", "", rimpay.PaymentStatusPending},
	}
	for _, c := range cases {
		nd := &NotificationData{Status: c.status, Error: c.errCode}
		if got := nd.ToPaymentStatus(); got != c.want {
			t.Errorf("status=%q err=%q: got %v want %v", c.status, c.errCode, got, c.want)
		}
	}
}
