# RimPay Library Examples

This directory contains comprehensive examples demonstrating how to use the RimPay payment library for processing payments through Mauritanian payment providers.

## 📋 Available Examples

### 1. 🏦 Basic Usage (`basic_usage.go`)
**What it demonstrates:**
- Creating a simple payment client
- Processing a basic payment
- Handling payment responses
- Error handling basics

**Key concepts:**
```go
// Create client with configuration
client, err := rimpay.NewClient(config)

// Process payment
response, err := client.ProcessPayment(ctx, request)
```

**Run:** `go run examples/basic_usage.go`

---

### 2. 💳 B-PAY Provider (`bpay_example.go`)
**What it demonstrates:**
- B-PAY specific configuration
- OAuth 2.0 authentication flow
- Payment processing for all Mauritanian operators (Mauritel, Mattel, Chinguitel)
- Transaction status checking
- Provider availability checking

**Key features:**
- Automatic token management
- Retry on authentication failures
- Operator-specific payment handling

**Run:** `go run examples/bpay_example.go`

---

### 3. 🌐 MASRVI Provider (`masrvi_example.go`)
**What it demonstrates:**
- MASRVI specific configuration
- Session management (5-minute validity)
- E-commerce payment form generation
- Webhook notification handling
- HTTP server for receiving notifications

**Key features:**
- Session-based authentication
- Payment form creation
- Real-time webhook processing
- Multiple redirect URLs

**Run:** `go run examples/masrvi_example.go`

---

### 4. 🔄 Multi-Provider (`multi_provider_example.go`)
**What it demonstrates:**
- Configuration with multiple providers
- Provider-specific payment routing
- Provider failover strategies
- Bulk payment processing
- Dynamic provider switching

**Key features:**
- Default and fallback providers
- Provider-specific request formatting
- Bulk payment orchestration

**Run:** `go run examples/multi_provider_example.go`

---

### 5. ⚠️ Error Handling (`error_handling_example.go`)
**What it demonstrates:**
- Network error handling with automatic retry
- Authentication error scenarios
- Validation error handling
- Business logic errors (insufficient funds, wrong PIN)
- Context timeout handling

**Key features:**
- Exponential backoff retry
- Error classification (retryable vs non-retryable)
- Timeout management
- Detailed error analysis

**Run:** `go run examples/error_handling_example.go`

---

### 6. ⚙️ Configuration (`configuration_example.go`)
**What it demonstrates:**
- Basic, production, and development configurations
- Environment-specific settings
- HTTP client tuning for performance
- Security configuration options
- Configuration validation

**Key features:**
- Environment-aware settings
- Performance optimization
- Security best practices
- Configuration validation

**Run:** `go run examples/configuration_example.go`

---

### 7. 🔄 Retry Mechanism (`retry_demo.go`)
**What it demonstrates:**
- Automatic retry functionality
- Exponential backoff with jitter
- Retry configuration options
- Retryable vs non-retryable errors

**Key features:**
- Smart retry logic
- Configurable retry policies
- Context cancellation support

**Run:** `go run examples/retry_demo.go`

---

## 🚀 Quick Start

1. **Install dependencies:**
   ```bash
   go mod tidy
   ```

2. **Update credentials:**
   Edit the example files and replace the placeholder credentials with your actual provider credentials:
   ```go
   Credentials: map[string]string{
       "username":  "your_actual_username",
       "password":  "your_actual_password", 
       "client_id": "your_actual_client_id",
   }
   ```

3. **Run an example:**
   ```bash
   go run examples/basic_usage.go
   ```

## 📖 Common Patterns

### Creating a Payment Request
```go
// Create phone number
phone, err := phone.NewPhone("22334455")
if err != nil {
    log.Fatal(err)
}

// Create amount
amount := money.New(decimal.NewFromFloat(100.00), money.MRU)

// Create request
request := &rimpay.PaymentRequest{
    Amount:      amount,
    PhoneNumber: phone,
    Reference:   "ORDER-123",
    Language:    rimpay.LanguageFrench,
    Passcode:    "1234", // For B-PAY
    Description: "Payment description",
}
```

### Error Handling
```go
response, err := client.ProcessPayment(ctx, request)
if err != nil {
    if paymentErr, ok := err.(*rimpay.PaymentError); ok {
        switch paymentErr.Code {
        case rimpay.ErrorCodeInsufficientFunds:
            // Handle insufficient funds
        case rimpay.ErrorCodeNetworkError:
            // Network error - will be retried automatically
        case rimpay.ErrorCodeAuthenticationFailed:
            // Check credentials
        }
    }
    return
}
```

### Provider Configuration
```go
config := &rimpay.Config{
    Environment:     rimpay.EnvironmentSandbox,
    DefaultProvider: "bpay",
    Providers: map[string]rimpay.ProviderConfig{
        "bpay": {
            Enabled: true,
            BaseURL: "https://ebankily-tst.appspot.com",
            Timeout: 30 * time.Second,
            Credentials: map[string]string{
                "username":  "your_username",
                "password":  "your_password",
                "client_id": "your_client_id",
            },
        },
    },
}
```

## 🔧 Testing

All examples can be run individually. For testing without real provider credentials:

1. **Use sandbox URLs** (already configured in examples)
2. **Mock responses** by replacing HTTP clients
3. **Test error scenarios** by using invalid credentials

## 📞 Mauritanian Phone Numbers

The library supports all Mauritanian operators:

| Operator | Prefixes | Example |
|----------|----------|---------|
| Mauritel | 2, 3, 4, 5 | 22334455 |
| Mattel   | 6, 7     | 66778899 |
| Chinguitel | 8, 9   | 88990011 |

Phone numbers can be provided in various formats:
- International: `+22222334455`
- With country code: `0022222334455`
- Local: `22334455`

## 💰 Money Handling

The library uses decimal precision for accurate money calculations:

```go
// Create from float
amount := money.FromFloat64(100.50, money.MRU)

// Create from string
amount, err := money.FromString("100.50", money.MRU)

// Create from cents
amount := money.FromCents(10050, money.MRU) // 100.50 MRU
```

## 🔐 Security Best Practices

1. **Never hardcode credentials** in production code
2. **Use environment variables** for sensitive data
3. **Enable TLS** in production
4. **Use request signing** when available
5. **Validate webhook signatures** for MASRVI notifications

## 🌍 Environment Configuration

| Environment | Use Case | Base URLs |
|-------------|----------|-----------|
| Sandbox | Development/Testing | Test provider URLs |
| Production | Live transactions | Production provider URLs |

## 📚 Additional Resources

- **API Documentation:** See PDF specifications in project root
- **Provider APIs:** B-PAY and MASRVI official documentation
- **Go Documentation:** `go doc` command for package details
- **Test Files:** `*_test.go` files for usage patterns

## 🆘 Troubleshooting

### Common Issues

1. **Authentication Failed**
   - Check credentials in provider configuration
   - Verify provider URLs are correct
   - Ensure credentials match environment (sandbox vs production)

2. **Network Errors**
   - Check internet connectivity
   - Verify provider service status
   - Increase timeout if needed

3. **Validation Errors**
   - Verify phone number format
   - Check amount is positive and not zero
   - Ensure reference is provided and valid

4. **Provider Not Available**
   - Check provider service status
   - Verify base URL is accessible
   - Test with simple HTTP request

### Getting Help

- Review error messages carefully
- Check provider-specific error codes
- Use debug logging for detailed information
- Test with minimal examples first

---

## 📝 Example Output

When you run the examples, you'll see detailed output showing:
- Configuration validation
- Payment processing steps
- Provider responses
- Error handling
- Retry attempts
- Success/failure notifications

This helps you understand exactly how the library works and how to integrate it into your applications.
