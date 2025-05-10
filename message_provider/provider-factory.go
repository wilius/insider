package message_provider

import (
	"errors"
	"insider/configs"
	"insider/constants"
)

var instance Provider

func Instance() Provider {
	return instance
}

func init() {
	providerConfig := configs.Instance().
		GetProviderConfig()

	switch providerConfig.GetType() {
	case constants.WebhookSite:
		instance = newWebhookSiteProvider(providerConfig.(configs.WebhookSiteConfig))
	default:
		panic(errors.New("couldn't find any matching provider"))
	}
}
