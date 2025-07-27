package providers

import (
	"github.com/CatoSystems/rim-pay/internal/providers/bpay"
	"github.com/CatoSystems/rim-pay/internal/providers/masrvi"
	"github.com/CatoSystems/rim-pay/pkg/rimpay"
)

// RegisterAll registers all available providers
func RegisterAll(registry *rimpay.ProviderRegistry) {
	// Register B-PAY provider
	registry.Register("bpay", func(config rimpay.ProviderConfig, logger rimpay.Logger) (rimpay.PaymentProvider, error) {
		return bpay.NewProvider(config, logger)
	})

	// Register MASRVI provider
	registry.Register("masrvi", func(config rimpay.ProviderConfig, logger rimpay.Logger) (rimpay.PaymentProvider, error) {
		return masrvi.NewProvider(config, logger)
	})
}
