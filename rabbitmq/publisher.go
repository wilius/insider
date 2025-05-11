package rabbitmq

import (
	"encoding/json"
	"github.com/rs/zerolog/log"
	"github.com/wagslane/go-rabbitmq"
	"insider/graceful_shutdown"
)

type Publisher struct {
	exchange  string
	publisher *rabbitmq.Publisher
}

func NewPublisher(
	exchange string,
	optionFunctions ...func(*rabbitmq.PublisherOptions)) (*Publisher, error) {
	configure()

	publisher, err := rabbitmq.NewPublisher(
		processingConnection,
		optionFunctions...,
	)

	if err != nil {
		if publisher != nil {
			publisher.Close()
		}

		return nil, err
	}

	graceful_shutdown.AddShutdownHook(publisher.Close)
	return &Publisher{
		exchange:  exchange,
		publisher: publisher,
	}, nil
}

func (p *Publisher) PublishEvent(
	object interface{},
	routingKey string,
) error {
	jsonData, err := json.Marshal(object)
	if err != nil {
		return err
	}

	err = p.publisher.Publish(
		jsonData,
		[]string{routingKey},
		rabbitmq.WithPublishOptionsExchange(p.exchange),
		rabbitmq.WithPublishOptionsContentType("application/json"),
		rabbitmq.WithPublishOptionsMandatory,
	)
	if err != nil {
		return err
	}

	log.Info().Msgf("Message published successfully!")
	return nil
}
