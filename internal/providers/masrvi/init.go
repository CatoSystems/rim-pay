package masrvi

import (
	"github.com/CatoSystems/rim-pay/pkg/rimpay"
)

// init registers the masrvi provider with the default registry
func init() {
	rimpay.DefaultRegistry.Register("masrvi", func(config rimpay.ProviderConfig, logger rimpay.Logger) (rimpay.PaymentProvider, error) {
		return NewProvider(config, logger)
	})
}
