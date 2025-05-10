package configs

type Config interface {
	GetServer() ServerConfig
	GetDatabase() DatabaseConfig
	GetProviderConfig() ProviderConfig
	GetScheduler() SchedulerConfig
}

type configImp struct {
	Server    serverConfigImp    `mapstructure:"server" validate:"required"`
	Database  databaseConfigImp  `mapstructure:"database" validate:"required"`
	Provider  ProviderConfig     `mapstructure:"-" validate:"required"`
	Scheduler schedulerConfigImp `mapstructure:"scheduler" validate:"required"`
}

func (c configImp) GetServer() ServerConfig {
	return c.Server
}

func (c configImp) GetDatabase() DatabaseConfig {
	return c.Database
}

func (c configImp) GetProviderConfig() ProviderConfig {
	return c.Provider
}

func (c configImp) GetScheduler() SchedulerConfig {
	return c.Scheduler
}
