# API Reference

## Client API

### Client Creation

#### `New(config *Config) (*Client, error)`

Creates a new RimPay client with the specified configuration.

**Parameters:**
- `config`: Configuration object containing provider settings, timeouts, and other options

**Returns:**
- `*Client`: Configured RimPay client
- `error`: Configuration validation error if any

**Example:**
```go
config := rimpay.DefaultConfig()
config.Providers["bpay"] = rimpay.ProviderConfig{
    Enabled: true,
    BaseURL: "https://api.bpay.mr",
    Credentials: map[string]string{
        "username": "your_username",
        "password": "your_password",
    },
}

client, err := rimpay.New(config)
if err != nil {
    log.Fatal(err)
}
```

#### `DefaultConfig() *Config`

Returns a default configuration with sensible defaults.

**Returns:**
- `*Config`: Default configuration object

## Payment Processing

### B-PAY Payments

#### `ProcessBPayPayment(ctx context.Context, request *BPayPaymentRequest) (*PaymentResponse, error)`

Processes a payment through the B-PAY provider.

**Parameters:**
- `ctx`: Context for request cancellation and timeout
- `request`: B-PAY specific payment request

**Returns:**
- `*PaymentResponse`: Payment processing result
- `error`: Processing error if any

**Example:**
```go
phone, _ := phone.Parse("+22233445566")
amount := money.New(decimal.NewFromInt(10000), money.MRU) // 100.00 MRU

request := &rimpay.BPayPaymentRequest{
    Amount:      amount,
    PhoneNumber: phone,
    Reference:   "ORDER-123",
    Description: "Test payment",
    Passcode:    "1234",
}

response, err := client.ProcessBPayPayment(ctx, request)
if err != nil {
    // Handle error
    return
}

fmt.Printf("Transaction ID: %s\n", response.TransactionID)
fmt.Printf("Status: %s\n", response.Status)
```

### MASRVI Payments

#### `ProcessMasrviPayment(ctx context.Context, request *MasrviPaymentRequest) (*PaymentResponse, error)`

Processes a payment through the MASRVI provider.

**Parameters:**
- `ctx`: Context for request cancellation and timeout
- `request`: MASRVI specific payment request

**Returns:**
- `*PaymentResponse`: Payment processing result
- `error`: Processing error if any

**Example:**
```go
phone, _ := phone.Parse("+22233445566")
amount := money.New(decimal.NewFromInt(10000), money.MRU) // 100.00 MRU

request := &rimpay.MasrviPaymentRequest{
    Amount:      amount,
    PhoneNumber: phone,
    Reference:   "ORDER-123",
    Description: "Test payment",
    CallbackURL: "https://yoursite.com/webhook",
    ReturnURL:   "https://yoursite.com/return",
}

response, err := client.ProcessMasrviPayment(ctx, request)
if err != nil {
    // Handle error
    return
}

// For MASRVI, redirect customer to PaymentURL
fmt.Printf("Redirect customer to: %s\n", response.PaymentURL)
```

### Status Checking

#### `GetPaymentStatus(ctx context.Context, transactionID string) (*PaymentStatus, error)`

Retrieves the current status of a payment transaction.

**Parameters:**
- `ctx`: Context for request cancellation and timeout
- `transactionID`: Unique transaction identifier

**Returns:**
- `*PaymentStatus`: Current payment status
- `error`: Status retrieval error if any

**Example:**
```go
status, err := client.GetPaymentStatus(ctx, "TXN123456")
if err != nil {
    // Handle error
    return
}

fmt.Printf("Status: %s\n", status.Status)
fmt.Printf("Amount: %s\n", status.Amount.String())
```

## Request Types

### BPayPaymentRequest

```go
type BPayPaymentRequest struct {
    Amount      money.Money   // Payment amount
    PhoneNumber phone.Number  // Customer phone number
    Reference   string        // Unique payment reference
    Description string        // Payment description
    Passcode    string        // Customer's mobile money PIN
}
```

**Field Validation:**
- `Amount`: Must be positive, supported currency (MRU)
- `PhoneNumber`: Must be valid Mauritanian phone number
- `Reference`: Must be unique, 1-50 characters
- `Description`: Optional, max 255 characters
- `Passcode`: Required, typically 4-6 digits

### MasrviPaymentRequest

```go
type MasrviPaymentRequest struct {
    Amount      money.Money   // Payment amount
    PhoneNumber phone.Number  // Customer phone number
    Reference   string        // Unique payment reference
    Description string        // Payment description
    CallbackURL string        // Webhook URL for status updates
    ReturnURL   string        // URL to redirect customer after payment
}
```

**Field Validation:**
- `Amount`: Must be positive, supported currency (MRU)
- `PhoneNumber`: Must be valid Mauritanian phone number
- `Reference`: Must be unique, 1-50 characters
- `Description`: Optional, max 255 characters
- `CallbackURL`: Must be valid HTTPS URL
- `ReturnURL`: Must be valid HTTPS URL

## Response Types

### PaymentResponse

```go
type PaymentResponse struct {
    TransactionID   string                 // Unique transaction identifier
    Status          PaymentStatus          // Current payment status
    Amount          money.Money           // Payment amount
    PhoneNumber     phone.Number          // Customer phone number
    Reference       string                // Payment reference
    Provider        string                // Payment provider used
    PaymentURL      string                // Redirect URL (MASRVI only)
    CreatedAt       time.Time            // Transaction creation time
    Metadata        map[string]interface{} // Provider-specific metadata
}
```

### PaymentStatus

```go
type PaymentStatus struct {
    Status      StatusType             // Payment status
    Amount      money.Money           // Payment amount
    Reference   string                // Payment reference
    UpdatedAt   time.Time            // Last status update
    Message     string               // Status message
    Metadata    map[string]interface{} // Additional status information
}
```

### StatusType

```go
type StatusType string

const (
    StatusPending   StatusType = "pending"   // Payment initiated
    StatusSuccess   StatusType = "success"   // Payment completed successfully
    StatusFailed    StatusType = "failed"    // Payment failed
    StatusCanceled  StatusType = "canceled"  // Payment canceled
    StatusExpired   StatusType = "expired"   // Payment session expired
)
```

## Configuration Types

### Config

```go
type Config struct {
    Environment    Environment              // Runtime environment
    DefaultTimeout time.Duration           // Default request timeout
    RetryConfig    RetryConfig             // Retry configuration
    Providers      map[string]ProviderConfig // Provider configurations
    Logging        LoggingConfig           // Logging configuration
}
```

### ProviderConfig

```go
type ProviderConfig struct {
    Enabled     bool                   // Whether provider is enabled
    BaseURL     string                // Provider API base URL
    Timeout     time.Duration         // Provider-specific timeout
    Credentials map[string]string     // Provider credentials
    Options     map[string]interface{} // Provider-specific options
}
```

### RetryConfig

```go
type RetryConfig struct {
    MaxAttempts     int           // Maximum retry attempts
    InitialDelay    time.Duration // Initial retry delay
    MaxDelay        time.Duration // Maximum retry delay
    BackoffFactor   float64       // Exponential backoff factor
    RetryableErrors []ErrorType   // Error types that trigger retry
}
```

### LoggingConfig

```go
type LoggingConfig struct {
    Level  LogLevel // Minimum log level
    Format string   // Log format (json, text)
    Output string   // Log output (stdout, stderr, file)
}
```

## Error Types

### ValidationError

```go
type ValidationError struct {
    Field   string // Field that failed validation
    Value   string // Invalid value
    Message string // Validation error message
}

func (e *ValidationError) Error() string
```

### ProviderError

```go
type ProviderError struct {
    Provider string // Provider name
    Code     string // Provider-specific error code
    Message  string // Error message
    Details  map[string]interface{} // Additional error details
}

func (e *ProviderError) Error() string
```

### NetworkError

```go
type NetworkError struct {
    Operation string // Network operation that failed
    URL       string // Target URL
    Message   string // Error message
    Cause     error  // Underlying error
}

func (e *NetworkError) Error() string
```

### AuthenticationError

```go
type AuthenticationError struct {
    Provider string // Provider name
    Message  string // Authentication error message
}

func (e *AuthenticationError) Error() string
```

## Utility Functions

### Error Type Checking

```go
// IsValidationError checks if error is a validation error
func IsValidationError(err error) bool

// IsProviderError checks if error is a provider error
func IsProviderError(err error) bool

// IsNetworkError checks if error is a network error
func IsNetworkError(err error) bool

// IsAuthenticationError checks if error is an authentication error
func IsAuthenticationError(err error) bool
```

### Status Helpers

```go
// IsPaymentSuccessful checks if payment completed successfully
func IsPaymentSuccessful(status StatusType) bool

// IsPaymentFinal checks if payment is in a final state
func IsPaymentFinal(status StatusType) bool

// IsPaymentPending checks if payment is still processing
func IsPaymentPending(status StatusType) bool
```

## Constants

### Environment Types

```go
type Environment string

const (
    EnvironmentDevelopment Environment = "development"
    EnvironmentSandbox     Environment = "sandbox"
    EnvironmentProduction  Environment = "production"
)
```

### Error Types

```go
type ErrorType string

const (
    ErrorTypeValidation     ErrorType = "validation"
    ErrorTypeProvider      ErrorType = "provider"
    ErrorTypeNetwork       ErrorType = "network"
    ErrorTypeAuthentication ErrorType = "authentication"
    ErrorTypeConfiguration ErrorType = "configuration"
)
```

### Log Levels

```go
type LogLevel string

const (
    LogLevelDebug LogLevel = "debug"
    LogLevelInfo  LogLevel = "info"
    LogLevelWarn  LogLevel = "warn"
    LogLevelError LogLevel = "error"
)
```

## Phone Package API

### phone.Parse(s string) (Number, error)

Parses a phone number string into a structured phone number.

**Parameters:**
- `s`: Phone number string in various formats

**Returns:**
- `Number`: Parsed and validated phone number
- `error`: Parsing error if invalid

**Example:**
```go
phone, err := phone.Parse("+22233445566")
if err != nil {
    // Handle parsing error
}

fmt.Printf("Country: %s\n", phone.Country())
fmt.Printf("National: %s\n", phone.National())
fmt.Printf("International: %s\n", phone.International())
```

### Number Methods

```go
type Number interface {
    // String returns the international format
    String() string
    
    // International returns international format (+22233445566)
    International() string
    
    // National returns national format (33445566)
    National() string
    
    // Country returns country code (+222)
    Country() string
    
    // IsValid returns true if number is valid
    IsValid() bool
    
    // Operator returns mobile operator
    Operator() string
}
```

## Money Package API

### money.New(amount decimal.Decimal, currency Currency) Money

Creates a new money value with the specified amount and currency.

**Parameters:**
- `amount`: Decimal amount
- `currency`: Currency code

**Returns:**
- `Money`: Money value

**Example:**
```go
amount := money.New(decimal.NewFromInt(10050), money.MRU) // 100.50 MRU
fmt.Printf("Amount: %s\n", amount.String()) // "100.50 MRU"
```

### money.FromFloat64(amount float64, currency Currency) Money

Creates a new money value from a float64 (less precise than decimal).

**Parameters:**
- `amount`: Float64 amount
- `currency`: Currency code

**Returns:**
- `Money`: Money value

### Money Methods

```go
type Money interface {
    // Amount returns the decimal amount
    Amount() decimal.Decimal
    
    // Currency returns the currency
    Currency() Currency
    
    // String returns formatted amount with currency
    String() string
    
    // Add adds another money value
    Add(other Money) (Money, error)
    
    // Subtract subtracts another money value
    Subtract(other Money) (Money, error)
    
    // Multiply multiplies by a decimal factor
    Multiply(factor decimal.Decimal) Money
    
    // IsZero returns true if amount is zero
    IsZero() bool
    
    // IsPositive returns true if amount is positive
    IsPositive() bool
    
    // IsNegative returns true if amount is negative
    IsNegative() bool
}
```

### Currency Constants

```go
type Currency string

const (
    MRU Currency = "MRU" // Mauritanian Ouguiya
)
```

This API reference provides comprehensive documentation for all public interfaces, types, and functions in the RimPay library.
