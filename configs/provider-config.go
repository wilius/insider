package configs

import (
	"fmt"
	"github.com/go-viper/mapstructure/v2"
	"insider/message_provider"
	"insider/util"
)

type ProviderConfig interface {
	GetType() message_provider.ProviderType
}

type providerConfigImp struct {
	Type message_provider.ProviderType `mapstructure:"type" validate:"required"`
}

func (p providerConfigImp) GetType() message_provider.ProviderType {
	return p.Type
}

func providerConfigFactory(data *map[string]interface{}) (ProviderConfig, error) {
	base := providerConfigImp{}

	if err := mapstructure.Decode(*data, &base); err != nil {
		return nil, err
	}

	provider, err := getProviderConfig(base.Type)
	if err != nil {
		return nil, err
	}

	if err = mapstructure.Decode(*data, &provider); err != nil {
		return nil, err
	}

	// Validate the constructed provider
	if err = util.Validate(provider); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	return provider, nil
}

func getProviderConfig(providerType message_provider.ProviderType) (ProviderConfig, error) {
	switch providerType {
	case message_provider.WebhookSite:
		return webhookSiteConfigImp{}, nil
	default:
		return nil, fmt.Errorf("unsupported provider type: %s", providerType)
	}
}
