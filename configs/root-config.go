package configs

type Config interface {
	GetServer() ServerConfig
	GetDatabase() DatabaseConfig
	GetRabbitMQ() RabbitmqConfig
	GetRedis() RedisConfig
	GetScheduler() SchedulerConfig
	GetProviderConfig() ProviderConfig
}

type configImp struct {
	Server      serverConfigImp    `mapstructure:"server" validate:"required"`
	Database    databaseConfigImp  `mapstructure:"database" validate:"required"`
	Scheduler   schedulerConfigImp `mapstructure:"scheduler" validate:"required"`
	RabbitMQ    rabbitmqConfigImp  `mapstructure:"rabbitmq" validate:"required"`
	RedisConfig redisImp           `mapstructure:"redis" validate:"required"`
	Provider    ProviderConfig     `mapstructure:"-" validate:"required"`
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

func (c configImp) GetRabbitMQ() RabbitmqConfig {
	return c.RabbitMQ
}
func (c configImp) GetRedis() RedisConfig {
	return c.RedisConfig
}

func (c configImp) GetScheduler() SchedulerConfig {
	return c.Scheduler
}
