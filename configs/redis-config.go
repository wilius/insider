package configs

type RedisConfig interface {
	GetHost() string
	GetPort() uint16
}

type redisImp struct {
	Host string `mapstructure:"host" validate:"required,hostname"`
	Port uint16 `mapstructure:"port" validate:"required,min=1,max=65535"`
}

func (d redisImp) GetHost() string {
	return d.Host
}

func (d redisImp) GetPort() uint16 {
	return d.Port
}
