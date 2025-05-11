package rabbitmq

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/streadway/amqp"
	"github.com/wagslane/go-rabbitmq"
	"insider/configs"
	"insider/graceful_shutdown"
	"sync"
)

var (
	adminConnection      *amqp.Connection
	channel              *amqp.Channel
	processingConnection *rabbitmq.Conn
	connectionOnce       sync.Once
)

func configure() {
	connectionOnce.Do(func() {
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
	})
}

func DeclareQueue(name string, durable, autoDelete, exclusive, noWait bool, args amqp.Table) error {
	configure()

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
