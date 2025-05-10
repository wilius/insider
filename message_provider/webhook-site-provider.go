package message_provider

import (
	"errors"
	"fmt"
	"github.com/sony/gobreaker"
	"insider/configs"
	"io/ioutil"
	"net/http"
)

func newWebhookSiteProvider(config configs.WebhookSiteConfig) Provider {
	client := &http.Client{
		Timeout: config.GetRequestTimeout(),
	}

	breakerConfig := config.GetCircuitBreakerConfig()
	settings := gobreaker.Settings{
		Name:        "HTTP Client Circuit Breaker",
		MaxRequests: breakerConfig.GetMaxRequests(),
		Interval:    breakerConfig.GetInterval(), //60 * time.Second
		Timeout:     breakerConfig.GetTimeout(),
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures > breakerConfig.GetMaxFailure()
		},
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			fmt.Printf("[Circuit Breaker] State change: %s → %s\n", from.String(), to.String())
		},
	}

	cb := gobreaker.NewCircuitBreaker(settings)
	return &senderImp{
		config:         config,
		client:         client,
		circuitBreaker: cb,
	}
}

type senderImp struct {
	config         configs.WebhookSiteConfig
	client         *http.Client
	circuitBreaker *gobreaker.CircuitBreaker
}

func (s *senderImp) Send(input *SendMessageInput) (*SendMessageOutput, error) {
	if input == nil {
		return nil, errors.New("input is nil")
	}

	_, err := s.circuitBreaker.Execute(func() (interface{}, error) {
		req, err := http.NewRequest("GET", s.config.GetUrl(), nil)
		if err != nil {
			return nil, err
		}

		resp, err := s.client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode >= 500 {
			return nil, fmt.Errorf("server error: %d", resp.StatusCode)
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		return string(body), nil
	})

	if err != nil {
		return nil, err
	}

	// TODO
	return nil, nil
}
