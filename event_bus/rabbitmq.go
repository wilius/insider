package event_bus

import (
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/streadway/amqp"
	"github.com/wagslane/go-rabbitmq"
	"insider/configs"
	"insider/graceful_shutdown"
)

var adminConnection *amqp.Connection
var processingConnection *rabbitmq.Conn
var channel *amqp.Channel

func init() {
	rabbitmqConfig := configs.Instance().
		GetRabbitMQ()

	var err error
	rabbitmqConnectionString := fmt.Sprintf(
		"amqp://%s:%s@%s:%d/%s",
		rabbitmqConfig.GetUsername(),
		rabbitmqConfig.GetPassword(),
		rabbitmqConfig.GetHost(),
		rabbitmqConfig.GetPort(),
		rabbitmqConfig.GetVHost(),
	)
	processingConnection, err = rabbitmq.NewConn(
		rabbitmqConnectionString,
		rabbitmq.WithConnectionOptionsLogging,
	)

	adminConnection, err = amqp.Dial(rabbitmqConnectionString)
	if err != nil {
		log.Fatal().
			Err(err).
			Msgf("%v", err)
	}

	if err != nil {
		log.Fatal().Err(err)
	}

	channel, err = adminConnection.Channel()
	if err != nil {
		log.Fatal().Err(err)
	}

	graceful_shutdown.AddShutdownHook(func() {
		_ = processingConnection.Close()
		_ = adminConnection.Close()
		_ = channel.Close()
	})
}

func DeclareQueue(name string, durable, autoDelete, exclusive, noWait bool, args amqp.Table) error {
	_, err := channel.QueueDeclare(
		name,
		durable,
		autoDelete,
		exclusive,
		noWait,
		args,
	)

	return err
}

func NewConsumer(
	queue string,
	optionFunctions ...func(*rabbitmq.ConsumerOptions),
) (*rabbitmq.Consumer, error) {
	optionFunctions = append(
		optionFunctions,
		rabbitmq.WithConsumerOptionsConsumerAutoAck(false),
	)

	consumer, err := rabbitmq.NewConsumer(
		processingConnection,
		queue,
		optionFunctions...,
	)
	if err != nil {
		if consumer != nil {
			consumer.Close()
		}

		return nil, err
	}

	graceful_shutdown.AddShutdownHook(consumer.Close)
	return consumer, nil
}

func NewPublisher(optionFunctions ...func(*rabbitmq.PublisherOptions)) (*rabbitmq.Publisher, error) {
	eventPublisher, err := rabbitmq.NewPublisher(
		processingConnection,
		optionFunctions...,
	)

	if err != nil {
		if eventPublisher != nil {
			eventPublisher.Close()
		}

		return nil, err
	}

	graceful_shutdown.AddShutdownHook(eventPublisher.Close)
	return eventPublisher, nil
}

func PublishEvent(
	publisher *rabbitmq.Publisher,
	object interface{},
	exchange string,
	routingKey string,
) {
	jsonData, err := json.Marshal(object)
	if err != nil {
		log.Error().Msgf("Failed to serialize object to JSON : %v", err)
	}

	err = publisher.Publish(
		jsonData,
		[]string{routingKey},
		rabbitmq.WithPublishOptionsExchange(exchange),
		rabbitmq.WithPublishOptionsContentType("application/json"),
		rabbitmq.WithPublishOptionsMandatory,
	)
	if err != nil {
		log.Error().Msgf("Failed to publish message: %v", err)
	} else {
		log.Info().Msgf("Message published successfully!")
	}

}
