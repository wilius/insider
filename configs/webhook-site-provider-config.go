package configs

import "insider/message_provider"

type webhookSiteConfigImp struct {
	providerConfigImp `mapstructure:",squash"`
	Url               string `mapstructure:"url" validate:"required,url"`
}

func (w webhookSiteConfigImp) GetType() message_provider.ProviderType {
	return w.Type
}
