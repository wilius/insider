package provider

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/sony/gobreaker"
	"insider/configs"
	"insider/constants"
	"io"
	"net/http"
)

type webhookSiteRequest struct {
	To      string `json:"to"`
	Content string `json:"content"`
}

type webhookSiteResponse struct {
	MessageId string `json:"messageId"`
}

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

func (s *senderImp) Type() constants.ProviderType {
	return constants.WebhookSite
}

func (s *senderImp) Send(input *SendMessageInput) (*SendMessageOutput, error) {
	if input == nil {
		return nil, errors.New("input is nil")
	}

	request := webhookSiteRequest{
		To:      input.PhoneNumber,
		Content: input.Message,
	}

	response, err := s.circuitBreaker.Execute(func() (interface{}, error) {
		reader, writer := io.Pipe()
		go func() {
			defer writer.Close()
			if err := json.NewEncoder(writer).Encode(request); err != nil {
				log.Err(err).
					Msg("failed to encode request")
			}
		}()

		log.Trace().
			Msgf("[HTTP] Sending request to %s", s.config.GetUrl())

		req, err := http.NewRequest("POST", s.config.GetUrl(), reader)
		if err != nil {
			return nil, err
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")

		resp, err := s.client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode < 200 ||
			resp.StatusCode >= 300 {
			return nil, fmt.Errorf("server error: %d", resp.StatusCode)
		}

		var result webhookSiteResponse
		err = json.NewDecoder(resp.Body).
			Decode(&result)

		if err != nil {
			return nil, fmt.Errorf("failed to decode response: %w", err)
		}

		return result, nil
	})

	if err != nil {
		return nil, err
	}

	casted := response.(webhookSiteResponse)
	return &SendMessageOutput{
		MessageId: casted.MessageId,
	}, nil
}
