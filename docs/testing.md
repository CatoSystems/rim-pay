# Testing Guide

This guide covers testing strategies, utilities, and best practices for RimPay.

## Overview

RimPay uses a comprehensive testing approach including:
- Unit tests for individual components
- Integration tests with provider sandboxes
- End-to-end tests for complete workflows
- Performance tests for scalability validation

## Test Structure

```
tests/
├── unit/           # Unit tests
│   ├── client/     # Client logic tests
│   ├── providers/  # Provider implementation tests
│   ├── validation/ # Validation logic tests
│   └── utils/      # Utility function tests
├── integration/    # Integration tests
│   ├── bpay/       # B-PAY integration tests
│   ├── masrvi/     # MASRVI integration tests
│   └── combined/   # Multi-provider tests
├── e2e/           # End-to-end tests
│   ├── scenarios/  # Complete payment scenarios
│   └── workflows/  # Business workflow tests
└── performance/   # Performance and load tests
    ├── load/      # Load testing
    └── stress/    # Stress testing
```

## Unit Testing

### Testing Client Logic

```go
package client_test

import (
    "context"
    "testing"
    "time"

    "github.com/CatoSystems/rim-pay/pkg/rimpay"
    "github.com/CatoSystems/rim-pay/pkg/money"
    "github.com/CatoSystems/rim-pay/pkg/phone"
    "github.com/shopspring/decimal"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestClient_ProcessBPayPayment(t *testing.T) {
    tests := []struct {
        name    string
        request *rimpay.BPayPaymentRequest
        wantErr bool
        errType string
    }{
        {
            name: "valid_payment",
            request: &rimpay.BPayPaymentRequest{
                Amount:      money.New(decimal.NewFromInt(10000), money.MRU),
                PhoneNumber: mustParsePhone("+22233445566"),
                Reference:   "TEST-123",
                Description: "Test payment",
                Passcode:    "1234",
            },
            wantErr: false,
        },
        {
            name: "invalid_phone",
            request: &rimpay.BPayPaymentRequest{
                Amount:      money.New(decimal.NewFromInt(10000), money.MRU),
                PhoneNumber: mustParsePhone("+22255445566"), // Invalid prefix
                Reference:   "TEST-123",
                Description: "Test payment",
                Passcode:    "1234",
            },
            wantErr: true,
            errType: "validation",
        },
        {
            name: "negative_amount",
            request: &rimpay.BPayPaymentRequest{
                Amount:      money.New(decimal.NewFromInt(-1000), money.MRU),
                PhoneNumber: mustParsePhone("+22233445566"),
                Reference:   "TEST-123",
                Description: "Test payment",
                Passcode:    "1234",
            },
            wantErr: true,
            errType: "validation",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            client := createTestClient()
            ctx := context.Background()

            response, err := client.ProcessBPayPayment(ctx, tt.request)

            if tt.wantErr {
                assert.Error(t, err)
                if tt.errType != "" {
                    switch tt.errType {
                    case "validation":
                        assert.True(t, rimpay.IsValidationError(err))
                    case "provider":
                        assert.True(t, rimpay.IsProviderError(err))
                    }
                }
            } else {
                assert.NoError(t, err)
                assert.NotNil(t, response)
                assert.NotEmpty(t, response.TransactionID)
                assert.Equal(t, tt.request.Reference, response.Reference)
            }
        })
    }
}

func createTestClient() *rimpay.Client {
    config := rimpay.DefaultConfig()
    config.Environment = rimpay.EnvironmentSandbox
    
    config.Providers["bpay"] = rimpay.ProviderConfig{
        Enabled: true,
        BaseURL: "https://sandbox-api.bpay.mr",
        Credentials: map[string]string{
            "username": "test_username",
            "password": "test_password",
        },
    }

    client, err := rimpay.New(config)
    if err != nil {
        panic(err)
    }
    return client
}

func mustParsePhone(s string) phone.Number {
    p, err := phone.Parse(s)
    if err != nil {
        panic(err)
    }
    return p
}
```

### Testing Provider Implementations

```go
package bpay_test

import (
    "context"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/CatoSystems/rim-pay/internal/providers/bpay"
    "github.com/stretchr/testify/assert"
)

func TestBPayProvider_ProcessPayment(t *testing.T) {
    tests := []struct {
        name           string
        mockResponse   string
        expectedStatus int
        wantErr        bool
    }{
        {
            name: "successful_payment",
            mockResponse: `{
                "status": "success",
                "transaction_id": "TXN123456",
                "amount": "100.00",
                "message": "Payment successful"
            }`,
            expectedStatus: http.StatusOK,
            wantErr:        false,
        },
        {
            name: "insufficient_funds",
            mockResponse: `{
                "status": "error",
                "error_code": "INSUFFICIENT_FUNDS",
                "message": "Insufficient funds in account"
            }`,
            expectedStatus: http.StatusBadRequest,
            wantErr:        true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                w.Header().Set("Content-Type", "application/json")
                w.WriteHeader(tt.expectedStatus)
                w.Write([]byte(tt.mockResponse))
            }))
            defer server.Close()

            provider := createTestBPayProvider(server.URL)
            request := createTestBPayRequest()

            response, err := provider.ProcessPayment(context.Background(), request)

            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
                assert.NotNil(t, response)
                assert.Equal(t, "TXN123456", response.TransactionID)
            }
        })
    }
}
```

### Testing Validation Logic

```go
package validation_test

import (
    "testing"

    "github.com/CatoSystems/rim-pay/pkg/phone"
    "github.com/stretchr/testify/assert"
)

func TestPhoneValidation(t *testing.T) {
    validNumbers := []string{
        "+22233445566",
        "22233445566",
        "33445566",
        "+22244556677",
        "+22222334455",
    }

    invalidNumbers := []string{
        "+22155445566",  // Invalid prefix (5)
        "+22166445566",  // Invalid prefix (6)
        "+223334455",    // Too short
        "+2223344556677", // Too long
        "+33233445566",  // Wrong country code
        "abc33445566",   // Non-numeric
        "",              // Empty
    }

    t.Run("valid_numbers", func(t *testing.T) {
        for _, number := range validNumbers {
            t.Run(number, func(t *testing.T) {
                phone, err := phone.Parse(number)
                assert.NoError(t, err)
                assert.True(t, phone.IsValid())
            })
        }
    })

    t.Run("invalid_numbers", func(t *testing.T) {
        for _, number := range invalidNumbers {
            t.Run(number, func(t *testing.T) {
                _, err := phone.Parse(number)
                assert.Error(t, err)
            })
        }
    })
}
```

## Integration Testing

### Provider Integration Tests

```go
package integration_test

import (
    "context"
    "os"
    "testing"
    "time"

    "github.com/CatoSystems/rim-pay/pkg/rimpay"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestBPayIntegration(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test in short mode")
    }

    // Skip if credentials not available
    username := os.Getenv("BPAY_TEST_USERNAME")
    password := os.Getenv("BPAY_TEST_PASSWORD")
    if username == "" || password == "" {
        t.Skip("B-PAY test credentials not available")
    }

    config := rimpay.DefaultConfig()
    config.Environment = rimpay.EnvironmentSandbox
    config.Providers["bpay"] = rimpay.ProviderConfig{
        Enabled: true,
        BaseURL: "https://sandbox-api.bpay.mr",
        Credentials: map[string]string{
            "username": username,
            "password": password,
        },
    }

    client, err := rimpay.New(config)
    require.NoError(t, err)

    t.Run("successful_payment", func(t *testing.T) {
        request := createValidBPayRequest()
        ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
        defer cancel()

        response, err := client.ProcessBPayPayment(ctx, request)
        
        assert.NoError(t, err)
        assert.NotNil(t, response)
        assert.NotEmpty(t, response.TransactionID)
        assert.Equal(t, request.Reference, response.Reference)
        
        // Verify payment status
        status, err := client.GetPaymentStatus(ctx, response.TransactionID)
        assert.NoError(t, err)
        assert.NotNil(t, status)
    })

    t.Run("invalid_pin", func(t *testing.T) {
        request := createValidBPayRequest()
        request.Passcode = "9999" // Invalid PIN
        
        ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
        defer cancel()

        _, err := client.ProcessBPayPayment(ctx, request)
        
        assert.Error(t, err)
        assert.True(t, rimpay.IsProviderError(err))
    })
}
```

### Multi-Provider Integration Tests

```go
func TestMultiProviderIntegration(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test in short mode")
    }

    client := createMultiProviderTestClient()

    tests := []struct {
        name     string
        provider string
        testFunc func(t *testing.T, client *rimpay.Client)
    }{
        {
            name:     "bpay_payment",
            provider: "bpay",
            testFunc: testBPayPayment,
        },
        {
            name:     "masrvi_payment",
            provider: "masrvi",
            testFunc: testMasrviPayment,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            tt.testFunc(t, client)
        })
    }
}
```

## End-to-End Testing

### Complete Payment Workflows

```go
package e2e_test

import (
    "context"
    "testing"
    "time"

    "github.com/CatoSystems/rim-pay/pkg/rimpay"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/suite"
)

type PaymentWorkflowSuite struct {
    suite.Suite
    client *rimpay.Client
}

func (s *PaymentWorkflowSuite) SetupSuite() {
    s.client = createE2ETestClient()
}

func (s *PaymentWorkflowSuite) TestCompletePaymentWorkflow() {
    ctx := context.Background()

    // Step 1: Process payment
    request := createValidPaymentRequest()
    response, err := s.client.ProcessBPayPayment(ctx, request)
    s.Require().NoError(err)
    s.Require().NotNil(response)

    transactionID := response.TransactionID
    s.NotEmpty(transactionID)

    // Step 2: Check initial status
    status, err := s.client.GetPaymentStatus(ctx, transactionID)
    s.Require().NoError(err)
    s.NotNil(status)

    // Step 3: Wait for final status (with timeout)
    finalStatus := s.waitForFinalStatus(ctx, transactionID, 60*time.Second)
    s.True(rimpay.IsPaymentFinal(finalStatus.Status))
}

func (s *PaymentWorkflowSuite) waitForFinalStatus(ctx context.Context, transactionID string, timeout time.Duration) *rimpay.PaymentStatus {
    ctx, cancel := context.WithTimeout(ctx, timeout)
    defer cancel()

    ticker := time.NewTicker(5 * time.Second)
    defer ticker.Stop()

    for {
        select {
        case <-ctx.Done():
            s.Fail("Timeout waiting for final status")
            return nil
        case <-ticker.C:
            status, err := s.client.GetPaymentStatus(ctx, transactionID)
            s.Require().NoError(err)
            
            if rimpay.IsPaymentFinal(status.Status) {
                return status
            }
        }
    }
}

func TestPaymentWorkflowSuite(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping E2E test in short mode")
    }
    suite.Run(t, new(PaymentWorkflowSuite))
}
```

## Performance Testing

### Load Testing

```go
package performance_test

import (
    "context"
    "sync"
    "testing"
    "time"

    "github.com/CatoSystems/rim-pay/pkg/rimpay"
    "github.com/stretchr/testify/assert"
)

func TestPaymentThroughput(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping performance test in short mode")
    }

    client := createPerformanceTestClient()
    concurrency := 10
    requestsPerWorker := 50
    totalRequests := concurrency * requestsPerWorker

    var wg sync.WaitGroup
    results := make(chan testResult, totalRequests)

    startTime := time.Now()

    // Start workers
    for i := 0; i < concurrency; i++ {
        wg.Add(1)
        go func(workerID int) {
            defer wg.Done()
            
            for j := 0; j < requestsPerWorker; j++ {
                result := processTestPayment(client, workerID, j)
                results <- result
            }
        }(i)
    }

    // Wait for completion
    go func() {
        wg.Wait()
        close(results)
    }()

    // Collect results
    var successful, failed int
    var totalDuration time.Duration

    for result := range results {
        totalDuration += result.duration
        if result.err != nil {
            failed++
        } else {
            successful++
        }
    }

    endTime := time.Now()
    totalTestTime := endTime.Sub(startTime)

    // Analyze results
    throughput := float64(totalRequests) / totalTestTime.Seconds()
    averageResponseTime := totalDuration / time.Duration(totalRequests)
    successRate := float64(successful) / float64(totalRequests) * 100

    t.Logf("Performance Test Results:")
    t.Logf("  Total Requests: %d", totalRequests)
    t.Logf("  Successful: %d", successful)
    t.Logf("  Failed: %d", failed)
    t.Logf("  Success Rate: %.2f%%", successRate)
    t.Logf("  Throughput: %.2f requests/second", throughput)
    t.Logf("  Average Response Time: %v", averageResponseTime)
    t.Logf("  Total Test Time: %v", totalTestTime)

    // Assertions
    assert.GreaterOrEqual(t, successRate, 95.0, "Success rate should be at least 95%")
    assert.GreaterOrEqual(t, throughput, 5.0, "Throughput should be at least 5 requests/second")
    assert.LessOrEqual(t, averageResponseTime, 10*time.Second, "Average response time should be under 10 seconds")
}

type testResult struct {
    duration time.Duration
    err      error
}

func processTestPayment(client *rimpay.Client, workerID, requestID int) testResult {
    start := time.Now()
    
    request := createUniquePaymentRequest(workerID, requestID)
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    _, err := client.ProcessBPayPayment(ctx, request)
    
    return testResult{
        duration: time.Since(start),
        err:      err,
    }
}
```

### Memory and Resource Testing

```go
func TestMemoryUsage(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping memory test in short mode")
    }

    var m1, m2 runtime.MemStats
    runtime.GC()
    runtime.ReadMemStats(&m1)

    client := createTestClient()
    
    // Process many payments
    for i := 0; i < 1000; i++ {
        request := createTestPaymentRequest(i)
        ctx := context.Background()
        
        _, err := client.ProcessBPayPayment(ctx, request)
        if err != nil {
            // Continue on error for memory testing
            continue
        }
    }

    runtime.GC()
    runtime.ReadMemStats(&m2)

    memoryIncrease := m2.Alloc - m1.Alloc
    t.Logf("Memory increase: %d bytes", memoryIncrease)
    
    // Assert reasonable memory usage (adjust threshold as needed)
    assert.Less(t, memoryIncrease, uint64(10*1024*1024), "Memory usage should be less than 10MB")
}
```

## Test Utilities

### Mock Providers

```go
package testutils

import (
    "context"
    
    "github.com/CatoSystems/rim-pay/pkg/rimpay"
)

type MockProvider struct {
    responses map[string]*rimpay.PaymentResponse
    errors    map[string]error
}

func NewMockProvider() *MockProvider {
    return &MockProvider{
        responses: make(map[string]*rimpay.PaymentResponse),
        errors:    make(map[string]error),
    }
}

func (m *MockProvider) SetResponse(reference string, response *rimpay.PaymentResponse) {
    m.responses[reference] = response
}

func (m *MockProvider) SetError(reference string, err error) {
    m.errors[reference] = err
}

func (m *MockProvider) ProcessPayment(ctx context.Context, request rimpay.PaymentRequest) (*rimpay.PaymentResponse, error) {
    // Implementation for mock responses
    reference := getReference(request)
    
    if err, exists := m.errors[reference]; exists {
        return nil, err
    }
    
    if response, exists := m.responses[reference]; exists {
        return response, nil
    }
    
    // Default successful response
    return &rimpay.PaymentResponse{
        TransactionID: "MOCK-" + reference,
        Status:        rimpay.StatusSuccess,
        Reference:     reference,
    }, nil
}
```

### Test Data Generators

```go
package testutils

import (
    "fmt"
    "time"

    "github.com/CatoSystems/rim-pay/pkg/money"
    "github.com/CatoSystems/rim-pay/pkg/phone"
    "github.com/CatoSystems/rim-pay/pkg/rimpay"
    "github.com/shopspring/decimal"
)

func CreateValidBPayRequest() *rimpay.BPayPaymentRequest {
    phone, _ := phone.Parse("+22233445566")
    amount := money.New(decimal.NewFromInt(10000), money.MRU)
    
    return &rimpay.BPayPaymentRequest{
        Amount:      amount,
        PhoneNumber: phone,
        Reference:   fmt.Sprintf("TEST-%d", time.Now().Unix()),
        Description: "Test payment",
        Passcode:    "1234",
    }
}

func CreateValidMasrviRequest() *rimpay.MasrviPaymentRequest {
    phone, _ := phone.Parse("+22233445566")
    amount := money.New(decimal.NewFromInt(10000), money.MRU)
    
    return &rimpay.MasrviPaymentRequest{
        Amount:      amount,
        PhoneNumber: phone,
        Reference:   fmt.Sprintf("TEST-%d", time.Now().Unix()),
        Description: "Test payment",
        CallbackURL: "https://example.com/webhook",
        ReturnURL:   "https://example.com/return",
    }
}

func CreateTestPhoneNumbers() []string {
    return []string{
        "+22233445566",
        "+22244556677",
        "+22222334455",
    }
}

func CreateTestAmounts() []money.Money {
    return []money.Money{
        money.New(decimal.NewFromInt(1000), money.MRU),    // 10.00 MRU
        money.New(decimal.NewFromInt(5000), money.MRU),    // 50.00 MRU
        money.New(decimal.NewFromInt(10000), money.MRU),   // 100.00 MRU
    }
}
```

## Running Tests

### Basic Test Commands

```bash
# Run all tests
go test ./...

# Run only unit tests
go test ./... -short

# Run with coverage
go test ./... -cover

# Run with detailed coverage
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out

# Run specific test
go test ./pkg/rimpay -run TestClient_ProcessBPayPayment

# Run with verbose output
go test ./... -v
```

### Integration Test Setup

```bash
# Set environment variables for integration tests
export BPAY_TEST_USERNAME="your_test_username"
export BPAY_TEST_PASSWORD="your_test_password"
export MASRVI_TEST_MERCHANT_ID="your_test_merchant_id"
export MASRVI_TEST_API_KEY="your_test_api_key"

# Run integration tests
go test ./tests/integration/... -v

# Run E2E tests
go test ./tests/e2e/... -v
```

### Performance Test Commands

```bash
# Run performance tests
go test ./tests/performance/... -v

# Run with memory profiling
go test ./tests/performance/... -memprofile=mem.prof

# Run with CPU profiling
go test ./tests/performance/... -cpuprofile=cpu.prof

# Analyze profiles
go tool pprof mem.prof
go tool pprof cpu.prof
```

## CI/CD Integration

### GitHub Actions Example

```yaml
name: Tests

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    
    services:
      redis:
        image: redis:6
        ports:
          - 6379:6379
    
    steps:
    - uses: actions/checkout@v3
    
    - name: Set up Go
      uses: actions/setup-go@v3
      with:
        go-version: 1.19
    
    - name: Run unit tests
      run: go test ./... -short -race -coverprofile=coverage.out
    
    - name: Run integration tests
      env:
        BPAY_TEST_USERNAME: ${{ secrets.BPAY_TEST_USERNAME }}
        BPAY_TEST_PASSWORD: ${{ secrets.BPAY_TEST_PASSWORD }}
      run: go test ./tests/integration/... -v
    
    - name: Upload coverage
      uses: codecov/codecov-action@v3
      with:
        file: ./coverage.out
```

This comprehensive testing guide ensures that RimPay is thoroughly tested at all levels, from individual units to complete end-to-end workflows, with proper performance validation and CI/CD integration.
