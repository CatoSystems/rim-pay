/*
Package rimpay provides a comprehensive payment processing library for Mauritania,
supporting multiple payment providers including B-PAY and MASRVI.

The library offers type-safe, provider-specific APIs with built-in validation,
error handling, and retry mechanisms.

# Quick Start

	import (
		"context"
		"github.com/CatoSystems/rim-pay/pkg/rimpay"
		"github.com/CatoSystems/rim-pay/pkg/phone"
		"github.com/CatoSystems/rim-pay/pkg/money"
		"github.com/shopspring/decimal"
	)

	// Create client
	config := rimpay.DefaultConfig()
	config.DefaultProvider = "bpay"
	client, err := rimpay.NewClient(config)

	// Create payment request
	phone, _ := phone.NewPhone("+22233445566")
	amount := money.New(decimal.NewFromInt(10000), "MRU")

	request := &rimpay.BPayPaymentRequest{
		Amount:      amount,
		PhoneNumber: phone,
		Reference:   "ORDER-12345",
		Passcode:    "1234",
	}

	// Process payment
	response, err := client.ProcessBPayPayment(context.Background(), request)

# Providers

The library supports multiple payment providers:

B-PAY: Mauritanian mobile payment provider requiring passcode authentication.
MASRVI: Web-based payment provider with callback/webhook support.

# Phone Numbers

Phone number validation is built-in for Mauritanian numbers:

	phone, err := phone.NewPhone("+22233445566") // Valid
	phone, err := phone.NewPhone("+22255667788") // Invalid (prefix 5)

Valid prefixes: 2, 3, 4 (8 digits including prefix)

# Money Handling

The library uses decimal precision for accurate financial calculations:

	amount := money.New(decimal.NewFromInt(10050), "MRU") // 100.50 MRU
	amount := money.FromFloat64(100.50, "MRU")           // Same as above

# Error Handling

Comprehensive error types are provided:

	if err != nil {
		switch e := err.(type) {
		case *rimpay.ValidationError:
			// Handle validation errors
		case *rimpay.ProviderError:
			// Handle provider-specific errors
		case *rimpay.NetworkError:
			// Handle network/connectivity errors
		}
	}

# Configuration

Flexible configuration system supports multiple environments:

	config := rimpay.DefaultConfig()
	config.Environment = rimpay.EnvironmentProduction
	config.Providers["bpay"] = rimpay.ProviderConfig{
		Enabled: true,
		BaseURL: "https://api.bpay.mr",
		Credentials: map[string]string{
			"username": "your_username",
			"password": "your_password",
		},
	}

For more detailed examples and documentation, see:
https://github.com/CatoSystems/rim-pay
*/
package rimpay
