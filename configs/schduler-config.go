package configs

import "time"

type SchedulerConfig interface {
	GetInterval() time.Duration
	GetItemCountPerCycle() uint
}

type schedulerConfigImp struct {
	Interval          time.Duration `mapstructure:"interval" validate:"required,min=1s"`
	ItemCountPerCycle uint          `mapstructure:"itemCountPerCycle" validate:"required,min=1"`
}

func (s schedulerConfigImp) GetInterval() time.Duration {
	return s.Interval
}

func (s schedulerConfigImp) GetItemCountPerCycle() uint {
	return s.ItemCountPerCycle
}
