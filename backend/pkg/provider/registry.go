package provider

import (
	"fmt"
	"lytemp/config"
)

// BuildProviders creates enabled providers from configuration.
// NOTE: Ensure config.Config has `Provider Provider `mapstructure:"provider"` field.
func BuildProviders() ([]Provider, error) {
	cfg := config.Get()

	var providers []Provider
	client := HTTPClient()

	if cfg.Provider.Provider1 != "" {
		providers = append(providers, NewJSONProvider(cfg.Provider.Provider1, client))
	}
	if cfg.Provider.Provider2 != "" {
		providers = append(providers, NewXMLProvider(cfg.Provider.Provider2, client))
	}

	if len(providers) == 0 {
		return nil, fmt.Errorf("no providers configured")
	}
	return providers, nil
}
