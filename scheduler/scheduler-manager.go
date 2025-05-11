package scheduler

import (
	"insider/configs"
	"sync"
)

type Manager struct {
	mu        sync.Mutex
	scheduler *scheduler
	worker    func()
}

func NewManager(worker func()) *Manager {
	managerInstance := &Manager{
		worker: worker,
	}

	managerInstance.Start()
	return managerInstance
}

func (s *Manager) Start() {
	if s.scheduler == nil {
		s.mu.Lock()
		defer s.mu.Unlock()
		if s.scheduler == nil {
			schedulerConfig := configs.Instance().
				GetScheduler()

			s.scheduler = &scheduler{
				interval: schedulerConfig.GetInterval(),
				stop:     make(chan struct{}),
				worker:   s.worker,
			}

			s.scheduler.Start()
		}
	}
}

func (s *Manager) Stop() {
	if s.scheduler != nil {
		s.mu.Lock()
		defer s.mu.Unlock()
		if s.scheduler != nil {
			s.scheduler.Stop()
			s.scheduler = nil
		}
	}
}
