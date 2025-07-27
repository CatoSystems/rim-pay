# Examples Documentation

This directory contains comprehensive examples demonstrating various aspects of the RimPay library.

## Available Examples

### Basic Examples

| Example | Description | Key Features |
|---------|-------------|--------------|
| [basic-usage](../examples/basic-usage/) | Simple payment processing | Client creation, basic B-PAY payment |
| [bpay](../examples/bpay/) | B-PAY specific features | OAuth authentication, status checking |
| [masrvi](../examples/masrvi/) | MASRVI specific features | API key auth, webhooks |

### Advanced Examples

| Example | Description | Key Features |
|---------|-------------|--------------|
| [multi-provider](../examples/multi-provider/) | Multiple providers | Provider switching, fallback strategies |
| [error-handling](../examples/error-handling/) | Error scenarios | Type-safe error handling, recovery |
| [configuration](../examples/configuration/) | Advanced config | Environment variables, multi-provider setup |

## Running Examples

### Prerequisites

```bash
# Install RimPay
go get github.com/CatoSystems/rim-pay@latest

# Navigate to examples directory
cd examples/
```

### Run Individual Examples

```bash
# Run basic usage example
cd basic-usage && go run main.go

# Run B-PAY example
cd ../bpay && go run main.go

# Run MASRVI example
cd ../masrvi && go run main.go
```

### Run All Examples

```bash
# From examples directory
make run-all

# Or manually
make basic-usage
make bpay
make masrvi
make multi-provider
make error-handling
make configuration
```

## Example Breakdown

### 1. Basic Usage (`basic-usage/`)

**What it demonstrates:**
- Client configuration
- Simple B-PAY payment
- Error handling basics

**Key code:**
```go
// Configure client
config := rimpay.DefaultConfig()
config.DefaultProvider = "bpay"
client, err := rimpay.NewClient(config)

// Create payment
request := &rimpay.BPayPaymentRequest{
    Amount:      money.New(decimal.NewFromInt(10000), money.MRU),
    PhoneNumber: phone,
    Reference:   "ORDER-123",
    Passcode:    "1234",
}

// Process payment
response, err := client.ProcessBPayPayment(ctx, request)
```

### 2. B-PAY Provider (`bpay/`)

**What it demonstrates:**
- B-PAY OAuth 2.0 authentication
- Provider availability checking
- Payment status monitoring
- Different operator support

**Key features:**
- Automatic token management
- Real-time status checking
- Operator-specific examples
- Error code handling

### 3. MASRVI Provider (`masrvi/`)

**What it demonstrates:**
- MASRVI API key authentication
- Webhook URL configuration
- Session-based payments
- Callback handling

**Key features:**
- Session management
- Webhook integration
- Return URL handling
- Web payment flows

### 4. Multi-Provider (`multi-provider/`)

**What it demonstrates:**
- Configuring multiple providers
- Provider selection strategies
- Fallback mechanisms
- Type-safe provider methods

**Strategies shown:**
- Amount-based provider selection
- Primary/fallback provider pattern
- Provider availability checking
- Load balancing approaches

### 5. Error Handling (`error-handling/`)

**What it demonstrates:**
- All error types and handling
- Error recovery strategies
- Logging and monitoring
- Retry mechanisms

**Error scenarios:**
- Validation errors
- Authentication failures
- Provider errors
- Network issues

### 6. Configuration (`configuration/`)

**What it demonstrates:**
- Environment-based configuration
- Credentials management
- Multi-environment setup
- Best practices

**Configuration methods:**
- Environment variables
- Configuration structs
- Provider-specific settings
- Security considerations

## Integration Patterns

### E-commerce Integration

```go
// Example: Order checkout with payment
func processOrderPayment(orderID string, customerPhone string, amount decimal.Decimal) error {
    // Create payment request
    phone, err := phone.NewPhone(customerPhone)
    if err != nil {
        return fmt.Errorf("invalid phone: %w", err)
    }
    
    request := &rimpay.BPayPaymentRequest{
        Amount:      money.New(amount, money.MRU),
        PhoneNumber: phone,
        Reference:   fmt.Sprintf("ORDER-%s", orderID),
        Description: fmt.Sprintf("Payment for order %s", orderID),
        Passcode:    getCustomerPIN(customerPhone), // Get from secure input
    }
    
    // Process payment
    response, err := client.ProcessBPayPayment(ctx, request)
    if err != nil {
        return handlePaymentError(err)
    }
    
    // Update order status
    return updateOrderStatus(orderID, response.TransactionID, "paid")
}
```

### Subscription Service Integration

```go
// Example: Recurring payment processing
func processSubscriptionPayment(subscriptionID string, amount decimal.Decimal) error {
    subscription := getSubscription(subscriptionID)
    
    request := &rimpay.BPayPaymentRequest{
        Amount:      money.New(amount, money.MRU),
        PhoneNumber: subscription.Phone,
        Reference:   fmt.Sprintf("SUB-%s-%s", subscriptionID, time.Now().Format("200601")),
        Description: "Monthly subscription payment",
        Passcode:    subscription.StoredPIN, // Securely stored
    }
    
    response, err := client.ProcessBPayPayment(ctx, request)
    if err != nil {
        // Handle subscription payment failure
        return handleSubscriptionPaymentError(subscriptionID, err)
    }
    
    // Update subscription
    return updateSubscriptionPayment(subscriptionID, response.TransactionID)
}
```

### Marketplace Integration

```go
// Example: Split payment processing
func processMarketplacePayment(orderID string, vendorID string, platformFee decimal.Decimal, vendorAmount decimal.Decimal) error {
    totalAmount := platformFee.Add(vendorAmount)
    
    // Process full payment first
    request := &rimpay.BPayPaymentRequest{
        Amount:      money.New(totalAmount, money.MRU),
        PhoneNumber: getCustomerPhone(orderID),
        Reference:   fmt.Sprintf("MARKETPLACE-%s", orderID),
        Description: "Marketplace purchase",
        Passcode:    getCustomerPIN(orderID),
    }
    
    response, err := client.ProcessBPayPayment(ctx, request)
    if err != nil {
        return err
    }
    
    // Process vendor payout (simplified)
    return processVendorPayout(vendorID, vendorAmount, response.TransactionID)
}
```

## Testing Examples

### Unit Testing

```go
func TestPaymentProcessing(t *testing.T) {
    // Create test client with mock provider
    client := createTestClient()
    
    request := &rimpay.BPayPaymentRequest{
        Amount:      money.New(decimal.NewFromInt(1000), money.MRU),
        PhoneNumber: createTestPhone(),
        Reference:   "TEST-123",
        Passcode:    "1234",
    }
    
    response, err := client.ProcessBPayPayment(context.Background(), request)
    
    assert.NoError(t, err)
    assert.NotEmpty(t, response.TransactionID)
    assert.Equal(t, rimpay.PaymentStatusSuccess, response.Status)
}
```

### Integration Testing

```go
func TestRealProviderIntegration(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test")
    }
    
    // Use sandbox configuration
    config := createSandboxConfig()
    client, err := rimpay.NewClient(config)
    require.NoError(t, err)
    
    // Test with small amount
    request := &rimpay.BPayPaymentRequest{
        Amount:      money.New(decimal.NewFromInt(100), money.MRU), // 1 MRU
        PhoneNumber: getSandboxPhone(),
        Reference:   fmt.Sprintf("INTEGRATION-TEST-%d", time.Now().Unix()),
        Passcode:    getSandboxPIN(),
    }
    
    response, err := client.ProcessBPayPayment(context.Background(), request)
    
    // Note: This might fail in sandbox, but should not panic
    if err != nil {
        t.Logf("Integration test failed (expected in sandbox): %v", err)
    } else {
        t.Logf("Integration test succeeded: %s", response.TransactionID)
    }
}
```

## Performance Examples

### Concurrent Processing

```go
func processConcurrentPayments(requests []*rimpay.BPayPaymentRequest) []PaymentResult {
    const maxConcurrency = 10
    
    semaphore := make(chan struct{}, maxConcurrency)
    results := make(chan PaymentResult, len(requests))
    
    for _, request := range requests {
        go func(req *rimpay.BPayPaymentRequest) {
            semaphore <- struct{}{} // Acquire
            defer func() { <-semaphore }() // Release
            
            response, err := client.ProcessBPayPayment(ctx, req)
            results <- PaymentResult{
                Request:  req,
                Response: response,
                Error:    err,
            }
        }(request)
    }
    
    // Collect results
    var paymentResults []PaymentResult
    for i := 0; i < len(requests); i++ {
        paymentResults = append(paymentResults, <-results)
    }
    
    return paymentResults
}
```

### Batch Processing

```go
func processBatchPayments(batch []PaymentItem) error {
    const batchSize = 50
    
    for i := 0; i < len(batch); i += batchSize {
        end := i + batchSize
        if end > len(batch) {
            end = len(batch)
        }
        
        currentBatch := batch[i:end]
        if err := processBatch(currentBatch); err != nil {
            return fmt.Errorf("batch %d-%d failed: %w", i, end, err)
        }
        
        // Rate limiting
        time.Sleep(100 * time.Millisecond)
    }
    
    return nil
}
```

## Monitoring Examples

### Metrics Collection

```go
import "github.com/prometheus/client_golang/prometheus"

var (
    paymentDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "rimpay_payment_duration_seconds",
            Help: "Payment processing duration",
        },
        []string{"provider", "status"},
    )
    
    paymentCounter = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "rimpay_payments_total",
            Help: "Total payments processed",
        },
        []string{"provider", "status"},
    )
)

func processPaymentWithMetrics(request *rimpay.BPayPaymentRequest) (*rimpay.PaymentResponse, error) {
    start := time.Now()
    
    response, err := client.ProcessBPayPayment(ctx, request)
    
    duration := time.Since(start).Seconds()
    status := "success"
    if err != nil {
        status = "failure"
    }
    
    paymentDuration.WithLabelValues("bpay", status).Observe(duration)
    paymentCounter.WithLabelValues("bpay", status).Inc()
    
    return response, err
}
```

## Security Examples

### PCI DSS Compliance

```go
// Example: Secure PIN handling
func processPaymentSecurely(request PaymentData) error {
    // Never log sensitive data
    logger.Info("Processing payment", 
        zap.String("reference", request.Reference),
        zap.String("amount", request.Amount.String()),
        // DO NOT log: request.Passcode
    )
    
    // Create request with secure PIN
    paymentRequest := &rimpay.BPayPaymentRequest{
        Amount:      request.Amount,
        PhoneNumber: request.PhoneNumber,
        Reference:   request.Reference,
        Passcode:    request.Passcode, // Handle securely
    }
    
    response, err := client.ProcessBPayPayment(ctx, paymentRequest)
    
    // Clear sensitive data immediately
    paymentRequest.Passcode = ""
    request.Passcode = ""
    
    return err
}
```

### Input Validation

```go
func validatePaymentInput(input PaymentInput) error {
    // Validate amount
    if input.Amount <= 0 {
        return errors.New("amount must be positive")
    }
    
    if input.Amount > maxTransactionAmount {
        return errors.New("amount exceeds maximum limit")
    }
    
    // Validate phone
    if _, err := phone.NewPhone(input.PhoneNumber); err != nil {
        return fmt.Errorf("invalid phone number: %w", err)
    }
    
    // Validate reference
    if len(input.Reference) < 5 || len(input.Reference) > 50 {
        return errors.New("reference must be 5-50 characters")
    }
    
    // Validate PIN format (without revealing actual value)
    if len(input.Passcode) < 4 || len(input.Passcode) > 8 {
        return errors.New("invalid PIN format")
    }
    
    return nil
}
```

## Next Steps

- [Run the examples](../examples/README.md)
- [Explore integration patterns](integration-patterns.md)
- [Learn about testing strategies](testing.md)
