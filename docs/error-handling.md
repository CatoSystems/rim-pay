# Error Handling

RimPay provides comprehensive error handling with specific error types that help you understand what went wrong and how to handle different failure scenarios.

## Error Types

### ValidationError

Occurs when input validation fails (phone numbers, amounts, etc.):

```go
type ValidationError struct {
    Field   string // The field that failed validation
    Message string // Human-readable error message
    Value   string // The invalid value
}
```

**Example:**
```go
request := &rimpay.BPayPaymentRequest{
    PhoneNumber: invalidPhone,
    // ... other fields
}

err := request.Validate()
if validationErr, ok := err.(*rimpay.ValidationError); ok {
    fmt.Printf("Validation failed for field '%s': %s\n", 
        validationErr.Field, validationErr.Message)
}
```

### ProviderError

Occurs when a payment provider returns an error:

```go
type ProviderError struct {
    Provider string // Provider name (bpay, masrvi)
    Code     string // Provider-specific error code
    Message  string // Error message from provider
}
```

**Example:**
```go
response, err := client.ProcessBPayPayment(ctx, request)
if providerErr, ok := err.(*rimpay.ProviderError); ok {
    switch providerErr.Code {
    case "INSUFFICIENT_FUNDS":
        fmt.Println("Customer needs to add funds to their account")
    case "INVALID_PIN":
        fmt.Println("Customer entered wrong PIN")
    case "ACCOUNT_BLOCKED":
        fmt.Println("Customer account is blocked")
    }
}
```

### NetworkError

Occurs when network communication fails:

```go
type NetworkError struct {
    Message   string // Error description
    Retryable bool   // Whether the error is retryable
    Cause     error  // Underlying error
}
```

**Example:**
```go
response, err := client.ProcessBPayPayment(ctx, request)
if networkErr, ok := err.(*rimpay.NetworkError); ok {
    if networkErr.Retryable {
        fmt.Println("Network error occurred, but will be retried automatically")
    } else {
        fmt.Println("Permanent network error:", networkErr.Message)
    }
}
```

### AuthenticationError

Occurs when provider authentication fails:

```go
type AuthenticationError struct {
    Provider string // Provider name
    Message  string // Authentication failure reason
}
```

**Example:**
```go
response, err := client.ProcessBPayPayment(ctx, request)
if authErr, ok := err.(*rimpay.AuthenticationError); ok {
    fmt.Printf("Authentication failed with %s: %s\n", 
        authErr.Provider, authErr.Message)
    // May need to update credentials
}
```

## Error Handling Patterns

### Pattern 1: Type-Based Error Handling

```go
func handlePaymentError(err error) {
    switch e := err.(type) {
    case *rimpay.ValidationError:
        // Handle validation errors - usually client-side issues
        log.Printf("Validation error in field %s: %s", e.Field, e.Message)
        // Return 400 Bad Request to client
        
    case *rimpay.AuthenticationError:
        // Handle authentication errors - configuration issue
        log.Printf("Auth error with %s: %s", e.Provider, e.Message)
        // Alert operations team, return 503 Service Unavailable
        
    case *rimpay.ProviderError:
        // Handle provider-specific errors
        log.Printf("Provider %s error %s: %s", e.Provider, e.Code, e.Message)
        handleProviderSpecificError(e)
        
    case *rimpay.NetworkError:
        // Handle network errors
        if e.Retryable {
            log.Printf("Retryable network error: %s", e.Message)
            // Error will be retried automatically
        } else {
            log.Printf("Permanent network error: %s", e.Message)
            // Return 502 Bad Gateway
        }
        
    default:
        // Handle unknown errors
        log.Printf("Unknown error: %v", err)
        // Return 500 Internal Server Error
    }
}
```

### Pattern 2: Error Code-Based Handling

```go
func handleProviderSpecificError(err *rimpay.ProviderError) {
    switch err.Provider {
    case "bpay":
        handleBPayError(err)
    case "masrvi":
        handleMasrviError(err)
    }
}

func handleBPayError(err *rimpay.ProviderError) {
    switch err.Code {
    case "INSUFFICIENT_FUNDS":
        // Customer needs to add money
        notifyCustomer("Please add funds to your mobile money account")
        
    case "INVALID_PIN":
        // Wrong PIN entered
        notifyCustomer("Incorrect PIN. Please try again")
        
    case "ACCOUNT_BLOCKED":
        // Account is blocked
        notifyCustomer("Your account is temporarily blocked. Please contact support")
        
    case "TRANSACTION_LIMIT_EXCEEDED":
        // Daily/monthly limit exceeded
        notifyCustomer("Transaction limit exceeded. Please try a smaller amount")
        
    case "SERVICE_UNAVAILABLE":
        // Provider service is down
        notifyCustomer("Payment service temporarily unavailable. Please try again later")
        
    default:
        // Unknown B-PAY error
        log.Printf("Unknown B-PAY error: %s - %s", err.Code, err.Message)
        notifyCustomer("Payment failed. Please try again or contact support")
    }
}
```

### Pattern 3: Graceful Degradation

```go
func processPaymentWithFallback(client *rimpay.Client, ctx context.Context) (*rimpay.PaymentResponse, error) {
    // Try primary provider (B-PAY)
    bpayRequest := createBPayRequest()
    response, err := client.ProcessBPayPayment(ctx, bpayRequest)
    
    if err == nil {
        return response, nil
    }
    
    // Check if error is recoverable
    if isRecoverableError(err) {
        log.Printf("B-PAY failed with recoverable error, trying MASRVI: %v", err)
        
        // Convert to MASRVI request and try alternative provider
        masrviRequest := convertToMasrviRequest(bpayRequest)
        return client.ProcessMasrviPayment(ctx, masrviRequest)
    }
    
    // Non-recoverable error
    return nil, err
}

func isRecoverableError(err error) bool {
    switch e := err.(type) {
    case *rimpay.NetworkError:
        return !e.Retryable // If not retryable by library, try different provider
    case *rimpay.ProviderError:
        // Some provider errors might be worth trying with different provider
        return e.Code == "SERVICE_UNAVAILABLE" || e.Code == "TIMEOUT"
    case *rimpay.AuthenticationError:
        return false // Auth errors are not recoverable with different provider
    default:
        return false
    }
}
```

## Error Logging and Monitoring

### Structured Logging

```go
import (
    "go.uber.org/zap"
    "github.com/CatoSystems/rim-pay/pkg/rimpay"
)

func logPaymentError(logger *zap.Logger, err error, request interface{}) {
    switch e := err.(type) {
    case *rimpay.ValidationError:
        logger.Warn("Payment validation failed",
            zap.String("error_type", "validation"),
            zap.String("field", e.Field),
            zap.String("message", e.Message),
            zap.String("value", e.Value),
        )
        
    case *rimpay.ProviderError:
        logger.Error("Provider error",
            zap.String("error_type", "provider"),
            zap.String("provider", e.Provider),
            zap.String("code", e.Code),
            zap.String("message", e.Message),
        )
        
    case *rimpay.NetworkError:
        logger.Error("Network error",
            zap.String("error_type", "network"),
            zap.Bool("retryable", e.Retryable),
            zap.String("message", e.Message),
            zap.Error(e.Cause),
        )
        
    case *rimpay.AuthenticationError:
        logger.Error("Authentication error",
            zap.String("error_type", "authentication"),
            zap.String("provider", e.Provider),
            zap.String("message", e.Message),
        )
        
    default:
        logger.Error("Unknown payment error",
            zap.String("error_type", "unknown"),
            zap.Error(err),
        )
    }
}
```

### Metrics Collection

```go
import "github.com/prometheus/client_golang/prometheus"

var (
    paymentErrors = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "rimpay_payment_errors_total",
            Help: "Total number of payment errors by type and provider",
        },
        []string{"error_type", "provider", "error_code"},
    )
)

func recordPaymentError(err error) {
    switch e := err.(type) {
    case *rimpay.ValidationError:
        paymentErrors.WithLabelValues("validation", "", e.Field).Inc()
        
    case *rimpay.ProviderError:
        paymentErrors.WithLabelValues("provider", e.Provider, e.Code).Inc()
        
    case *rimpay.NetworkError:
        retryable := "false"
        if e.Retryable {
            retryable = "true"
        }
        paymentErrors.WithLabelValues("network", "", retryable).Inc()
        
    case *rimpay.AuthenticationError:
        paymentErrors.WithLabelValues("authentication", e.Provider, "").Inc()
        
    default:
        paymentErrors.WithLabelValues("unknown", "", "").Inc()
    }
}
```

## Error Recovery Strategies

### Retry Logic

RimPay includes automatic retry for certain types of errors:

```go
// Automatic retry configuration
config.Retry = rimpay.RetryConfig{
    MaxAttempts:        3,
    InitialDelay:       1 * time.Second,
    MaxDelay:          10 * time.Second,
    BackoffMultiplier: 2.0,
}
```

**Automatically retried errors:**
- Network timeouts
- Connection errors  
- HTTP 5xx server errors
- Authentication token expiration

**Not automatically retried:**
- Validation errors
- Business logic errors (insufficient funds, etc.)
- HTTP 4xx client errors

### Manual Retry

For errors not automatically retried, you can implement manual retry:

```go
func retryPayment(client *rimpay.Client, ctx context.Context, request *rimpay.BPayPaymentRequest, maxRetries int) (*rimpay.PaymentResponse, error) {
    var lastErr error
    
    for attempt := 0; attempt <= maxRetries; attempt++ {
        response, err := client.ProcessBPayPayment(ctx, request)
        if err == nil {
            return response, nil
        }
        
        lastErr = err
        
        // Check if error is worth retrying
        if !shouldRetry(err) {
            break
        }
        
        if attempt < maxRetries {
            // Wait before retry with exponential backoff
            delay := time.Duration(attempt+1) * time.Second
            time.Sleep(delay)
        }
    }
    
    return nil, lastErr
}

func shouldRetry(err error) bool {
    switch e := err.(type) {
    case *rimpay.NetworkError:
        return e.Retryable
    case *rimpay.ProviderError:
        // Retry on service unavailable or timeout
        return e.Code == "SERVICE_UNAVAILABLE" || e.Code == "TIMEOUT"
    default:
        return false
    }
}
```

## Testing Error Scenarios

### Unit Tests

```go
func TestErrorHandling(t *testing.T) {
    tests := []struct {
        name          string
        setupError    func() error
        expectedType  string
        expectedField string
    }{
        {
            name: "validation error",
            setupError: func() error {
                return &rimpay.ValidationError{
                    Field:   "phone_number",
                    Message: "invalid phone number",
                    Value:   "invalid",
                }
            },
            expectedType:  "validation",
            expectedField: "phone_number",
        },
        // More test cases...
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := tt.setupError()
            
            switch e := err.(type) {
            case *rimpay.ValidationError:
                assert.Equal(t, tt.expectedField, e.Field)
            }
        })
    }
}
```

### Integration Tests

```go
func TestPaymentErrorScenarios(t *testing.T) {
    client := createTestClient()
    ctx := context.Background()
    
    t.Run("invalid phone number", func(t *testing.T) {
        request := &rimpay.BPayPaymentRequest{
            PhoneNumber: createInvalidPhone(), // This will cause validation error
            Amount:      validAmount,
            Reference:   "TEST-123",
            Passcode:    "1234",
        }
        
        _, err := client.ProcessBPayPayment(ctx, request)
        assert.Error(t, err)
        
        validationErr, ok := err.(*rimpay.ValidationError)
        assert.True(t, ok)
        assert.Equal(t, "phone_number", validationErr.Field)
    })
}
```

## Best Practices

### 1. Always Handle Errors

```go
// ❌ Don't ignore errors
response, _ := client.ProcessBPayPayment(ctx, request)

// ✅ Always handle errors
response, err := client.ProcessBPayPayment(ctx, request)
if err != nil {
    return handlePaymentError(err)
}
```

### 2. Provide Meaningful Error Messages

```go
// ❌ Generic error message
return errors.New("payment failed")

// ✅ Specific, actionable error message
if validationErr, ok := err.(*rimpay.ValidationError); ok {
    return fmt.Errorf("invalid %s: %s", validationErr.Field, validationErr.Message)
}
```

### 3. Log Errors Appropriately

```go
// ❌ Log sensitive information
log.Printf("Payment failed: %+v", request) // Contains passcode!

// ✅ Log safely without sensitive data
log.Printf("Payment failed for reference %s: %v", request.Reference, err)
```

### 4. Use Circuit Breakers for Provider Failures

```go
type CircuitBreaker struct {
    failureCount int
    threshold    int
    timeout      time.Duration
    lastFailure  time.Time
}

func (cb *CircuitBreaker) Call(fn func() error) error {
    if cb.isOpen() {
        return errors.New("circuit breaker is open")
    }
    
    err := fn()
    if err != nil {
        cb.recordFailure()
    } else {
        cb.recordSuccess()
    }
    
    return err
}
```

## Next Steps

- Learn about [retry mechanisms](retry-mechanisms.md)
- Explore [security considerations](security.md)
- Check out [performance tuning](performance.md)
