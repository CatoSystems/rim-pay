# Quick Start Guide

This guide will help you process your first payment with RimPay in just a few minutes.

## Step 1: Import RimPay

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
```

## Step 2: Configure RimPay

```go
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
```

## Step 3: Create a Payment Request

```go
    // Create phone number (Mauritanian format)
    phone, err := phone.NewPhone("+22233445566")
    if err != nil {
        log.Fatal(err)
    }
    
    // Create amount (100.00 MRU)
    amount := money.New(decimal.NewFromInt(10000), money.MRU)
    
    // Create B-PAY payment request
    request := &rimpay.BPayPaymentRequest{
        Amount:      amount,
        PhoneNumber: phone,
        Reference:   "ORDER-12345",
        Description: "Payment for order 12345",
        Passcode:    "1234", // Customer's mobile money PIN
    }
```

## Step 4: Process the Payment

```go
    // Process payment
    ctx := context.Background()
    response, err := client.ProcessBPayPayment(ctx, request)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Payment successful!\n")
    fmt.Printf("Transaction ID: %s\n", response.TransactionID)
    fmt.Printf("Status: %s\n", response.Status)
}
```

## Complete Example

Here's the complete working example:

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
    
    // Create phone number
    phone, err := phone.NewPhone("+22233445566")
    if err != nil {
        log.Fatal(err)
    }
    
    // Create amount
    amount := money.New(decimal.NewFromInt(10000), money.MRU)
    
    // Create payment request
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
    
    fmt.Printf("Payment successful!\n")
    fmt.Printf("Transaction ID: %s\n", response.TransactionID)
    fmt.Printf("Status: %s\n", response.Status)
}
```

## Running the Example

Save the code as `main.go` and run:

```bash
go mod init my-payment-app
go get github.com/CatoSystems/rim-pay@latest
go run main.go
```

## What's Next?

- Learn about [different payment providers](providers/README.md)
- Explore [advanced configuration options](configuration.md)
- Handle [errors gracefully](error-handling.md)
- Check out [more examples](examples/README.md)

## Common Issues

### Invalid Phone Number
Make sure your phone numbers use Mauritanian prefixes (2, 3, or 4):
- ✅ `+22233445566` (Mattel)
- ✅ `+22222334455` (Chinguitel)  
- ✅ `+22244556677` (Mauritel)
- ❌ `+22255667788` (Invalid prefix 5)

### Authentication Errors
Ensure your provider credentials are correct and your account is active in the sandbox/production environment you're targeting.

### Network Errors
RimPay includes automatic retry mechanisms, but ensure your network connectivity is stable.
