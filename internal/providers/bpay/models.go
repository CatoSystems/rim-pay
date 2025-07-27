package bpay

import "github.com/CatoSystems/rim-pay/pkg/rimpay"

type AuthRequest struct {
	GrantType string `json:"grant_type"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	ClientID  string `json:"client_id"`
}

// AuthResponse represents B-PAY authentication response
type AuthResponse struct {
	AccessToken      string `json:"access_token"`
	ExpiresIn        string `json:"expires_in"`
	RefreshToken     string `json:"refresh_token"`
	RefreshExpiresIn string `json:"refresh_expires_in"`
}

// PaymentRequest represents B-PAY payment request
type PaymentRequest struct {
	ClientPhone string `json:"clientPhone"`
	Passcode    string `json:"passcode"`
	OperationID string `json:"operationId"`
	Amount      string `json:"amount"`
	Language    string `json:"language,omitempty"`
}

// PaymentResponse represents B-PAY payment response
type PaymentResponse struct {
	ErrorCode     string `json:"errorCode"`
	ErrorMessage  string `json:"errorMessage"`
	TransactionID string `json:"transactionId"`
}

// CheckTransactionRequest represents status check request
type CheckTransactionRequest struct {
	OperationID string `json:"operationID"`
}

// CheckTransactionResponse represents status check response
type CheckTransactionResponse struct {
	ErrorCode     string `json:"errorCode"`
	ErrorMessage  string `json:"errorMessage"`
	TransactionID string `json:"transactionId"`
	Status        string `json:"status"`
}

// convertErrorCodeToStatus converts B-PAY error code to payment status
func convertErrorCodeToStatus(errorCode string) rimpay.PaymentStatus {
	switch errorCode {
	case "0":
		return rimpay.PaymentStatusSuccess
	case "2":
		return rimpay.PaymentStatusFailed // Invalid token
	case "4":
		return rimpay.PaymentStatusFailed // Operation ID required
	case "1":
		return rimpay.PaymentStatusFailed // Other error
	default:
		return rimpay.PaymentStatusPending
	}
}

// convertTransactionStatus converts B-PAY status to payment status
func convertTransactionStatus(status string) rimpay.PaymentStatus {
	switch status {
	case "TS":
		return rimpay.PaymentStatusSuccess // Transaction success
	case "TF":
		return rimpay.PaymentStatusFailed // Transaction failed
	case "TA":
		return rimpay.PaymentStatusPending // Transaction pending
	default:
		return rimpay.PaymentStatusPending
	}
}
