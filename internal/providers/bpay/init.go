package bpay

import (
	"github.com/CatoSystems/rim-pay/pkg/rimpay"
)

// init registers the bpay provider with the default registry
func init() {
	rimpay.DefaultRegistry.Register("bpay", func(config rimpay.ProviderConfig, logger rimpay.Logger) (rimpay.PaymentProvider, error) {
		return NewProvider(config, logger)
	})
}
