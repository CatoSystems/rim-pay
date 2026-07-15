package click

import (
	"context"
	"net/url"
	"strings"
	"time"

	"github.com/CatoSystems/rim-pay/internal/providers/common"
	"github.com/CatoSystems/rim-pay/pkg/rimpay"
)

// PaymentProcessor handles CLICK payment operations.
type PaymentProcessor struct {
	config         rimpay.ProviderConfig
	httpClient     common.HTTPClient
	sessionManager *SessionManager
	logger         rimpay.Logger
	baseURL        string
}

// NewPaymentProcessor creates a new CLICK payment processor.
func NewPaymentProcessor(config rimpay.ProviderConfig, httpClient common.HTTPClient, sessionManager *SessionManager, logger rimpay.Logger) *PaymentProcessor {
	return &PaymentProcessor{
		config:         config,
		httpClient:     httpClient,
		sessionManager: sessionManager,
		logger:         logger,
		baseURL:        strings.TrimRight(config.BaseURL, "/"),
	}
}

// ProcessPayment creates a CLICK session and builds the order form.
func (pp *PaymentProcessor) ProcessPayment(ctx context.Context, request *rimpay.PaymentRequest) (*rimpay.PaymentResponse, error) {
	sessionID, err := pp.sessionManager.GetSessionID(ctx)
	if err != nil {
		return nil, rimpay.NewPaymentError(rimpay.ErrorCodeProviderError, "failed to get session ID", "click", true)
	}

	formData := pp.createFormData(sessionID, request)
	paymentURL := pp.baseURL + "/online/online.php"

	pp.logger.Info("CLICK payment created",
		"reference", request.Reference,
		"amount", request.Amount.String(),
	)

	return &rimpay.PaymentResponse{
		TransactionID: request.Reference,
		Status:        rimpay.PaymentStatusPending,
		Amount:        request.Amount,
		Reference:     request.Reference,
		Provider:      "click",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		PaymentURL:    paymentURL,
		Metadata: map[string]interface{}{
			"session_id":  sessionID,
			"form_data":   formData,
			"payment_url": paymentURL,
			"message":     "Payment initiated, redirect user to payment URL",
		},
	}, nil
}

// createFormData builds the lowercase TagPay order form.
func (pp *PaymentProcessor) createFormData(sessionID string, request *rimpay.PaymentRequest) url.Values {
	form := url.Values{}
	form.Set("sessionid", sessionID)
	form.Set("merchantid", pp.config.Credentials["merchant_id"])
	form.Set("amount", request.Amount.ToProviderAmount(true)) // cents
	form.Set("currency", request.Amount.GetCurrencyCode())    // ISO 4217 numeric
	form.Set("purchaseref", request.Reference)
	if request.Description != "" {
		form.Set("description", request.Description)
	}
	if request.PhoneNumber != nil {
		form.Set("phonenumber", request.PhoneNumber.LocalFormat())
	}
	if request.SuccessURL != "" {
		form.Set("accepturl", request.SuccessURL)
	}
	if request.FailureURL != "" {
		form.Set("declineurl", request.FailureURL)
	}
	if request.CancelURL != "" {
		form.Set("cancelurl", request.CancelURL)
	}
	if brand, ok := request.Metadata["brand"].(string); ok && brand != "" {
		form.Set("brand", brand)
	} else if brand, ok := pp.config.Options["brand_name"].(string); ok && brand != "" {
		form.Set("brand", brand)
	}
	if request.Language != "" {
		form.Set("language", strings.ToLower(string(request.Language)))
	}
	return form
}

// HandleNotification converts a TagPay notification into a TransactionStatus.
func (pp *PaymentProcessor) HandleNotification(notification *NotificationData) (*rimpay.TransactionStatus, error) {
	if notification == nil {
		return nil, rimpay.NewValidationError("notification", "is required")
	}

	status := notification.ToPaymentStatus()
	message := "Payment notification received"
	if status != rimpay.PaymentStatusSuccess {
		if notification.Reason != "" {
			message = notification.Reason
		} else if notification.Error != "" {
			message = notification.Error
		}
	}

	pp.logger.Info("CLICK notification processed",
		"reference", notification.PurchaseRef,
		"status", status,
		"payment_ref", notification.PaymentRef,
	)

	ts := &rimpay.TransactionStatus{
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
			"error":       notification.Error,
			"reason":      notification.Reason,
		},
	}
	ts.AddEvent(status, message)
	return ts, nil
}
