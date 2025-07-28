# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.2.0] - 2025-07-28

### 🔒 Security
- **BREAKING CHANGE**: BPay payments now automatically generate secure 4-digit passcodes
- Implement mandatory passcode generation using `crypto/rand` for cryptographic security
- Library now ignores any user-provided passcodes to prevent weak or predictable codes
- All generated passcodes are guaranteed to be in secure range 1000-9999

### ✨ Added
- Automatic secure passcode generation for all BPay payments
- Comprehensive test suite including integration and security tests
- Passcode generation validation and uniqueness testing
- Enhanced logging for passcode generation debugging
- Generated passcode returned in response metadata for payer use

### 🔧 Changed
- **BREAKING**: `BPayPaymentRequest.Passcode` field is now ignored during processing
- **BREAKING**: Removed user passcode validation (no longer needed)
- Updated examples and documentation to reflect new passcode behavior
- Enhanced BPay payment processor to always generate fresh passcodes

### 🐛 Fixed  
- Corrected currency code from deprecated MRO to current MRU standard
- Improved documentation structure and comprehensive examples

### 📚 Documentation
- Added comprehensive documentation and restructured examples
- Updated API documentation to reflect passcode generation behavior
- Added security notes about automatic passcode generation
- Enhanced integration examples with detailed explanations

### 🔄 Migration Guide
For users upgrading from v0.1.0 to v0.2.0:

```go
// Before (v0.1.0) - passcode was required
bpayRequest := &rimpay.BPayPaymentRequest{
    PhoneNumber: phoneNum,
    Amount:      amount,
    Passcode:    "1234", // This was required
    // ...
}

// After (v0.2.0) - passcode is auto-generated and ignored
bpayRequest := &rimpay.BPayPaymentRequest{
    PhoneNumber: phoneNum,
    Amount:      amount,
    // Passcode field can be omitted or will be ignored
    // ...
}

// Extract generated passcode from response
if passcode, exists := response.Metadata["passcode"]; exists {
    // Use the generated passcode for the payer
    fmt.Printf("Generated passcode: %s", passcode)
}
```

## [0.1.0] - 2025-07-27

### Added
- Initial release of RimPay library
- B-PAY payment provider support
- MASRVI payment provider support
- Provider-specific request types (`BPayPaymentRequest`, `MasrviPaymentRequest`)
- Type-safe client methods (`ProcessBPayPayment`, `ProcessMasrviPayment`) 
- Mauritanian phone number validation with prefixes 2, 3, 4
- Decimal-based money handling with MRU currency support
- Comprehensive error handling with specific error types
- Configurable retry mechanisms with exponential backoff
- Multi-provider configuration system
- Payment status checking for B-PAY
- Complete test suite
- Comprehensive examples and documentation

### Features
- **Multi-Provider Architecture**: Support for multiple payment providers with type-safe APIs
- **Phone Validation**: Built-in validation for Mauritanian phone numbers
- **Money Precision**: Decimal-based calculations for accurate money handling
- **Error Handling**: Detailed error types for validation, provider, and network errors
- **Retry Logic**: Configurable retry mechanisms for resilient payment processing
- **Configuration**: Flexible configuration for different environments and providers

### Documentation
- Complete README with usage examples
- API reference documentation
- Provider-specific usage guides
- Error handling patterns
- Configuration examples
- Comprehensive example applications

[0.2.0]: https://github.com/CatoSystems/rim-pay/releases/tag/v0.2.0
[0.1.0]: https://github.com/CatoSystems/rim-pay/releases/tag/v0.1.0
