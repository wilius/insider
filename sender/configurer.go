package sender

import (
	"context"
	"github.com/rs/zerolog/log"
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

		go func() {
			doRun() // Run the task immediately
			ticker := time.NewTicker(instance.interval)
			defer ticker.Stop()

			for {
				select {
				case <-ticker.C:
					ticker.Stop()
					doRun()
					ticker.Reset(instance.interval)
				case <-instance.stop:
					log.Info().
						Msg("Scheduler stopped")
					return
				}
			}
		}()
	})
}

func doRun() {
	log.Info().
		Msg("Running scheduler")
	count := configs.Instance().
		GetScheduler().
		GetItemCountPerCycle()

	items, err := messageService.Fetch(context.Background(), count)
	if err != nil {
		log.Err(err).
			Msg("Error while fetching items")
		return
	}

	for _, item := range *items {
		log.Info().
			Msgf("Sending message to %s with content %s", item.PhoneNumber, item.Message)
	}
}
