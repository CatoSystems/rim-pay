# RimPay Documentation

Welcome to the comprehensive documentation for RimPay, a Go library for payment processing in Mauritania.

## Table of Contents

### Getting Started
- [Installation Guide](installation.md) - How to install and set up RimPay
- [Quick Start Guide](quick-start.md) - Get up and running in minutes
- [Configuration](configuration.md) - Complete configuration options

### Core Concepts
- [Phone Number Validation](phone-validation.md) - Mauritanian phone number handling
- [Error Handling](error-handling.md) - Comprehensive error management
- [Architecture Overview](architecture.md) - System design and architecture

### Payment Providers
- [Provider Overview](providers/README.md) - Supported payment providers and comparison

### API Reference
- [API Documentation](api/README.md) - Complete API reference overview
- [Complete Reference](api/reference.md) - Full API reference

### Examples
- [Code Examples](examples/README.md) - Practical usage examples and patterns

### Advanced Topics
- [Testing Guide](testing.md) - Testing strategies and utilities
- [FAQ](faq.md) - Frequently asked questions

## Quick Reference

### Supported Providers
| Provider | Type | Features | Status |
|----------|------|----------|---------|
| B-PAY | Mobile Money | Real-time, PIN-based | ✅ Active |
| MASRVI | Web Payment | Session-based, Redirect | ✅ Active |

### Supported Countries
- 🇲🇷 **Mauritania** (Primary focus)

### Currency Support
- **MRU** (Mauritanian Ouguiya) - Current currency since 2018

## Quick Start

```go
import "github.com/CatoSystems/rim-pay/pkg/rimpay"

// Create client
config := rimpay.DefaultConfig()
client, err := rimpay.New(config)

// Process payment
phone, _ := phone.Parse("+22233445566")
amount := money.New(decimal.NewFromInt(10000), money.MRU)

request := &rimpay.BPayPaymentRequest{
    Amount:      amount,
    PhoneNumber: phone,
    Reference:   "ORDER-123",
    Passcode:    "1234",
}

response, err := client.ProcessBPayPayment(ctx, request)
```

## Documentation Structure

```
docs/
├── README.md              # This file - main documentation index
├── installation.md        # Installation and setup
├── quick-start.md         # Getting started tutorial
├── configuration.md       # Configuration options
├── phone-validation.md    # Phone number validation
├── error-handling.md      # Error handling patterns
├── architecture.md        # System architecture
├── testing.md            # Testing guide
├── faq.md                # Frequently asked questions
├── providers/            # Provider-specific documentation
│   └── README.md         # Provider overview
├── api/                 # API reference documentation
│   ├── README.md        # API overview
│   └── reference.md     # Complete API reference
└── examples/            # Example documentation
    └── README.md        # Examples overview
```

## Contributing to Documentation

We welcome contributions to improve our documentation:

1. **Fix typos or errors** - Submit a pull request with corrections
2. **Add examples** - Share real-world usage patterns
3. **Improve explanations** - Make complex topics clearer
4. **Update outdated content** - Keep documentation current

See our [Contributing Guide](../CONTRIBUTING.md) for detailed instructions.

## Quick Links

- 📦 [GitHub Repository](https://github.com/CatoSystems/rim-pay)
- 🚀 [Quick Start Guide](quick-start.md)
- 📖 [API Reference](api/reference.md)
- 💡 [Examples](examples/README.md)
- ❓ [FAQ](faq.md)
- 🐛 [Report Issues](https://github.com/CatoSystems/rim-pay/issues)
- 💬 [Discussions](https://github.com/CatoSystems/rim-pay/discussions)

## Support

For questions, issues, or contributions:

1. Check the [FAQ](faq.md)
2. Search [existing issues](https://github.com/CatoSystems/rim-pay/issues)
3. Create a new issue with detailed information
4. For security issues, contact us privately

## License

This project is licensed under the MIT License - see the [LICENSE](../LICENSE) file for details.
