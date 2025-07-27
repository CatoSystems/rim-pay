# Frequently Asked Questions (FAQ)

## General Questions

### What is RimPay?

RimPay is a comprehensive Go library for payment processing in Mauritania. It provides a unified API for integrating with multiple payment providers including B-PAY and MASRVI, with built-in support for phone number validation, money handling, error management, and retry mechanisms.

### Which programming language is RimPay written in?

RimPay is written in Go (Golang) and requires Go 1.19 or later.

### Is RimPay free to use?

Yes, RimPay is open-source software released under the MIT License. You can use it freely in both personal and commercial projects.

### Which payment providers does RimPay support?

Currently, RimPay supports:
- **B-PAY**: Mauritanian mobile money provider
- **MASRVI**: Web-based payment provider

More providers may be added in future releases.

## Installation and Setup

### How do I install RimPay?

```bash
go get github.com/CatoSystems/rim-pay@latest
```

### What are the system requirements?

- Go 1.19 or later
- Internet connection for API calls
- Valid credentials for your chosen payment provider(s)

### Do I need credentials for both providers?

No, you only need credentials for the providers you plan to use. You can configure just B-PAY, just MASRVI, or both.

### How do I get provider credentials?

- **B-PAY**: Contact B-PAY directly to set up a merchant account
- **MASRVI**: Contact MASRVI to obtain API credentials

Both providers typically offer sandbox environments for testing.

## Phone Number Validation

### What phone number formats are supported?

RimPay supports Mauritanian phone numbers in various formats:
- International: `+22233445566`
- National: `22233445566`
- Local: `33445566` (automatically prefixed with +222)

### Which operator prefixes are valid?

Valid Mauritanian mobile prefixes:
- **2**: Mauritel
- **3**: Chinguitel  
- **4**: Mattel

Prefixes 5, 6, 7, 8, 9 are not supported.

### Why is my phone number being rejected?

Common reasons:
- Invalid prefix (must be 2, 3, or 4)
- Wrong length (must be 8 digits after country code)
- Contains non-numeric characters
- Wrong country code (must be +222)

## Money Handling

### What currency does RimPay support?

RimPay currently supports MRU (Mauritanian Ouguiya), the current currency of Mauritania since January 1, 2018.

### How are amounts represented?

RimPay uses decimal-based arithmetic to avoid floating-point precision issues:
```go
// 100.50 MRU
amount := money.New(decimal.NewFromInt(10050), money.MRU)
```

### Can I use floating-point numbers for amounts?

While possible, it's not recommended due to precision issues:
```go
// Less precise
amount := money.FromFloat64(100.50, money.MRU)

// More precise (recommended)
amount := money.New(decimal.NewFromFloat(100.50), money.MRU)
```

## Payment Processing

### How do I process a B-PAY payment?

```go
request := &rimpay.BPayPaymentRequest{
    Amount:      money.New(decimal.NewFromInt(10000), money.MRU),
    PhoneNumber: phone,
    Reference:   "ORDER-123",
    Description: "Payment description",
    Passcode:    "1234", // Customer's mobile money PIN
}

response, err := client.ProcessBPayPayment(ctx, request)
```

### How do I process a MASRVI payment?

```go
request := &rimpay.MasrviPaymentRequest{
    Amount:      money.New(decimal.NewFromInt(10000), money.MRU),
    PhoneNumber: phone,
    Reference:   "ORDER-123",
    Description: "Payment description",
    CallbackURL: "https://yoursite.com/webhook",
    ReturnURL:   "https://yoursite.com/return",
}

response, err := client.ProcessMasrviPayment(ctx, request)
```

### How do I check payment status?

For B-PAY (real-time status checking):
```go
status, err := client.GetPaymentStatus(ctx, transactionID)
```

For MASRVI (webhook-based):
```go
// Set up webhook endpoint to receive status updates
// Status checking not available via API
```

### What's the maximum transaction amount?

Transaction limits depend on:
- Your provider agreement
- Customer account limits
- Regulatory limits
- Daily/monthly transaction limits

Check with your payment provider for specific limits.

## Error Handling

### What types of errors can occur?

RimPay provides specific error types:
- **ValidationError**: Input validation failures
- **ProviderError**: Payment provider errors
- **NetworkError**: Connection/network issues
- **AuthenticationError**: Provider authentication failures

### How do I handle different error types?

```go
response, err := client.ProcessBPayPayment(ctx, request)
if err != nil {
    switch e := err.(type) {
    case *rimpay.ValidationError:
        // Handle validation error
        fmt.Printf("Validation error in field %s: %s", e.Field, e.Message)
    case *rimpay.ProviderError:
        // Handle provider error
        fmt.Printf("Provider %s error %s: %s", e.Provider, e.Code, e.Message)
    // ... other error types
    }
}
```

### Are payments automatically retried?

Yes, RimPay automatically retries certain types of failures:
- Network timeouts
- Connection errors
- HTTP 5xx server errors
- Authentication token expiration

Business logic errors (insufficient funds, wrong PIN) are not retried.

## Configuration

### How do I configure multiple providers?

```go
config := rimpay.DefaultConfig()

// Configure B-PAY
config.Providers["bpay"] = rimpay.ProviderConfig{
    Enabled: true,
    BaseURL: "https://api.bpay.mr",
    Credentials: map[string]string{
        "username": "your_bpay_username",
        "password": "your_bpay_password",
    },
}

// Configure MASRVI
config.Providers["masrvi"] = rimpay.ProviderConfig{
    Enabled: true,
    BaseURL: "https://api.masrvi.mr",
    Credentials: map[string]string{
        "merchant_id": "your_merchant_id",
        "api_key":     "your_api_key",
    },
}
```

### How do I use environment variables for configuration?

```go
config.Providers["bpay"] = rimpay.ProviderConfig{
    Enabled: true,
    BaseURL: os.Getenv("BPAY_BASE_URL"),
    Credentials: map[string]string{
        "username": os.Getenv("BPAY_USERNAME"),
        "password": os.Getenv("BPAY_PASSWORD"),
    },
}
```

### What's the difference between sandbox and production?

- **Sandbox**: For development and testing, no real money
- **Production**: For live transactions with real money

Always test thoroughly in sandbox before going to production.

## Security

### How should I handle customer PINs?

- Never log or store PINs
- Use HTTPS for all communications
- Clear PIN data from memory after use
- Implement rate limiting for PIN attempts

```go
// Good practice
paymentRequest.Passcode = customerPIN
response, err := client.ProcessBPayPayment(ctx, paymentRequest)
paymentRequest.Passcode = "" // Clear immediately
```

### Is RimPay PCI compliant?

RimPay itself doesn't handle card data, but follows security best practices:
- No sensitive data logging
- Secure credential handling
- HTTPS enforcement
- Memory cleanup for sensitive data

### How do I secure my API credentials?

- Use environment variables, not hardcoded values
- Implement proper secret management
- Rotate credentials regularly
- Use different credentials for sandbox/production
- Monitor credential usage

## Performance

### How many concurrent payments can I process?

This depends on:
- Your server resources
- Provider rate limits
- Network capacity
- Client configuration

Start with 10-20 concurrent requests and adjust based on performance.

### How do I optimize performance?

- Use connection pooling
- Implement appropriate timeouts
- Monitor provider response times
- Use batch processing where possible
- Implement circuit breakers for provider failures

### What are typical response times?

Response times vary by provider and network conditions:
- **B-PAY**: Usually 2-10 seconds
- **MASRVI**: Usually 1-5 seconds for session creation

Monitor your specific use case for accurate metrics.

## Testing

### How do I test payments without real money?

Use sandbox environments:
```go
config.Environment = rimpay.EnvironmentSandbox
```

Both B-PAY and MASRVI offer sandbox environments for testing.

### How do I write unit tests?

```go
func TestPaymentProcessing(t *testing.T) {
    // Create test client
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
}
```

### Are there mock providers available?

Currently, RimPay doesn't include built-in mocks, but you can create your own test implementations of the provider interfaces.

## Deployment

### Can I use RimPay in production?

Yes, RimPay is designed for production use. Make sure to:
- Use production provider credentials
- Set `Environment` to `production`
- Implement proper logging and monitoring
- Test thoroughly in sandbox first

### How do I monitor RimPay in production?

- Log payment attempts and results
- Monitor error rates by type and provider
- Track response times
- Set up alerts for high error rates
- Monitor provider availability

### What about high availability?

- Use multiple provider configurations
- Implement circuit breakers
- Set up health checks
- Use load balancing
- Plan for graceful degradation

## Troubleshooting

### Why am I getting authentication errors?

Common causes:
- Incorrect credentials
- Expired credentials
- Wrong environment (sandbox vs production)
- Network connectivity issues
- Provider service issues

### Why are payments failing?

Check for:
- Invalid phone numbers
- Insufficient customer funds
- Incorrect PINs
- Provider service issues
- Network problems
- Configuration errors

### How do I debug payment issues?

1. Enable detailed logging
2. Check error types and messages
3. Verify configuration
4. Test in sandbox
5. Check provider status pages
6. Contact provider support if needed

### Why are responses slow?

Possible causes:
- Network latency
- Provider performance issues
- Incorrect timeout settings
- Server resource constraints
- High concurrent load

## Migration and Updates

### How do I upgrade RimPay?

```bash
go get github.com/CatoSystems/rim-pay@latest
go mod tidy
```

Check the [CHANGELOG](../CHANGELOG.md) for breaking changes.

### Will my code break with updates?

RimPay follows semantic versioning:
- Patch versions (0.1.1): Bug fixes, no breaking changes
- Minor versions (0.2.0): New features, backward compatible
- Major versions (1.0.0): May include breaking changes

### How do I migrate from single provider to multi-provider?

See the [migration guide](api/README.md#migration-guide) in the API documentation.

## Getting Help

### Where can I get support?

1. Check this FAQ
2. Read the [documentation](README.md)
3. Look at [examples](examples/README.md)
4. Search [GitHub issues](https://github.com/CatoSystems/rim-pay/issues)
5. Create a new issue with detailed information

### How do I report bugs?

1. Check if it's already reported
2. Create a GitHub issue with:
   - Clear description
   - Steps to reproduce
   - Expected vs actual behavior
   - RimPay version
   - Go version
   - Sample code (without sensitive data)

### How do I request features?

Create a GitHub issue with:
- Clear feature description
- Use case and benefits
- Proposed API (if applicable)
- Willingness to contribute

### Is commercial support available?

Currently, support is community-based through GitHub issues. For commercial support needs, contact the maintainers.

## Contributing

### How can I contribute?

- Report bugs
- Suggest features
- Improve documentation
- Submit pull requests
- Help with testing

See [CONTRIBUTING.md](../CONTRIBUTING.md) for detailed guidelines.

### What kind of contributions are welcome?

- Bug fixes
- New provider integrations
- Performance improvements
- Documentation improvements
- Test coverage improvements
- Example applications

### How do I add a new payment provider?

1. Implement the `PaymentProvider` interface
2. Create provider-specific request types
3. Add comprehensive tests
4. Update documentation
5. Submit a pull request

See the existing provider implementations for examples.
