# RimPay Architecture

This document describes the architecture and design principles of the RimPay payment processing library.

## Overview

RimPay is designed as a modular, extensible payment processing library that provides a unified interface for multiple payment providers while maintaining type safety, reliability, and ease of use.

## High-Level Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    Client Application                       │
└─────────────────────┬───────────────────────────────────────┘
                      │
┌─────────────────────▼───────────────────────────────────────┐
│                 RimPay Client                               │
│  ┌─────────────────┬─────────────────┬─────────────────────┐│
│  │  Configuration  │   Validation    │   Error Handling    ││
│  └─────────────────┴─────────────────┴─────────────────────┘│
└─────────────────────┬───────────────────────────────────────┘
                      │
┌─────────────────────▼───────────────────────────────────────┐
│                Provider Abstraction                         │
│  ┌─────────────────┬─────────────────┬─────────────────────┐│
│  │  Provider A     │   Provider B    │   Provider C        ││
│  │  (B-PAY)        │   (MASRVI)      │   (Future)          ││
│  └─────────────────┴─────────────────┴─────────────────────┘│
└─────────────────────┬───────────────────────────────────────┘
                      │
┌─────────────────────▼───────────────────────────────────────┐
│              Supporting Packages                            │
│  ┌─────────────────┬─────────────────┬─────────────────────┐│
│  │   Phone Pkg     │   Money Pkg     │   Common Utils      ││
│  │  (Validation)   │ (Decimal Math)  │ (HTTP, Retry, etc)  ││
│  └─────────────────┴─────────────────┴─────────────────────┘│
└─────────────────────────────────────────────────────────────┘
```

## Core Components

### 1. Client Layer (`pkg/rimpay`)

The client layer provides the main interface for applications to interact with RimPay.

```go
type Client struct {
    config    *Config
    providers map[string]PaymentProvider
    logger    Logger
}
```

**Key Responsibilities:**
- Request validation and preprocessing
- Provider selection and routing
- Response handling and normalization
- Error handling and retry logic
- Configuration management

### 2. Provider Abstraction (`internal/providers`)

The provider layer abstracts different payment providers behind a common interface:

```go
type PaymentProvider interface {
    ProcessPayment(ctx context.Context, request PaymentRequest) (*PaymentResponse, error)
    GetPaymentStatus(ctx context.Context, transactionID string) (*PaymentStatus, error)
    GetCapabilities() ProviderCapabilities
}
```

**Design Benefits:**
- Uniform API across providers
- Easy addition of new providers
- Provider-specific optimizations
- Isolated provider logic

### 3. Request/Response Types

Provider-specific request types ensure type safety while maintaining flexibility:

```go
// Provider-specific requests
type BPayPaymentRequest struct {
    Amount      money.Money
    PhoneNumber phone.Number
    Reference   string
    Description string
    Passcode    string  // B-PAY specific
}

type MasrviPaymentRequest struct {
    Amount      money.Money
    PhoneNumber phone.Number
    Reference   string
    Description string
    CallbackURL string  // MASRVI specific
    ReturnURL   string  // MASRVI specific
}
```

### 4. Supporting Packages

#### Phone Package (`pkg/phone`)
- Mauritanian phone number validation
- Format normalization
- Operator detection

#### Money Package (`pkg/money`)
- Decimal-based arithmetic
- Currency support (MRU)
- Precision handling

#### Common Utilities (`internal/providers/common`)
- HTTP client with retries
- Rate limiting
- Logging utilities

## Provider Implementations

### B-PAY Provider (`internal/providers/bpay`)

B-PAY follows a session-based authentication model:

```
┌─────────┐    ┌─────────┐    ┌─────────┐    ┌─────────┐
│ Client  │    │ RimPay  │    │ B-PAY   │    │ Mobile  │
│         │    │         │    │ API     │    │ Network │
└────┬────┘    └────┬────┘    └────┬────┘    └────┬────┘
     │              │              │              │
     │ PaymentReq   │              │              │
     │─────────────▶│              │              │
     │              │ Auth Request │              │
     │              │─────────────▶│              │
     │              │ Auth Token   │              │
     │              │◀─────────────│              │
     │              │ Payment Req  │              │
     │              │─────────────▶│              │
     │              │              │ Mobile Req   │
     │              │              │─────────────▶│
     │              │              │ Mobile Resp  │
     │              │              │◀─────────────│
     │              │ Payment Resp │              │
     │              │◀─────────────│              │
     │ PaymentResp  │              │              │
     │◀─────────────│              │              │
```

**Key Features:**
- Username/password authentication
- Real-time payment processing
- Direct mobile network integration
- Status checking support

### MASRVI Provider (`internal/providers/masrvi`)

MASRVI uses a web-based redirect flow:

```
┌─────────┐    ┌─────────┐    ┌─────────┐    ┌─────────┐
│ Client  │    │ RimPay  │    │ MASRVI  │    │ Browser │
│         │    │         │    │ API     │    │         │
└────┬────┘    └────┬────┘    └────┬────┘    └────┬────┘
     │              │              │              │
     │ PaymentReq   │              │              │
     │─────────────▶│              │              │
     │              │ Session Req  │              │
     │              │─────────────▶│              │
     │              │ Session Resp │              │
     │              │◀─────────────│              │
     │ PaymentResp  │              │              │
     │◀─────────────│              │              │
     │              │              │              │
     │ Redirect     │              │              │
     │─────────────────────────────────────────▶│
     │              │              │ Payment UI   │
     │              │              │◀─────────────│
     │              │              │ Completion   │
     │              │              │─────────────▶│
     │              │ Webhook      │              │
     │              │◀─────────────│              │
```

**Key Features:**
- Merchant ID/API key authentication
- Session-based payments
- Web redirect flow
- Webhook notifications

## Error Handling Strategy

RimPay implements a comprehensive error handling strategy with specific error types:

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

### Error Hierarchy

```
Error (interface)
├── ValidationError
│   ├── PhoneValidationError
│   ├── AmountValidationError
│   └── FieldValidationError
├── ProviderError
│   ├── BPayError
│   ├── MasrviError
│   └── GenericProviderError
├── NetworkError
│   ├── TimeoutError
│   ├── ConnectionError
│   └── HTTPError
└── ConfigurationError
    ├── MissingCredentialsError
    ├── InvalidConfigError
    └── ProviderNotConfiguredError
```

### Retry Strategy

RimPay implements intelligent retry logic:

1. **Immediate Retry**: For transient network errors
2. **Exponential Backoff**: For rate limiting and server errors
3. **No Retry**: For validation and business logic errors

```go
type RetryConfig struct {
    MaxAttempts     int
    InitialDelay    time.Duration
    MaxDelay        time.Duration
    BackoffFactor   float64
    RetryableErrors []ErrorType
}
```

## Configuration Management

### Configuration Structure

```go
type Config struct {
    Environment    Environment
    DefaultTimeout time.Duration
    RetryConfig    RetryConfig
    Providers      map[string]ProviderConfig
    Logging        LoggingConfig
}

type ProviderConfig struct {
    Enabled     bool
    BaseURL     string
    Timeout     time.Duration
    Credentials map[string]string
    Options     map[string]interface{}
}
```

### Environment Support

- **Development**: Relaxed validation, verbose logging
- **Sandbox**: Provider test environments
- **Production**: Strict validation, minimal logging

## Security Considerations

### Credential Management

1. **Separation**: Production and sandbox credentials are separate
2. **Encryption**: Credentials encrypted in memory when possible
3. **Rotation**: Support for credential rotation without downtime
4. **Audit**: All credential usage is logged

### Data Protection

1. **PII Handling**: Phone numbers are validated but not logged
2. **Financial Data**: Amounts are handled with decimal precision
3. **Sensitive Data**: PINs and passcodes are cleared from memory
4. **Transport Security**: All communications use HTTPS

### Rate Limiting

```go
type RateLimiter struct {
    requests    int
    window      time.Duration
    resetTime   time.Time
    mutex       sync.Mutex
}
```

## Performance Characteristics

### Throughput

- **B-PAY**: ~50 requests/second per merchant
- **MASRVI**: ~100 requests/second per merchant
- **Concurrent**: Supports concurrent requests with proper rate limiting

### Latency

- **B-PAY**: 2-10 seconds (mobile network dependent)
- **MASRVI**: 1-5 seconds (session creation)
- **Validation**: <1ms (local validation)

### Memory Usage

- **Base**: ~10MB for library initialization
- **Per Request**: ~1KB additional memory
- **Connection Pooling**: Reduces overhead for multiple requests

## Extensibility

### Adding New Providers

1. Implement the `PaymentProvider` interface
2. Create provider-specific request/response types
3. Add provider configuration
4. Implement error mapping
5. Add comprehensive tests

```go
type NewProvider struct {
    config     ProviderConfig
    httpClient *http.Client
    logger     Logger
}

func (p *NewProvider) ProcessPayment(ctx context.Context, request PaymentRequest) (*PaymentResponse, error) {
    // Implementation
}
```

### Custom Validation Rules

Extend the validation system with custom rules:

```go
func RegisterCustomValidator(field string, validator func(interface{}) error) {
    // Register custom validation logic
}
```

### Middleware Support

Add custom middleware for logging, metrics, etc.:

```go
type Middleware func(HandlerFunc) HandlerFunc

func (c *Client) Use(middleware Middleware) {
    // Add middleware to processing chain
}
```

## Testing Strategy

### Unit Tests

- **Provider Tests**: Mock HTTP responses
- **Validation Tests**: Comprehensive input validation
- **Error Handling Tests**: All error scenarios
- **Configuration Tests**: Various config combinations

### Integration Tests

- **Sandbox Testing**: Real provider sandbox APIs
- **End-to-End**: Complete payment flows
- **Error Recovery**: Network failure scenarios
- **Performance**: Load and stress testing

### Test Coverage

Current test coverage targets:
- **Unit Tests**: >90% code coverage
- **Integration Tests**: All happy paths and major error cases
- **Performance Tests**: Baseline performance metrics

## Monitoring and Observability

### Logging

```go
type Logger interface {
    Debug(msg string, fields ...Field)
    Info(msg string, fields ...Field)
    Warn(msg string, fields ...Field)
    Error(msg string, fields ...Field)
}
```

### Metrics

Key metrics to monitor:
- Request count by provider
- Response time percentiles
- Error rate by type
- Provider availability
- Configuration changes

### Tracing

Support for distributed tracing:
- Request correlation IDs
- Provider call tracing
- Error propagation tracking

## Future Enhancements

### Short Term (Next Release)

1. **Circuit Breaker**: Automatic provider failover
2. **Metrics Export**: Prometheus/StatsD support
3. **Connection Pooling**: Improved HTTP performance
4. **Caching**: Response caching for status checks

### Medium Term

1. **Additional Providers**: More Mauritanian payment providers
2. **Async Processing**: Background payment processing
3. **Batch Operations**: Multiple payments in single request
4. **Admin API**: Configuration management API

### Long Term

1. **Multi-Region**: Support for other African countries
2. **Event Sourcing**: Complete audit trail
3. **GraphQL API**: Alternative API interface
4. **Mobile SDK**: Direct mobile integration

## Design Principles

### 1. Type Safety

All operations are type-safe with minimal runtime type assertions:
- Provider-specific request types
- Compile-time interface verification
- Generic error handling

### 2. Fail Fast

Invalid configurations and requests fail immediately:
- Validation at client creation
- Request validation before provider calls
- Clear error messages

### 3. Extensibility

Easy to extend without breaking existing code:
- Interface-based design
- Plugin architecture for providers
- Middleware support

### 4. Performance

Optimized for high-throughput environments:
- Connection pooling
- Efficient retry mechanisms
- Minimal memory allocations

### 5. Reliability

Built for production use:
- Comprehensive error handling
- Automatic retries
- Circuit breaker patterns
- Graceful degradation

This architecture ensures RimPay is maintainable, extensible, and production-ready while providing a simple and consistent API for developers.
