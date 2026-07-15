package click

import (
	"context"
	"testing"
	"time"

	"github.com/CatoSystems/rim-pay/internal/providers/common"
	"github.com/CatoSystems/rim-pay/pkg/rimpay"
)

type stubHTTP struct {
	body string
	err  error
	last *common.HTTPRequest
}

func (s *stubHTTP) Do(req *common.HTTPRequest) (*common.HTTPResponse, error) {
	s.last = req
	if s.err != nil {
		return nil, s.err
	}
	return &common.HTTPResponse{StatusCode: 200, Body: []byte(s.body)}, nil
}

type nopLogger struct{}

func (nopLogger) Debug(string, ...interface{}) {}
func (nopLogger) Info(string, ...interface{})  {}
func (nopLogger) Warn(string, ...interface{})  {}
func (nopLogger) Error(string, ...interface{}) {}

func testConfig() rimpay.ProviderConfig {
	return rimpay.ProviderConfig{
		BaseURL:     "https://tagpay.test",
		Credentials: map[string]string{"merchant_id": "0896353536734538"},
		Timeout:     5 * time.Second,
	}
}

func TestSessionParsesOKPrefix(t *testing.T) {
	stub := &stubHTTP{body: "OK:27875690759565722269644474422394"}
	sm := NewSessionManager(testConfig(), stub, nopLogger{})
	id, err := sm.GetSessionID(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id != "27875690759565722269644474422394" {
		t.Errorf("session id = %q", id)
	}
}

func TestSessionRejectsNOK(t *testing.T) {
	stub := &stubHTTP{body: "NOK:UNKNOWN_MERCHANT"}
	sm := NewSessionManager(testConfig(), stub, nopLogger{})
	if _, err := sm.GetSessionID(context.Background()); err == nil {
		t.Error("expected error for NOK response")
	}
}
