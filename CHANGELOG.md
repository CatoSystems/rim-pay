# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

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

[0.1.0]: https://github.com/CatoSystems/rim-pay/releases/tag/v0.1.0
