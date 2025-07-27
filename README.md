# RimPay

RimPay is a comprehensive Go library for payment processing in Mauritania, supporting multiple payment providers including B-PAY and MASRVI. The library provides a type-safe, provider-specific API for seamless payment integration.

## Features

- **Multi-Provider Support**: B-PAY and MASRVI payment providers
- **Type-Safe API**: Provider-specific request types with validation
- **Mauritanian Phone Validation**: Built-in validation for Mauritanian phone numbers
- **Money Handling**: Precise decimal-based money calculations
- **Retry Logic**: Configurable retry mechanisms with exponential backoff
- **Error Handling**: Comprehensive error types and handling
- **Configuration**: Flexible configuration system for different environments

## Installation

```bash
go get github.com/CatoSystems/rim-pay
```

## Quick Start

### Basic Setup

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"
    
    "github.com/CatoSystems/rim-pay/pkg/rimpay"
    "github.com/CatoSystems/rim-pay/pkg/phone"
    "github.com/CatoSystems/rim-pay/pkg/money"
    "github.com/shopspring/decimal"
)

func main() {
    // Create configuration
    config := rimpay.DefaultConfig()
    config.Environment = rimpay.EnvironmentSandbox
    config.DefaultProvider = "bpay"
    
    // Configure B-PAY provider
    config.Providers["bpay"] = rimpay.ProviderConfig{
        Enabled: true,
        BaseURL: "https://api.bpay.mr",
        Credentials: map[string]string{
            "username": "your_username",
            "password": "your_password",
        },
        Timeout: 30 * time.Second,
    }
    
    // Create client
    client, err := rimpay.NewClient(config)
    if err != nil {
        log.Fatal(err)
    }
    
    // Create payment request
    phone, _ := phone.NewPhone("+22233445566")
    amount := money.New(decimal.NewFromInt(10000), "MRU") // 100.00 MRU
    
    request := &rimpay.BPayPaymentRequest{
        Amount:      amount,
        PhoneNumber: phone,
        Reference:   "ORDER-12345",
        Description: "Payment for order 12345",
        Passcode:    "1234",
    }
    
    // Process payment
    ctx := context.Background()
    response, err := client.ProcessBPayPayment(ctx, request)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Payment successful! Transaction ID: %s\n", response.TransactionID)
}
```

## Provider-Specific Usage

### B-PAY Payments

```go
// Create B-PAY specific request
bpayRequest := &rimpay.BPayPaymentRequest{
    Amount:      money.New(decimal.NewFromInt(10000), "MRU"),
    PhoneNumber: phone,
    Reference:   "BPAY-ORDER-123",
    Description: "B-PAY payment",
    Passcode:    "1234", // Required for B-PAY
}

// Process with type-safe method
response, err := client.ProcessBPayPayment(ctx, bpayRequest)
```

### MASRVI Payments

```go
// Create MASRVI specific request
masrviRequest := &rimpay.MasrviPaymentRequest{
    Amount:      money.New(decimal.NewFromInt(50000), "MRU"),
    PhoneNumber: phone,
    Reference:   "MASRVI-ORDER-456",
    Description: "MASRVI payment",
    CallbackURL: "https://yoursite.com/webhook",
    ReturnURL:   "https://yoursite.com/return",
}

// Process with type-safe method
response, err := client.ProcessMasrviPayment(ctx, masrviRequest)
```

## Phone Number Validation

RimPay includes built-in validation for Mauritanian phone numbers:

```go
import "github.com/CatoSystems/rim-pay/pkg/phone"

// Valid Mauritanian phone numbers (prefixes: 2, 3, 4)
validNumbers := []string{
    "+22222334455",  // Mauritantel
    "+22233445566",  // Chinguitel  
    "+22244556677",  // Mattel
}

for _, number := range validNumbers {
    phone, err := phone.NewPhone(number)
    if err != nil {
        log.Printf("Invalid phone: %v", err)
        continue
    }
    fmt.Printf("Valid phone: %s\n", phone.String())
}
```

## Money Handling

The library uses decimal-based precision for accurate money calculations:

```go
import (
    "github.com/CatoSystems/rim-pay/pkg/money"
    "github.com/shopspring/decimal"
)

// Create money amounts
amount1 := money.New(decimal.NewFromInt(10050), "MRU")  // 100.50 MRU
amount2 := money.FromFloat64(75.25, "MRU")              // 75.25 MRU

// Money operations
sum := amount1.Add(amount2)         // 175.75 MRU
diff := amount1.Subtract(amount2)   // 25.25 MRU

fmt.Printf("Amount: %s\n", amount1.String()) // "100.50 MRU"
fmt.Printf("In cents: %d\n", amount1.Cents()) // 10050
```

## Configuration

### Environment Configuration

```go
config := rimpay.DefaultConfig()
config.Environment = rimpay.EnvironmentProduction // or EnvironmentSandbox
config.DefaultProvider = "bpay"

// Configure retry settings
config.Retry.MaxAttempts = 3
config.Retry.InitialDelay = 1 * time.Second
config.Retry.MaxDelay = 10 * time.Second
config.Retry.BackoffMultiplier = 2.0
```

### Provider Configuration

```go
// B-PAY Configuration
config.Providers["bpay"] = rimpay.ProviderConfig{
    Enabled: true,
    BaseURL: "https://api.bpay.mr",
    Credentials: map[string]string{
        "username": "your_bpay_username",
        "password": "your_bpay_password",
    },
    Timeout: 30 * time.Second,
}

// MASRVI Configuration  
config.Providers["masrvi"] = rimpay.ProviderConfig{
    Enabled: true,
    BaseURL: "https://api.masrvi.mr",
    Credentials: map[string]string{
        "merchant_id": "your_merchant_id",
        "api_key":     "your_api_key",
    },
    Timeout: 45 * time.Second,
}
```

```

## Error Handling

RimPay provides comprehensive error handling with detailed error types:

```go
response, err := client.ProcessBPayPayment(ctx, request)
if err != nil {
    switch e := err.(type) {
    case *rimpay.ValidationError:
        fmt.Printf("Validation error: %s (field: %s)\n", e.Message, e.Field)
    case *rimpay.ProviderError:
        fmt.Printf("Provider error: %s (code: %s, provider: %s)\n", 
            e.Message, e.Code, e.Provider)
    case *rimpay.NetworkError:
        fmt.Printf("Network error: %s (retryable: %v)\n", e.Message, e.Retryable)
    default:
        fmt.Printf("Unknown error: %v\n", err)
    }
    return
}
```

## Payment Status Checking

For B-PAY payments, you can check the payment status:

```go
// After successful payment
status, err := client.GetPaymentStatus(ctx, response.TransactionID)
if err != nil {
    log.Printf("Failed to get status: %v", err)
    return
}

fmt.Printf("Payment Status: %s\n", status.Status)
fmt.Printf("Message: %s\n", status.Message)
```

## Multi-Provider Setup

```go
// Configure multiple providers
config := rimpay.DefaultConfig()

// Enable both providers
config.Providers["bpay"] = rimpay.ProviderConfig{
    Enabled: true,
    BaseURL: "https://api.bpay.mr",
    Credentials: map[string]string{
        "username": "bpay_user",
        "password": "bpay_pass",
    },
}

config.Providers["masrvi"] = rimpay.ProviderConfig{
    Enabled: true,
    BaseURL: "https://api.masrvi.mr", 
    Credentials: map[string]string{
        "merchant_id": "masrvi_merchant",
        "api_key":     "masrvi_key",
    },
}

// Use different providers for different requests
client, _ := rimpay.NewClient(config)

// B-PAY payment
bpayResponse, err := client.ProcessBPayPayment(ctx, bpayRequest)

// MASRVI payment  
masrviResponse, err := client.ProcessMasrviPayment(ctx, masrviRequest)
```

## API Reference

### Core Types

#### BPayPaymentRequest
```go
type BPayPaymentRequest struct {
    Amount      *money.Money  // Payment amount
    PhoneNumber *phone.Phone  // Recipient phone number
    Reference   string        // Unique payment reference
    Description string        // Payment description (optional)
    Passcode    string        // B-PAY passcode (required)
}
```

#### MasrviPaymentRequest
```go
type MasrviPaymentRequest struct {
    Amount      *money.Money  // Payment amount
    PhoneNumber *phone.Phone  // Recipient phone number  
    Reference   string        // Unique payment reference
    Description string        // Payment description (optional)
    CallbackURL string        // Webhook URL (optional)
    ReturnURL   string        // Return URL (optional)
}
```

#### PaymentResponse
```go
type PaymentResponse struct {
    TransactionID string     // Unique transaction identifier
    Status        string     // Payment status
    Provider      string     // Payment provider used
    CreatedAt     time.Time  // Transaction creation time
}
```

### Client Methods

#### ProcessBPayPayment
```go
func (c *Client) ProcessBPayPayment(ctx context.Context, request *BPayPaymentRequest) (*PaymentResponse, error)
```

#### ProcessMasrviPayment  
```go
func (c *Client) ProcessMasrviPayment(ctx context.Context, request *MasrviPaymentRequest) (*PaymentResponse, error)
```

#### GetPaymentStatus (B-PAY only)
```go
func (c *Client) GetPaymentStatus(ctx context.Context, transactionID string) (*PaymentStatus, error)
```

## Examples

See the [examples](./examples/) directory for complete working examples:

- [`basic_usage.go`](./examples/basic_usage.go) - Basic payment processing
- [`bpay_example.go`](./examples/bpay_example.go) - B-PAY specific features
- [`masrvi_example.go`](./examples/masrvi_example.go) - MASRVI specific features  
- [`multi_provider_example.go`](./examples/multi_provider_example.go) - Multi-provider usage
- [`error_handling_example.go`](./examples/error_handling_example.go) - Error handling patterns
- [`retry_demo.go`](./examples/retry_demo.go) - Retry configuration
- [`configuration_example.go`](./examples/configuration_example.go) - Advanced configuration

## Testing

Run all tests:
```bash
make test
```

Run tests for specific packages:
```bash
go test ./pkg/...
go test ./internal/...
```

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Development Guidelines

- Write tests for new features
- Follow Go conventions and best practices
- Update documentation for API changes
- Ensure all tests pass before submitting PR

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

For questions, issues, or contributions, please:

1. Check existing [Issues](https://github.com/CatoSystems/rim-pay/issues)
2. Create a new issue with detailed information
3. For security issues, contact us privately

## Changelog

### v0.1.0 (Latest)
- Initial release
- B-PAY and MASRVI provider support
- Provider-specific request types
- Mauritanian phone number validation
- Comprehensive error handling
- Retry mechanisms
- Multi-provider configuration
