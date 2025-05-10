package sender

import (
	"github.com/wagslane/go-rabbitmq"
	"insider/configs"
	"insider/graceful_shutdown"
	"insider/message"
	"sync"
	"time"
)

var instance *scheduler
var connectionOnce sync.Once
var messageService message.UnsentMessageService

type scheduler struct {
	interval time.Duration
	stop     chan struct{}
}

func Start() {
	connectionOnce.Do(func() {
		publisher := configureRabbitMQ()
		configureScheduler(publisher)
	})
}

func configureScheduler(publisher *rabbitmq.Publisher) {
	messageService = message.GetUnsentMessageService()

	graceful_shutdown.AddShutdownHook(func() {
		close(instance.stop)
	})

	schedulerConfig := configs.Instance().
		GetScheduler()

	instance = &scheduler{
		interval: schedulerConfig.GetInterval(),
		stop:     make(chan struct{}),
	}

	go runner(publisher)
}
