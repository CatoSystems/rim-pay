package masrvi

import "github.com/CatoSystems/rim-pay/pkg/rimpay"

type SessionResponse struct {
	SessionID string `json:"session_id"`
}

// FormData represents MASRVI form data
type FormData struct {
	SessionID   string `json:"sessionid"`
	MerchantID  string `json:"merchantid"`
	Amount      string `json:"amount"`
	Currency    string `json:"currency"`
	PurchaseRef string `json:"purchaseref"`
	Description string `json:"description"`
	PhoneNumber string `json:"phonenumber,omitempty"`
	Brand       string `json:"brand,omitempty"`
	AcceptURL   string `json:"accepturl,omitempty"`
	DeclineURL  string `json:"declineurl,omitempty"`
	CancelURL   string `json:"cancelurl,omitempty"`
	Text        string `json:"text,omitempty"`
}

// NotificationData represents MASRVI webhook notification
type NotificationData struct {
	Status      string `json:"status"` // OK/NOK
	ClientID    string `json:"clientid"`
	ClientName  string `json:"cname"`
	Mobile      string `json:"mobile"`
	PurchaseRef string `json:"purchaseref"`
	PaymentRef  string `json:"paymentref"`
	PayID       string `json:"payid"`
	Timestamp   string `json:"timestamp"`
	IPAddress   string `json:"ipaddr"`
	Error       string `json:"error,omitempty"`
}

// ToPaymentStatus converts notification status to payment status
func (nd *NotificationData) ToPaymentStatus() rimpay.PaymentStatus {
	switch nd.Status {
	case "Ok":
		return rimpay.PaymentStatusSuccess
	case "NOK":
		return rimpay.PaymentStatusFailed
	default:
		return rimpay.PaymentStatusPending
	}
}
