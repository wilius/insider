package message_provider

/*func providerFactory(data map[string]interface{}) (interface{}, error) {
	// Decode the base config to identify the type
	var base configs.ProviderConfig
	if err := mapstructure.Decode(data, &base); err != nil {
		return nil, err
	}

	// Determine the type and decode to the correct struct
	switch base.Type {
	case WebhookSite:
		var config configs.webhookSiteConfigImp
		if err := mapstructure.Decode(data, &config); err != nil {
			return nil, err
		}
		return config, nil

	default:
		return nil, fmt.Errorf("unknown provider type: %s", base.Type)
	}
}
*/
