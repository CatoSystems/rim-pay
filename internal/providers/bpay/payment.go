package bpay

import (
	"context"
	"encoding/json"
	"github.com/CatoSystems/rim-pay/internal/providers/common"
	"github.com/CatoSystems/rim-pay/pkg/rimpay"
	"time"
)

// PaymentProcessor handles B-PAY payment operations
type PaymentProcessor struct {
	config      rimpay.ProviderConfig
	httpClient  common.HTTPClient
	authManager *AuthManager
	logger      rimpay.Logger
	baseURL     string
}

// NewPaymentProcessor creates new payment processor
func NewPaymentProcessor(config rimpay.ProviderConfig, httpClient common.HTTPClient, authManager *AuthManager, logger rimpay.Logger) *PaymentProcessor {
	return &PaymentProcessor{
		config:      config,
		httpClient:  httpClient,
		authManager: authManager,
		logger:      logger,
		baseURL:     config.BaseURL,
	}
}

// ProcessPayment processes a payment request
func (pp *PaymentProcessor) ProcessPayment(ctx context.Context, request *rimpay.PaymentRequest) (*rimpay.PaymentResponse, error) {
	// Get access token
	token, err := pp.authManager.GetAccessToken(ctx)
	if err != nil {
		return nil, rimpay.NewPaymentError(
			rimpay.ErrorCodeAuthenticationFailed,
			"failed to get access token",
			"bpay",
			true,
		)
	}

	// Create B-PAY specific request
	bpayReq := &PaymentRequest{
		ClientPhone: request.PhoneNumber.ForProvider(false),
		Passcode:    request.Passcode,
		OperationID: request.Reference,
		Amount:      request.Amount.ToProviderAmount(false),
		Language:    convertLanguage(request.GetLanguage()),
	}

	// Marshal request
	payload, err := json.Marshal(bpayReq)
	if err != nil {
		return nil, rimpay.NewPaymentError(
			rimpay.ErrorCodeInvalidRequest,
			"failed to marshal payment request",
			"bpay",
			false,
		)
	}

	// Create HTTP request
	httpReq := &common.HTTPRequest{
		Method: "POST",
		URL:    pp.baseURL + "/payment",
		Headers: map[string]string{
			"Content-Type":  "application/json",
			"Authorization": "Bearer " + token,
		},
		Body:    payload,
		Timeout: pp.config.Timeout,
	}

	pp.logger.Info("Making B-PAY payment request",
		"operation_id", bpayReq.OperationID,
		"amount", bpayReq.Amount,
	)

	// Execute request
	resp, err := pp.httpClient.Do(httpReq)
	if err != nil {
		return nil, rimpay.NewPaymentError(
			rimpay.ErrorCodeNetworkError,
			"payment request failed",
			"bpay",
			true,
		)
	}

	// Parse response
	var bpayResp PaymentResponse
	if err := json.Unmarshal(resp.Body, &bpayResp); err != nil {
		return nil, rimpay.NewPaymentError(
			rimpay.ErrorCodeProviderError,
			"failed to decode payment response",
			"bpay",
			false,
		)
	}

	// Convert to standard response
	status := convertErrorCodeToStatus(bpayResp.ErrorCode)

	response := &rimpay.PaymentResponse{
		TransactionID: bpayResp.TransactionID,
		Status:        status,
		Amount:        request.Amount,
		Reference:     request.Reference,
		Provider:      "bpay",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		Metadata: map[string]interface{}{
			"error_code":        bpayResp.ErrorCode,
			"error_message":     bpayResp.ErrorMessage,
			"transaction_id":    bpayResp.TransactionID,
			"provider_reference": bpayResp.TransactionID,
		},
	}

	pp.logger.Info("B-PAY payment response received",
		"transaction_id", response.TransactionID,
		"status", response.Status,
	)

	return response, nil
}

// CheckPaymentStatus checks payment status
func (pp *PaymentProcessor) CheckPaymentStatus(ctx context.Context, transactionID string) (*rimpay.TransactionStatus, error) {
	// Get access token
	token, err := pp.authManager.GetAccessToken(ctx)
	if err != nil {
		return nil, rimpay.NewPaymentError(
			rimpay.ErrorCodeAuthenticationFailed,
			"failed to get access token",
			"bpay",
			true,
		)
	}

	// Create check request
	checkReq := &CheckTransactionRequest{
		OperationID: transactionID,
	}

	payload, err := json.Marshal(checkReq)
	if err != nil {
		return nil, rimpay.NewPaymentError(
			rimpay.ErrorCodeInvalidRequest,
			"failed to marshal check request",
			"bpay",
			false,
		)
	}

	// Create HTTP request
	httpReq := &common.HTTPRequest{
		Method: "POST",
		URL:    pp.baseURL + "/checkTransaction",
		Headers: map[string]string{
			"Content-Type":  "application/json",
			"Authorization": "Bearer " + token,
		},
		Body:    payload,
		Timeout: pp.config.Timeout,
	}

	// Execute request
	resp, err := pp.httpClient.Do(httpReq)
	if err != nil {
		return nil, rimpay.NewPaymentError(
			rimpay.ErrorCodeNetworkError,
			"status check failed",
			"bpay",
			true,
		)
	}

	// Parse response
	var checkResp CheckTransactionResponse
	if err := json.Unmarshal(resp.Body, &checkResp); err != nil {
		return nil, rimpay.NewPaymentError(
			rimpay.ErrorCodeProviderError,
			"failed to decode status response",
			"bpay",
			false,
		)
	}

	// Convert to standard response
	status := &rimpay.TransactionStatus{
		TransactionID:     checkResp.TransactionID,
		Status:            convertTransactionStatus(checkResp.Status),
		Reference:         transactionID,
		ProviderReference: checkResp.TransactionID,
		Message:           checkResp.ErrorMessage,
		LastUpdated:       time.Now(),
		ProviderData: map[string]interface{}{
			"error_code":     checkResp.ErrorCode,
			"error_message":  checkResp.ErrorMessage,
			"status":         checkResp.Status,
			"transaction_id": checkResp.TransactionID,
		},
	}

	return status, nil
}

// convertLanguage converts rimpay.Language to B-PAY format
func convertLanguage(lang rimpay.Language) string {
	switch lang {
	case rimpay.LanguageEnglish:
		return "EN"
	case rimpay.LanguageFrench:
		return "FR"
	case rimpay.LanguageArabic:
		return "AR"
	default:
		return "FR"
	}
}
