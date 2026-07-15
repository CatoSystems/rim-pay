package click

import "github.com/CatoSystems/rim-pay/pkg/rimpay"

// NotificationData is a TagPay server-to-server notification for CLICK.
type NotificationData struct {
	Status      string // OK / NOK
	ClientID    string
	ClientName  string
	Mobile      string
	PurchaseRef string
	PaymentRef  string
	PayID       string
	Timestamp   string
	IPAddress   string
	Error       string
	Reason      string
	Amount      string
	Currency    string
}

// ToPaymentStatus maps the TagPay status + error code to a rimpay status.
func (nd *NotificationData) ToPaymentStatus() rimpay.PaymentStatus {
	switch nd.Status {
	case "OK":
		return rimpay.PaymentStatusSuccess
	case "NOK":
		switch nd.Error {
		case "CANCEL":
			return rimpay.PaymentStatusCancelled
		case "EXPIRED_SESSION":
			return rimpay.PaymentStatusExpired
		case "AUTHENTICATION", "PAYMENT_FAILED":
			return rimpay.PaymentStatusFailed
		default:
			return rimpay.PaymentStatusFailed
		}
	default:
		return rimpay.PaymentStatusPending
	}
}
