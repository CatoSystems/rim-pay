# Installation

## Prerequisites

- Go 1.19 or later
- Git

## Install RimPay

### Using go get (Recommended)

```bash
go get github.com/CatoSystems/rim-pay@latest
```

### Using go mod

Add to your `go.mod` file:

```go
require github.com/CatoSystems/rim-pay v0.1.0
```

Then run:

```bash
go mod tidy
```

### From Source

```bash
git clone https://github.com/CatoSystems/rim-pay.git
cd rim-pay
go install ./...
```

## Verify Installation

Create a simple test file to verify the installation:

```go
// test-install.go
package main

import (
    "fmt"
    "github.com/CatoSystems/rim-pay/pkg/rimpay"
    "github.com/CatoSystems/rim-pay/pkg/phone"
    "github.com/CatoSystems/rim-pay/pkg/money"
    "github.com/shopspring/decimal"
)

func main() {
    // Test phone validation
    phone, err := phone.NewPhone("+22233445566")
    if err != nil {
        panic(err)
    }
    fmt.Printf("✅ Phone: %s\n", phone.String())

    // Test money creation
    amount := money.New(decimal.NewFromInt(10000), money.MRU)
    fmt.Printf("✅ Amount: %s\n", amount.String())

    // Test configuration
    config := rimpay.DefaultConfig()
    fmt.Printf("✅ Default environment: %s\n", config.Environment)

    fmt.Println("🎉 RimPay installation verified!")
}
```

Run the test:

```bash
go run test-install.go
```

Expected output:
```
✅ Phone: +22233445566
✅ Amount: 100.00 MRU
✅ Default environment: sandbox
🎉 RimPay installation verified!
```

## Dependencies

RimPay has minimal dependencies:

- `github.com/shopspring/decimal` - For precise decimal arithmetic
- `github.com/stretchr/testify` - For testing (dev dependency)

## IDE Setup

### VS Code

Install the Go extension:

1. Open VS Code
2. Go to Extensions (Ctrl+Shift+X)
3. Search for "Go" by Google
4. Install the extension

### GoLand/IntelliJ IDEA

The Go plugin should work out of the box with RimPay.

## Next Steps

- [Quick Start Guide](quick-start.md)
- [Configuration](configuration.md)
- [Examples](examples/README.md)
