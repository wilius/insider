package sender

import (
	"insider/configs"
	"insider/graceful_shutdown"
	"insider/message"
	"insider/rabbitmq"
	"sync"
	"time"
)

var (
	instance       *scheduler
	connectionOnce sync.Once
	messageService message.UnsentMessageService
	publisher      *rabbitmq.Publisher
)

type scheduler struct {
	interval time.Duration
	stop     chan struct{}
}

func Start() {
	connectionOnce.Do(func() {
		configureRabbitMQ()
		configureScheduler()
	})
}

func configureScheduler() {
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

	go runner()
}
