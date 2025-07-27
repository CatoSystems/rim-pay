package rimpay

import (
	"fmt"
)

// ProviderFactory creates payment providers
type ProviderFactory func(config ProviderConfig, logger Logger) (PaymentProvider, error)

// ProviderRegistry manages payment provider factories
type ProviderRegistry struct {
	factories map[string]ProviderFactory
}

// NewProviderRegistry creates a new provider registry
func NewProviderRegistry() *ProviderRegistry {
	return &ProviderRegistry{
		factories: make(map[string]ProviderFactory),
	}
}

// Register registers a provider factory
func (r *ProviderRegistry) Register(name string, factory ProviderFactory) {
	r.factories[name] = factory
}

// Create creates a provider instance
func (r *ProviderRegistry) Create(name string, config ProviderConfig, logger Logger) (PaymentProvider, error) {
	factory, exists := r.factories[name]
	if !exists {
		return nil, fmt.Errorf("unknown provider: %s", name)
	}
	return factory(config, logger)
}

// GetRegisteredProviders returns list of registered provider names
func (r *ProviderRegistry) GetRegisteredProviders() []string {
	names := make([]string, 0, len(r.factories))
	for name := range r.factories {
		names = append(names, name)
	}
	return names
}

// DefaultRegistry is the default global provider registry
var DefaultRegistry = NewProviderRegistry()
