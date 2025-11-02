// Package providers provides initialization for all payment providers.package providers

// Import this package to register all available providers with the RimPay client.
package providers

import (
	// Import provider packages to trigger their init() functions
	// which register the providers with the RimPay client
	_ "github.com/CatoSystems/rim-pay/internal/providers/bpay"
	_ "github.com/CatoSystems/rim-pay/internal/providers/masrvi"
)
