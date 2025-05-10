package configs

import (
	"insider/constants"
	"time"
)

type WebhookSiteConfig interface {
	GetUrl() string
	GetRequestTimeout() time.Duration
	GetCircuitBreakerConfig() WebhookSiteCircuitBreakerConfig
}

type WebhookSiteCircuitBreakerConfig interface {
	GetMaxRequests() uint32
	GetMaxFailure() uint32
	GetInterval() time.Duration
	GetTimeout() time.Duration
}

type webhookSiteConfigImp struct {
	providerConfigImp `mapstructure:",squash"`
	Url               string                             `mapstructure:"url" validate:"required,url"`
	RequestTimeout    time.Duration                      `mapstructure:"requestTimeout" validate:"required,min=0s,max=5s"`
	CircuitBreaker    WebhookSiteCircuitBreakerConfigImp `mapstructure:"circuitBreaker" validate:"required"`
}

func (w webhookSiteConfigImp) GetUrl() string {
	return w.Url
}

func (w webhookSiteConfigImp) GetType() constants.ProviderType {
	return w.Type
}

func (w webhookSiteConfigImp) GetRequestTimeout() time.Duration {
	return w.RequestTimeout
}

func (w webhookSiteConfigImp) GetCircuitBreakerConfig() WebhookSiteCircuitBreakerConfig {
	return w.CircuitBreaker
}

type WebhookSiteCircuitBreakerConfigImp struct {
	MaxRequests uint32        `mapstructure:"maxRequests" validate:"required"`
	MaxFailure  uint32        `mapstructure:"maxFailure" validate:"required"`
	Interval    time.Duration `mapstructure:"interval" validate:"required"`
	Timeout     time.Duration `mapstructure:"timeout" validate:"required"`
}

func (w WebhookSiteCircuitBreakerConfigImp) GetMaxRequests() uint32 {
	return w.MaxRequests
}

func (w WebhookSiteCircuitBreakerConfigImp) GetMaxFailure() uint32 {
	return w.MaxFailure
}

func (w WebhookSiteCircuitBreakerConfigImp) GetInterval() time.Duration {
	return w.Interval
}

func (w WebhookSiteCircuitBreakerConfigImp) GetTimeout() time.Duration {
	return w.Timeout
}
