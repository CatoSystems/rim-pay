package rimpay

import "fmt"

// AddBPayProvider adds a B-PAY provider to the client
func (c *Client) AddBPayProvider(config ProviderConfig) error {
	if createBPayProvider == nil {
		return fmt.Errorf("B-PAY provider not registered")
	}
	
	// Create provider using the registered factory
	provider, err := createBPayProvider(config, c.logger)
	if err != nil {
		return err
	}
	return c.AddProvider(ProviderBPay, provider)
}

// AddMasrviProvider adds a MASRVI provider to the client
func (c *Client) AddMasrviProvider(config ProviderConfig) error {
	if createMasrviProvider == nil {
		return fmt.Errorf("MASRVI provider not registered")
	}
	
	// Create provider using the registered factory
	provider, err := createMasrviProvider(config, c.logger)
	if err != nil {
		return err
	}
	return c.AddProvider(ProviderMasrvi, provider)
}

// GetBPayProvider returns the B-PAY provider if available
func (c *Client) GetBPayProvider() (BPayProvider, error) {
	provider, ok := c.providers[ProviderBPay]
	if !ok {
		return nil, ErrProviderNotFound
	}

	bpayProvider, ok := provider.(BPayProvider)
	if !ok {
		return nil, ErrInvalidProvider
	}

	return bpayProvider, nil
}

// GetMasrviProvider returns the MASRVI provider if available
func (c *Client) GetMasrviProvider() (MasrviProvider, error) {
	provider, ok := c.providers[ProviderMasrvi]
	if !ok {
		return nil, ErrProviderNotFound
	}

	masrviProvider, ok := provider.(MasrviProvider)
	if !ok {
		return nil, ErrInvalidProvider
	}

	return masrviProvider, nil
}
