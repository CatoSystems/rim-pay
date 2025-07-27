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

## RimPay Examples

This directory contains comprehensive examples demonstrating how to use the RimPay library for payment processing in Mauritania.

## Structure

Each example is now organized in its own package to avoid conflicts and make them easier to run independently:

```
examples/
├── README.md                     # This file
├── basic-usage/                  # Basic usage patterns
│   └── main.go
├── bpay/                         # B-PAY specific examples
│   └── main.go
├── masrvi/                       # MASRVI specific examples
│   └── main.go
├── multi-provider/               # Using multiple providers
│   └── main.go
├── error-handling/               # Error handling and retry
│   └── main.go
├── configuration/                # Configuration examples
│   └── main.go
└── [legacy files...]            # Old single-file examples
```

## Running Examples

Each example can be run independently using Go:

```bash
# Basic usage example
cd examples/basic-usage && go run main.go

# B-PAY specific example
cd examples/bpay && go run main.go

# MASRVI specific example
cd examples/masrvi && go run main.go

# Multi-provider example
cd examples/multi-provider && go run main.go

# Error handling example
cd examples/error-handling && go run main.go

# Configuration example
cd examples/configuration && go run main.go
```

## Example Descriptions

### 1. Basic Usage (`basic-usage/`)
- Shows fundamental payment processing
- Demonstrates both B-PAY and MASRVI workflows
- Basic error handling
- Configuration setup

### 2. B-PAY Examples (`bpay/`)
- B-PAY specific authentication (OAuth 2.0)
- Payment processing for all Mauritanian operators
- Transaction status checking
- B-PAY specific error handling

### 3. MASRVI Examples (`masrvi/`)
- MASRVI session management
- Payment form generation
- Webhook notification handling
- E-commerce integration patterns

### 4. Multi-Provider (`multi-provider/`)
- Using multiple providers in one application
- Dynamic provider selection logic
- Type-safe provider-specific requests
- Unified error handling across providers

### 5. Error Handling (`error-handling/`)
- Comprehensive error handling patterns
- Automatic retry mechanisms
- Error classification and recovery
- Context timeout handling

### 6. Configuration (`configuration/`)
- Default vs custom configurations
- Environment-specific settings
- Timeout and connection management
- Logging configuration

## Prerequisites

Before running any examples, ensure you have:

1. **Go 1.19+** installed
2. **RimPay library** installed:
   ```bash
   go mod init your-project
   go get github.com/CatoSystems/rim-pay
   ```
3. **Provider credentials** (for actual payments):
   - B-PAY: username, password, client_id
   - MASRVI: merchant_id

## Configuration

Update the credentials in each example file before running:

```go
// B-PAY Configuration
"bpay": {
    Enabled: true,
    BaseURL: "https://ebankily-tst.appspot.com", // Sandbox
    Credentials: map[string]string{
        "username":  "your_bpay_username",
        "password":  "your_bpay_password", 
        "client_id": "your_bpay_client_id",
    },
},

// MASRVI Configuration
"masrvi": {
    Enabled: true,
    BaseURL: "https://masrviapp.mr/online", // Production
    Credentials: map[string]string{
        "merchant_id": "your_masrvi_merchant_id",
    },
},
```

## Phone Number Format

All examples use valid Mauritanian phone numbers with prefixes 2, 3, or 4:
- `22334455` - Mauritel
- `33445566` - Chinguitel
- `44556677` - Mattel

## Testing

The examples are configured to work with sandbox/test environments. For production use:

1. Update the `BaseURL` in provider configurations
2. Use production credentials
3. Set `Environment: rimpay.EnvironmentProduction`

## Troubleshooting

### Common Issues

1. **"Provider not available"** - Check credentials and network connectivity
2. **"Invalid phone number"** - Ensure phone numbers use valid Mauritanian prefixes (2, 3, 4)
3. **"Payment failed"** - Check provider-specific error messages and logs

### Debug Mode

Enable debug logging by setting:
```go
config.Logging.Level = "debug"
```

This will show detailed HTTP requests/responses and authentication flows.

## Contributing

When adding new examples:

1. Create a new directory under `examples/`
2. Add a `main.go` file with a focused example
3. Update this README with the new example description
4. Test the example thoroughly

## Support

For questions about these examples:
1. Check the main [README.md](../README.md)
2. Review the [API documentation](../pkg/rimpay/)
3. Create an issue on GitHub

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
