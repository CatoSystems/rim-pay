# RimPay Examples

This directory contains comprehensive examples demonstrating how to use the RimPay library for mobile payment processing in Mauritania.

## Available Examples

### 1. Provider-Specific API Examples (`provider_specific_example.go`)
**NEW**: Demonstrates the improved type-safe API with provider-specific request types.

- **B-PAY Specific Payments**: Using `BPayPaymentRequest` with provider-specific fields like `Passcode`
- **MASRVI Specific Payments**: Using `MasrviPaymentRequest` with `CallbackURL` and `ReturnURL`
- **Type Safety**: Compile-time safety preventing mismatched provider requests
- **Multi-Provider Support**: Managing multiple providers with type safety

**Key Features:**
- Provider-specific request types (`BPayPaymentRequest`, `MasrviPaymentRequest`)
- Type-safe API preventing provider mismatches
- Improved developer experience with IntelliSense support
- Provider-specific validation and error handling

### 2. Basic Usage (`basic_usage.go`)
Simple introduction to the RimPay library with basic payment processing.

- Client initialization
- Provider configuration
- Simple payment processing
- Error handling

### 3. B-PAY Provider Example (`bpay_example.go`)
Comprehensive B-PAY implementation with OAuth 2.0 authentication.

- OAuth 2.0 authentication flow
- Provider availability checking
- Operator-specific payments (Mauritel, Mattel, Chinguitel)
- Payment status checking
- Error handling with retries

### 4. MASRVI Provider Example (`masrvi_example.go`)
MASRVI implementation with webhook server and session management.

- Session-based authentication
- Webhook notification server
- Payment form creation
- Notification handling and processing
- HTTP server setup for webhook endpoints

### 5. Error Handling (`error_handling_example.go`)
Comprehensive error handling patterns and recovery strategies.

- Different error types (network, authentication, validation)
- Error recovery strategies
- Logging and monitoring
- Custom error handling logic

### 6. Configuration Examples (`configuration_example.go`)
Various configuration patterns and best practices.

- Environment-based configuration
- Multiple provider setup
- Timeout and retry configuration
- Logging configuration
- Security best practices

### 7. Multi-Provider Support (`multi_provider_example.go`)
Managing multiple payment providers with fallback strategies.

- Multiple provider registration
- Intelligent provider selection
- Fallback mechanisms
- Load balancing strategies
- Provider health monitoring

### 8. Retry Mechanisms (`retry_demo.go`)
Advanced retry patterns for handling transient failures.

- Exponential backoff strategies
- Custom retry conditions
- Circuit breaker patterns
- Retry metrics and monitoring
- Timeout handling

## Running the Examples

### Prerequisites
```bash
# Install dependencies
go mod tidy

# Set environment variables (if needed)
export BPAY_CLIENT_ID="your_client_id"
export BPAY_CLIENT_SECRET="your_client_secret"
export MASRVI_MERCHANT_ID="your_merchant_id"
```

### Running Individual Examples

```bash
# Provider-specific API (recommended starting point)
go run examples/provider_specific_example.go

# Basic usage
go run examples/basic_usage.go

# B-PAY specific example
go run examples/bpay_example.go

# MASRVI with webhook server
go run examples/masrvi_example.go

# Error handling patterns
go run examples/error_handling_example.go

# Configuration examples
go run examples/configuration_example.go

# Multi-provider setup
go run examples/multi_provider_example.go

# Retry mechanisms
go run examples/retry_demo.go
```

### Running All Examples
```bash
# Build and run all examples
make examples

# Or manually
for example in examples/*.go; do
    echo "Running $example..."
    go run "$example"
    echo "---"
done
```

## Example Output

When running the provider-specific example, you'll see output like:

```
=== RimPay Provider-Specific API Examples ===

--- B-PAY Specific Payment Example ---
B-PAY Payment Response:
  Transaction ID: bpay_tx_1234567890
  Status: pending
  Reference: BPAY-TEST-001
  Provider: bpay

--- MASRVI Specific Payment Example ---
MASRVI Payment Response:
  Transaction ID: masrvi_tx_0987654321
  Status: pending
  Reference: MASRVI-TEST-001
  Provider: masrvi
  Payment URL: https://pay.masrvi.mr/payment/xyz

--- Multi-Provider Type-Safe Example ---
Multi-Provider B-PAY Result: bpay_tx_1111111111 (Status: pending)
Multi-Provider MASRVI Result: masrvi_tx_2222222222 (Status: pending)
✅ Multi-provider payments completed with type safety!
```

## Key Concepts

### Provider-Specific Types (NEW)
The library now provides type-safe, provider-specific request types:

```go
// B-PAY specific request
bpayRequest := &bpay.BPayPaymentRequest{
    PhoneNumber: phone.MauritanianPhone{Number: "22123456"},
    Amount:      money.NewMRU(1000),
    Passcode:    "1234",      // B-PAY specific
    // ... other fields
}

// MASRVI specific request  
masrviRequest := &masrvi.MasrviPaymentRequest{
    PhoneNumber: phone.MauritanianPhone{Number: "33987654"},
    Amount:      money.NewMRU(500),
    CallbackURL: "https://myapp.com/webhook", // MASRVI specific
    ReturnURL:   "https://myapp.com/return",  // MASRVI specific
    // ... other fields
}
```

### Type Safety Benefits
- **Compile-time Error Prevention**: Cannot mix provider request types
- **IntelliSense Support**: Better IDE support with provider-specific fields
- **Validation**: Provider-specific validation rules
- **Documentation**: Clear API contracts for each provider

### Migration from Generic API
The generic `PaymentRequest` is still supported but deprecated:

```go
// Old way (still works, but deprecated)
request := &rimpay.PaymentRequest{...}
response, err := client.ProcessPayment(ctx, request)

// New way (recommended)
bpayRequest := &bpay.BPayPaymentRequest{...}
response, err := client.ProcessBPayPayment(ctx, bpayRequest)
```

## Best Practices

1. **Use Provider-Specific Types**: Always use the new provider-specific request types for better type safety
2. **Error Handling**: Always handle errors appropriately with retry logic
3. **Validation**: Validate requests before processing
4. **Logging**: Implement proper logging for debugging and monitoring
5. **Configuration**: Use environment variables for sensitive configuration
6. **Testing**: Test with mock providers before production deployment
7. **Webhooks**: For MASRVI, always implement webhook handlers for status updates

## Testing

All examples are designed to work with mock/test configurations. For production use:

1. Replace test credentials with real ones
2. Update base URLs to production endpoints
3. Implement proper error handling
4. Set up monitoring and logging
5. Configure webhook endpoints (for MASRVI)

## Support

For questions or issues with the examples:
1. Check the main README.md in the project root
2. Review the internal documentation
3. Examine the test files for additional usage patterns
