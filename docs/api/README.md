# API Reference

This section provides detailed API documentation for all RimPay components.

## Core Components

### [Client API](client.md)
The main client interface for processing payments and managing providers.

### [Request Types](requests.md)
Provider-specific payment request structures and validation.

### [Response Types](responses.md)
Payment response structures and status information.

### [Error Types](errors.md)
Comprehensive error handling with specific error types.

## Package Structure

```
github.com/CatoSystems/rim-pay/
├── pkg/
│   ├── rimpay/          # Main client and types
│   ├── phone/           # Phone number validation
│   └── money/           # Money handling
└── internal/
    └── providers/       # Provider implementations
        ├── bpay/        # B-PAY provider
        ├── masrvi/      # MASRVI provider
        └── common/      # Shared utilities
```

## Quick Reference

### Client Creation

```go
// Default configuration
config := rimpay.DefaultConfig()
client, err := rimpay.NewClient(config)

// Custom configuration
config := &rimpay.Config{
    Environment:     rimpay.EnvironmentProduction,
    DefaultProvider: "bpay",
    Providers:       map[string]rimpay.ProviderConfig{...},
}
client, err := rimpay.NewClient(config)
```

### Payment Processing

```go
// B-PAY payment
bpayRequest := &rimpay.BPayPaymentRequest{...}
response, err := client.ProcessBPayPayment(ctx, bpayRequest)

// MASRVI payment
masrviRequest := &rimpay.MasrviPaymentRequest{...}
response, err := client.ProcessMasrviPayment(ctx, masrviRequest)
```

### Status Checking

```go
// B-PAY status (supported)
status, err := client.GetPaymentStatus(ctx, transactionID)

// MASRVI uses webhooks for status updates
```

### Phone Number Validation

```go
import "github.com/CatoSystems/rim-pay/pkg/phone"

phone, err := phone.NewPhone("+22233445566")
fmt.Println(phone.String())           // "+22233445566"
fmt.Println(phone.ForProvider(false)) // "33445566"
```

### Money Handling

```go
import "github.com/CatoSystems/rim-pay/pkg/money"

amount := money.New(decimal.NewFromInt(10000), money.MRU)
fmt.Println(amount.String())       // "100.00 MRU"
fmt.Println(amount.Cents())        // 10000
```

## Type Definitions

### Core Types

```go
// Configuration
type Config struct {
    Environment     Environment
    DefaultProvider string
    Providers       map[string]ProviderConfig
    Retry          RetryConfig
}

// Provider configuration
type ProviderConfig struct {
    Enabled     bool
    BaseURL     string
    Credentials map[string]string
    Timeout     time.Duration
}

// Payment response
type PaymentResponse struct {
    TransactionID string
    Status        PaymentStatus
    Provider      string
    CreatedAt     time.Time
}
```

### Request Types

```go
// B-PAY request
type BPayPaymentRequest struct {
    Amount      *money.Money
    PhoneNumber *phone.Phone
    Reference   string
    Description string
    Passcode    string
}

// MASRVI request
type MasrviPaymentRequest struct {
    Amount      *money.Money
    PhoneNumber *phone.Phone
    Reference   string
    Description string
    CallbackURL string
    ReturnURL   string
}
```

### Error Types

```go
type ValidationError struct {
    Field   string
    Message string
    Value   string
}

type ProviderError struct {
    Provider string
    Code     string
    Message  string
}

type NetworkError struct {
    Message   string
    Retryable bool
    Cause     error
}

type AuthenticationError struct {
    Provider string
    Message  string
}
```

## Constants

### Environments

```go
const (
    EnvironmentSandbox    Environment = "sandbox"
    EnvironmentProduction Environment = "production"
)
```

### Payment Status

```go
const (
    PaymentStatusPending PaymentStatus = "pending"
    PaymentStatusSuccess PaymentStatus = "success"
    PaymentStatusFailed  PaymentStatus = "failed"
)
```

### Error Codes

```go
const (
    ErrorCodeValidationFailed      = "VALIDATION_FAILED"
    ErrorCodeAuthenticationFailed  = "AUTHENTICATION_FAILED"
    ErrorCodeInsufficientFunds     = "INSUFFICIENT_FUNDS"
    ErrorCodePaymentDeclined       = "PAYMENT_DECLINED"
    ErrorCodeNetworkError          = "NETWORK_ERROR"
    ErrorCodeProviderError         = "PROVIDER_ERROR"
)
```

## Interface Definitions

### Payment Provider Interface

```go
type PaymentProvider interface {
    ProcessPayment(ctx context.Context, request PaymentRequest) (*PaymentResponse, error)
    GetPaymentStatus(ctx context.Context, transactionID string) (*PaymentStatus, error)
    ValidateCredentials(ctx context.Context) error
}
```

### Payment Request Interface

```go
type PaymentRequest interface {
    Validate() error
    GetAmount() *money.Money
    GetPhoneNumber() *phone.Phone
    GetReference() string
    ToGenericRequest() *GenericPaymentRequest
}
```

## Version Information

- **Current Version**: v0.1.0
- **Go Version**: 1.19+
- **API Stability**: Beta

## Compatibility

### Go Versions
- Minimum: Go 1.19
- Tested: Go 1.19, 1.20, 1.21, 1.22
- Recommended: Go 1.21+

### Dependencies
- `github.com/shopspring/decimal` v1.4.0+
- `github.com/stretchr/testify` v1.10.0+ (testing only)

## Migration Guide

### From v0.0.x to v0.1.0

1. **Update import paths**:
   ```go
   // Old
   import "github.com/CatoSystems/rim-pay/rimpay"
   
   // New
   import "github.com/CatoSystems/rim-pay/pkg/rimpay"
   ```

2. **Use provider-specific requests**:
   ```go
   // Old
   request := &rimpay.PaymentRequest{...}
   response, err := client.ProcessPayment(ctx, request)
   
   // New
   bpayRequest := &rimpay.BPayPaymentRequest{...}
   response, err := client.ProcessBPayPayment(ctx, bpayRequest)
   ```

3. **Update phone number handling**:
   ```go
   // Old
   request.PhoneNumber = "+22233445566"
   
   // New
   phone, err := phone.NewPhone("+22233445566")
   request.PhoneNumber = phone
   ```

## Next Steps

- [Client API Reference](client.md)
- [Request Types Reference](requests.md)
- [Response Types Reference](responses.md)
- [Error Types Reference](errors.md)
