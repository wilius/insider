package configs

import (
	"fmt"
	"github.com/go-viper/mapstructure/v2"
	"insider/constants"
	"insider/util"
)

type ProviderConfig interface {
	GetType() constants.ProviderType
}

type providerConfigImp struct {
	Type constants.ProviderType `mapstructure:"type" validate:"required"`
}

func (p providerConfigImp) GetType() constants.ProviderType {
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

	decoderConfig := &mapstructure.DecoderConfig{
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			stringToDurationParserHook(),
		),
		Result:           &provider,
		TagName:          "mapstructure",
		WeaklyTypedInput: true, // Allow Viper to cast
	}

	decoder, err := mapstructure.NewDecoder(decoderConfig)
	if err != nil {
		panic(fmt.Errorf("failed to create decoder: %w", err))
	}

	if err = decoder.Decode(*data); err != nil {
		return nil, err
	}

	// Validate the constructed provider
	if err = util.Validate(provider); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	return provider, nil
}

func getProviderConfig(providerType constants.ProviderType) (ProviderConfig, error) {
	switch providerType {
	case constants.WebhookSite:
		return webhookSiteConfigImp{}, nil
	default:
		return nil, fmt.Errorf("unsupported provider type: %s", providerType)
	}
}
