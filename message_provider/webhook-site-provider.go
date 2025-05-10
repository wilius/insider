package message_provider

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/sony/gobreaker"
)

func newWebhookSiteProvider() Provider {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Configure Circuit Breaker
	settings := gobreaker.Settings{
		Name:        "HTTP Client Circuit Breaker",
		MaxRequests: 3,
		Interval:    60 * time.Second,
		Timeout:     5 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures > 5
		},
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			fmt.Printf("[Circuit Breaker] State change: %s → %s\n", from.String(), to.String())
		},
	}

	cb := gobreaker.NewCircuitBreaker(settings)
	return &senderImp{
		client:         client,
		circuitBreaker: cb,
	}
}

type senderImp struct {
	client         *http.Client
	circuitBreaker *gobreaker.CircuitBreaker
}

func (s *senderImp) Send(input *SendMessageInput) (*SendMessageOutput, error) {
	if input == nil {
		return nil, errors.New("input is nil")
	}

	_, err := s.circuitBreaker.Execute(func() (interface{}, error) {
		// TODO read url from the config
		req, err := http.NewRequest("GET", "url", nil)
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
