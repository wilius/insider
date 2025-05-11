package scheduler

import (
	"github.com/rs/zerolog/log"
	"time"
)

type scheduler struct {
	interval time.Duration
	worker   func()
	stop     chan struct{}
}

func (s scheduler) Start() {
	go s.doStart()
}

func (s scheduler) doStart() {
	s.worker()
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			ticker.Stop()
			s.worker()
			ticker.Reset(s.interval)
		case <-s.stop:
			log.Info().
				Msg("Scheduler stopped")
			return
		}
	}
}

func (s scheduler) Stop() {
	close(s.stop)
}
