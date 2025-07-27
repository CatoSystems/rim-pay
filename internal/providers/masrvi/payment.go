package masrvi

import (
	"context"
	"github.com/CatoSystems/rim-pay/internal/providers/common"
	"github.com/CatoSystems/rim-pay/pkg/rimpay"
	"net/url"
	"time"
)

// PaymentProcessor handles MASRVI payment operations
type PaymentProcessor struct {
	config         rimpay.ProviderConfig
	httpClient     common.HTTPClient
	sessionManager *SessionManager
	logger         rimpay.Logger
	baseURL        string
}

// NewPaymentProcessor creates new payment processor
func NewPaymentProcessor(config rimpay.ProviderConfig, httpClient common.HTTPClient, sessionManager *SessionManager, logger rimpay.Logger) *PaymentProcessor {
	return &PaymentProcessor{
		config:         config,
		httpClient:     httpClient,
		sessionManager: sessionManager,
		logger:         logger,
		baseURL:        config.BaseURL,
	}
}

// ProcessPayment processes a payment request
func (pp *PaymentProcessor) ProcessPayment(ctx context.Context, request *rimpay.PaymentRequest) (*rimpay.PaymentResponse, error) {
	// Get session ID
	sessionID, err := pp.sessionManager.GetSessionID(ctx)
	if err != nil {
		return nil, rimpay.NewPaymentError(
			rimpay.ErrorCodeProviderError,
			"failed to get session ID",
			"masrvi",
			true,
		)
	}

	// Create form data
	formData := pp.createFormData(sessionID, request)

	// Create payment URL
	paymentURL := pp.baseURL + "/online/online.php"

	pp.logger.Info("MASRVI payment created",
		"reference", request.Reference,
		"session_id", sessionID,
		"amount", request.Amount.String(),
	)

	// Create response
	response := &rimpay.PaymentResponse{
		TransactionID: request.Reference, // Use reference as initial transaction ID
		Status:        rimpay.PaymentStatusPending,
		Amount:        request.Amount,
		Reference:     request.Reference,
		Provider:      "masrvi",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		PaymentURL:    paymentURL,
		Metadata: map[string]interface{}{
			"session_id":  sessionID,
			"form_data":   formData,
			"payment_url": paymentURL,
			"message":     "Payment initiated, redirect user to payment URL",
		},
	}

	return response, nil
}

// createFormData creates form data for MASRVI
func (pp *PaymentProcessor) createFormData(sessionID string, request *rimpay.PaymentRequest) url.Values {
	formData := url.Values{}
	formData.Set("sessionid", sessionID)
	formData.Set("merchantid", pp.config.Credentials["merchant_id"])
	formData.Set("amount", request.Amount.ToProviderAmount(true)) // MASRVI uses cents
	formData.Set("currency", request.Amount.GetCurrencyCode())
	formData.Set("purchaseref", request.Reference)
	formData.Set("description", request.Description)

	// Optional fields
	if request.PhoneNumber != nil {
		formData.Set("phonenumber", request.PhoneNumber.LocalFormat())
	}

	if request.SuccessURL != "" {
		formData.Set("accepturl", request.SuccessURL)
	}

	if request.FailureURL != "" {
		formData.Set("declineurl", request.FailureURL)
	}

	if request.CancelURL != "" {
		formData.Set("cancelurl", request.CancelURL)
	}

	// Brand name from config or request metadata
	if brandName, exists := pp.config.Options["brand_name"].(string); exists {
		formData.Set("brand", brandName)
	}

	return formData
}

// HandleNotification handles webhook notifications
func (pp *PaymentProcessor) HandleNotification(notification *NotificationData) (*rimpay.TransactionStatus, error) {
	if notification == nil {
		return nil, rimpay.NewValidationError("notification", "is required")
	}

	status := notification.ToPaymentStatus()
	message := "Payment notification received"

	if status == rimpay.PaymentStatusFailed && notification.Error != "" {
		message = notification.Error
	}

	pp.logger.Info("MASRVI notification processed",
		"reference", notification.PurchaseRef,
		"status", status,
		"payment_ref", notification.PaymentRef,
	)

	transactionStatus := &rimpay.TransactionStatus{
		TransactionID:     notification.PayID,
		Status:            status,
		Reference:         notification.PurchaseRef,
		ProviderReference: notification.PaymentRef,
		Message:           message,
		LastUpdated:       time.Now(),
		ProviderData: map[string]interface{}{
			"client_id":   notification.ClientID,
			"client_name": notification.ClientName,
			"mobile":      notification.Mobile,
			"payment_ref": notification.PaymentRef,
			"pay_id":      notification.PayID,
			"timestamp":   notification.Timestamp,
			"ip_address":  notification.IPAddress,
			"status":      notification.Status,
		},
	}

	// Add status event
	transactionStatus.AddEvent(status, message)

	return transactionStatus, nil
}
