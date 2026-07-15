package bpay

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/CatoSystems/rim-pay/internal/providers/common"
	"github.com/CatoSystems/rim-pay/pkg/money"
	"github.com/CatoSystems/rim-pay/pkg/phone"
	"github.com/CatoSystems/rim-pay/pkg/rimpay"
)

type passcodeTestLogger struct{}

func (passcodeTestLogger) Debug(string, ...interface{}) {}
func (passcodeTestLogger) Info(string, ...interface{})  {}
func (passcodeTestLogger) Warn(string, ...interface{})  {}
func (passcodeTestLogger) Error(string, ...interface{}) {}

// routingStub returns an auth token for the auth endpoint and a canned payment
// response for the payment endpoint, capturing the payment request body.
type routingStub struct {
	capturedPayment *common.HTTPRequest
}

func (s *routingStub) Do(req *common.HTTPRequest) (*common.HTTPResponse, error) {
	if strings.Contains(req.URL, "/authentification") {
		return &common.HTTPResponse{
			StatusCode: 200,
			Body:       []byte(`{"access_token":"test-token","expires_in":"3600"}`),
		}, nil
	}
	s.capturedPayment = req
	return &common.HTTPResponse{
		StatusCode: 200,
		Body:       []byte(`{"errorCode":"0","errorMessage":"","transactionId":"TX-1"}`),
	}, nil
}

func TestBPayForwardsCallerPasscode(t *testing.T) {
	phoneNum, err := phone.NewPhone("+22220000000")
	if err != nil {
		t.Fatalf("failed to create phone: %v", err)
	}

	stub := &routingStub{}
	config := rimpay.ProviderConfig{
		BaseURL:     "https://example.test",
		Credentials: map[string]string{"username": "u", "password": "p", "client_id": "e-bankily"},
		Timeout:     5 * time.Second,
	}
	auth := NewAuthManager(config, stub, passcodeTestLogger{})
	pp := NewPaymentProcessor(config, stub, auth, passcodeTestLogger{})

	const callerPasscode = "4321"
	req := &rimpay.PaymentRequest{
		PhoneNumber: phoneNum,
		Amount:      money.FromFloat64(50.00, money.MRU),
		Reference:   "REF-1",
		Passcode:    callerPasscode,
		Language:    rimpay.LanguageFrench,
	}

	resp, err := pp.ProcessPayment(context.Background(), req)
	if err != nil {
		t.Fatalf("ProcessPayment failed: %v", err)
	}

	// The outgoing payment request must carry the caller's passcode verbatim.
	if stub.capturedPayment == nil {
		t.Fatal("payment request was never sent")
	}
	var sent PaymentRequest
	if err := json.Unmarshal(stub.capturedPayment.Body, &sent); err != nil {
		t.Fatalf("failed to decode sent body: %v", err)
	}
	if sent.Passcode != callerPasscode {
		t.Errorf("sent passcode = %q, want %q", sent.Passcode, callerPasscode)
	}

	// The passcode must NOT be echoed back in response metadata.
	if _, leaked := resp.Metadata["passcode"]; leaked {
		t.Error("passcode must not appear in response metadata")
	}
}
