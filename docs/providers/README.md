# Payment Providers

RimPay supports multiple payment providers in Mauritania, each with its own strengths and use cases.

## Supported Providers

| Provider | Type | Authentication | Status Checking | Webhooks |
|----------|------|----------------|-----------------|----------|
| [B-PAY](bpay.md) | Mobile Money | OAuth 2.0 | ✅ | ❌ |
| [MASRVI](masrvi.md) | Web Payment | API Key | ❌ | ✅ |

## Provider Selection Guide

### Choose B-PAY when:
- You need mobile money integration
- Your customers primarily use mobile wallets
- You need real-time payment status checking
- You're building mobile applications
- Your target audience is mobile-first

### Choose MASRVI when:
- You need web-based payment processing
- You want webhook notifications
- You're building e-commerce platforms
- You need callback URL support
- Your customers prefer web payments

## Multi-Provider Architecture

RimPay is designed to support multiple providers simultaneously:

```go
// Configure multiple providers
config.Providers["bpay"] = bpayConfig
config.Providers["masrvi"] = masrviConfig

// Use type-safe methods for each provider
bpayResponse, err := client.ProcessBPayPayment(ctx, bpayRequest)
masrviResponse, err := client.ProcessMasrviPayment(ctx, masrviRequest)
```

## Provider-Specific Features

### Request Types
Each provider has its own request type with provider-specific fields:

- **BPayPaymentRequest**: Includes `Passcode` field for mobile money PIN
- **MasrviPaymentRequest**: Includes `CallbackURL` and `ReturnURL` fields

### Validation Rules
Each provider has specific validation requirements:

- **B-PAY**: Requires passcode, supports all Mauritanian operators
- **MASRVI**: Requires callback URLs, supports web payment flows

### Error Handling
Provider-specific error codes and messages:

- **B-PAY**: Mobile money specific errors (insufficient funds, wrong PIN)
- **MASRVI**: Web payment specific errors (session timeout, redirect issues)

## Provider Comparison

### Technical Comparison

| Feature | B-PAY | MASRVI |
|---------|-------|--------|
| **Authentication** | OAuth 2.0 with automatic token refresh | API Key based |
| **Request Method** | HTTP POST with JSON | HTTP POST with JSON |
| **Response Format** | JSON with transaction details | JSON with session details |
| **Status Checking** | Real-time status API | Webhook notifications |
| **Timeout Handling** | Configurable (default 30s) | Configurable (default 45s) |
| **Retry Support** | Automatic with exponential backoff | Automatic with exponential backoff |

### Business Comparison

| Aspect | B-PAY | MASRVI |
|--------|-------|--------|
| **Target Market** | Mobile money users | Web payment users |
| **User Experience** | Mobile-first | Web-first |
| **Integration Complexity** | Medium (OAuth) | Low (API Key) |
| **Settlement Time** | Instant | 1-3 business days |
| **Transaction Fees** | Competitive mobile rates | Competitive web rates |
| **Customer Support** | Mobile money focused | Web payment focused |

## Provider Integration Patterns

### Single Provider Integration

```go
// Simple B-PAY only integration
config := rimpay.DefaultConfig()
config.DefaultProvider = "bpay"
config.Providers["bpay"] = bpayConfig

client, err := rimpay.NewClient(config)
response, err := client.ProcessBPayPayment(ctx, bpayRequest)
```

### Multi-Provider with Fallback

```go
// Try B-PAY first, fallback to MASRVI
func processPaymentWithFallback(client *rimpay.Client, ctx context.Context, 
    bpayReq *rimpay.BPayPaymentRequest, masrviReq *rimpay.MasrviPaymentRequest) (*rimpay.PaymentResponse, error) {
    
    // Try B-PAY first
    response, err := client.ProcessBPayPayment(ctx, bpayReq)
    if err == nil {
        return response, nil
    }
    
    // If B-PAY fails, try MASRVI
    log.Printf("B-PAY failed, trying MASRVI: %v", err)
    return client.ProcessMasrviPayment(ctx, masrviReq)
}
```

### Provider Selection by Amount

```go
func selectProvider(amount *money.Money) string {
    // Use B-PAY for small amounts (mobile money)
    if amount.Cents() < 50000 { // Less than 500 MRU
        return "bpay"
    }
    // Use MASRVI for larger amounts (web payment)
    return "masrvi"
}
```

## Provider Status and Monitoring

### Health Checks

```go
func checkProviderHealth(client *rimpay.Client, ctx context.Context) map[string]bool {
    status := make(map[string]bool)
    
    // Check B-PAY
    if err := client.CheckBPayHealth(ctx); err == nil {
        status["bpay"] = true
    } else {
        status["bpay"] = false
    }
    
    // Check MASRVI  
    if err := client.CheckMasrviHealth(ctx); err == nil {
        status["masrvi"] = true
    } else {
        status["masrvi"] = false
    }
    
    return status
}
```

### Circuit Breaker Pattern

```go
type ProviderCircuitBreaker struct {
    failures map[string]int
    maxFailures int
    resetTime time.Duration
    lastFailure map[string]time.Time
}

func (cb *ProviderCircuitBreaker) IsOpen(provider string) bool {
    failures := cb.failures[provider]
    lastFailure := cb.lastFailure[provider]
    
    if failures >= cb.maxFailures {
        if time.Since(lastFailure) > cb.resetTime {
            // Reset circuit breaker
            cb.failures[provider] = 0
            return false
        }
        return true // Circuit is open
    }
    return false
}
```

## Best Practices

### 1. Provider Configuration
- Always configure providers for your target environment
- Use environment variables for credentials
- Set appropriate timeouts for each provider

### 2. Error Handling
- Handle provider-specific errors appropriately
- Implement fallback strategies for critical payments
- Log provider performance and errors

### 3. Testing
- Test with both providers in sandbox environment
- Validate provider-specific fields and requirements
- Monitor provider response times and success rates

### 4. Monitoring
- Track provider success rates
- Monitor provider response times
- Set up alerts for provider failures

## Migration Guide

### From Single Provider to Multi-Provider

1. **Update Configuration**:
   ```go
   // Old: Single provider
   config.DefaultProvider = "bpay"
   
   // New: Multi-provider
   config.Providers["bpay"] = bpayConfig
   config.Providers["masrvi"] = masrviConfig
   ```

2. **Update Request Creation**:
   ```go
   // Old: Generic request
   request := &rimpay.PaymentRequest{...}
   
   // New: Provider-specific requests
   bpayRequest := &rimpay.BPayPaymentRequest{...}
   masrviRequest := &rimpay.MasrviPaymentRequest{...}
   ```

3. **Update Processing**:
   ```go
   // Old: Generic processing
   response, err := client.ProcessPayment(ctx, request)
   
   // New: Provider-specific processing
   bpayResponse, err := client.ProcessBPayPayment(ctx, bpayRequest)
   masrviResponse, err := client.ProcessMasrviPayment(ctx, masrviRequest)
   ```

## Next Steps

- Learn about [B-PAY provider](bpay.md) specifics
- Learn about [MASRVI provider](masrvi.md) specifics  
- Compare [provider features](comparison.md) in detail
