package click

import "github.com/CatoSystems/rim-pay/pkg/rimpay"

// init registers the click provider with the default registry.
func init() {
	rimpay.DefaultRegistry.Register("click", func(config rimpay.ProviderConfig, logger rimpay.Logger) (rimpay.PaymentProvider, error) {
		return NewProvider(config, logger)
	})
}
