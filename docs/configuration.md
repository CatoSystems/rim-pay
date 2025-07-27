# Configuration

RimPay provides a flexible configuration system that allows you to customize behavior for different environments and providers.

## Basic Configuration

### Default Configuration

```go
config := rimpay.DefaultConfig()
// Returns a configuration with sensible defaults:
// - Environment: sandbox
// - DefaultProvider: "" (must be set)
// - Retry settings with exponential backoff
// - Empty providers map (must be configured)
```

### Custom Configuration

```go
config := &rimpay.Config{
    Environment:     rimpay.EnvironmentProduction,
    DefaultProvider: "bpay",
    Providers:       make(map[string]rimpay.ProviderConfig),
    Retry: rimpay.RetryConfig{
        MaxAttempts:        3,
        InitialDelay:       1 * time.Second,
        MaxDelay:          10 * time.Second,
        BackoffMultiplier: 2.0,
    },
}
```

## Environments

RimPay supports two environments:

### Sandbox Environment
```go
config.Environment = rimpay.EnvironmentSandbox
```
- Use for development and testing
- No real money transactions
- Test credentials and endpoints

### Production Environment
```go
config.Environment = rimpay.EnvironmentProduction
```
- Use for live transactions
- Real money processing
- Production credentials required

## Provider Configuration

### B-PAY Provider

```go
config.Providers["bpay"] = rimpay.ProviderConfig{
    Enabled: true,
    BaseURL: "https://api.bpay.mr", // Production
    // BaseURL: "https://sandbox.bpay.mr", // Sandbox
    Credentials: map[string]string{
        "username":  "your_bpay_username",
        "password":  "your_bpay_password",
        "client_id": "your_client_id", // Optional
    },
    Timeout: 30 * time.Second,
}
```

### MASRVI Provider

```go
config.Providers["masrvi"] = rimpay.ProviderConfig{
    Enabled: true,
    BaseURL: "https://api.masrvi.mr", // Production
    // BaseURL: "https://sandbox.masrvi.mr", // Sandbox
    Credentials: map[string]string{
        "merchant_id": "your_merchant_id",
        "api_key":     "your_api_key",
        "secret_key":  "your_secret_key", // Optional
    },
    Timeout: 45 * time.Second,
}
```

## Retry Configuration

RimPay includes built-in retry mechanisms for handling transient failures:

```go
config.Retry = rimpay.RetryConfig{
    MaxAttempts:        5,                // Maximum retry attempts
    InitialDelay:       1 * time.Second,  // Initial delay between retries
    MaxDelay:          30 * time.Second,  // Maximum delay between retries
    BackoffMultiplier: 2.0,               // Exponential backoff multiplier
}
```

### Retry Behavior

1. **Initial attempt**: No delay
2. **First retry**: 1 second delay
3. **Second retry**: 2 seconds delay (1 × 2.0)
4. **Third retry**: 4 seconds delay (2 × 2.0)
5. **Fourth retry**: 8 seconds delay (4 × 2.0)
6. **Fifth retry**: 16 seconds delay (capped at MaxDelay)

### Retryable Conditions

RimPay automatically retries on:
- Network timeouts
- Connection errors
- HTTP 5xx server errors
- Authentication token expiration (with re-authentication)

## Multi-Provider Setup

You can configure multiple providers and switch between them:

```go
config := rimpay.DefaultConfig()
config.Environment = rimpay.EnvironmentSandbox

// Configure B-PAY
config.Providers["bpay"] = rimpay.ProviderConfig{
    Enabled: true,
    BaseURL: "https://sandbox.bpay.mr",
    Credentials: map[string]string{
        "username": "bpay_test_user",
        "password": "bpay_test_pass",
    },
    Timeout: 30 * time.Second,
}

// Configure MASRVI
config.Providers["masrvi"] = rimpay.ProviderConfig{
    Enabled: true,
    BaseURL: "https://sandbox.masrvi.mr",
    Credentials: map[string]string{
        "merchant_id": "masrvi_test_merchant",
        "api_key":     "masrvi_test_key",
    },
    Timeout: 45 * time.Second,
}

// Set default provider
config.DefaultProvider = "bpay"

client, err := rimpay.NewClient(config)
```

## Environment Variables

You can use environment variables for sensitive configuration:

```go
import "os"

config := rimpay.DefaultConfig()
config.Providers["bpay"] = rimpay.ProviderConfig{
    Enabled: true,
    BaseURL: os.Getenv("BPAY_BASE_URL"),
    Credentials: map[string]string{
        "username": os.Getenv("BPAY_USERNAME"),
        "password": os.Getenv("BPAY_PASSWORD"),
    },
    Timeout: 30 * time.Second,
}
```

### Recommended Environment Variables

```bash
# B-PAY Configuration
export BPAY_BASE_URL="https://api.bpay.mr"
export BPAY_USERNAME="your_username"
export BPAY_PASSWORD="your_password"
export BPAY_CLIENT_ID="your_client_id"

# MASRVI Configuration
export MASRVI_BASE_URL="https://api.masrvi.mr"
export MASRVI_MERCHANT_ID="your_merchant_id"
export MASRVI_API_KEY="your_api_key"
export MASRVI_SECRET_KEY="your_secret_key"

# General Configuration
export RIMPAY_ENVIRONMENT="production"
export RIMPAY_DEFAULT_PROVIDER="bpay"
```

## Configuration Validation

RimPay validates configuration at client creation:

```go
client, err := rimpay.NewClient(config)
if err != nil {
    // Handle configuration errors
    switch e := err.(type) {
    case *rimpay.ConfigError:
        fmt.Printf("Configuration error: %s\n", e.Message)
    default:
        fmt.Printf("Unknown error: %v\n", err)
    }
}
```

### Common Configuration Errors

- **Missing provider configuration**: No providers configured
- **Invalid environment**: Environment must be "sandbox" or "production"
- **Missing credentials**: Required credentials not provided
- **Invalid timeout**: Timeout must be positive
- **Provider not found**: DefaultProvider not in Providers map

## Security Best Practices

1. **Never hardcode credentials** in source code
2. **Use environment variables** for sensitive data
3. **Rotate credentials regularly**
4. **Use different credentials** for sandbox and production
5. **Monitor credential usage** and access logs
6. **Implement proper secret management** in production

## Example: Production Configuration

```go
package main

import (
    "os"
    "time"
    "github.com/CatoSystems/rim-pay/pkg/rimpay"
)

func createProductionConfig() *rimpay.Config {
    config := &rimpay.Config{
        Environment:     rimpay.EnvironmentProduction,
        DefaultProvider: os.Getenv("RIMPAY_DEFAULT_PROVIDER"),
        Providers:       make(map[string]rimpay.ProviderConfig),
        Retry: rimpay.RetryConfig{
            MaxAttempts:        3,
            InitialDelay:       2 * time.Second,
            MaxDelay:          30 * time.Second,
            BackoffMultiplier: 2.0,
        },
    }

    // B-PAY Configuration
    if os.Getenv("BPAY_USERNAME") != "" {
        config.Providers["bpay"] = rimpay.ProviderConfig{
            Enabled: true,
            BaseURL: os.Getenv("BPAY_BASE_URL"),
            Credentials: map[string]string{
                "username":  os.Getenv("BPAY_USERNAME"),
                "password":  os.Getenv("BPAY_PASSWORD"),
                "client_id": os.Getenv("BPAY_CLIENT_ID"),
            },
            Timeout: 30 * time.Second,
        }
    }

    // MASRVI Configuration
    if os.Getenv("MASRVI_MERCHANT_ID") != "" {
        config.Providers["masrvi"] = rimpay.ProviderConfig{
            Enabled: true,
            BaseURL: os.Getenv("MASRVI_BASE_URL"),
            Credentials: map[string]string{
                "merchant_id": os.Getenv("MASRVI_MERCHANT_ID"),
                "api_key":     os.Getenv("MASRVI_API_KEY"),
                "secret_key":  os.Getenv("MASRVI_SECRET_KEY"),
            },
            Timeout: 45 * time.Second,
        }
    }

    return config
}
```

## Next Steps

- Learn about [specific payment providers](providers/README.md)
- Understand [error handling](error-handling.md)
- Explore [retry mechanisms](retry-mechanisms.md)
