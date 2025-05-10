package configs

type RabbitmqConfig interface {
	GetHost() string
	GetPort() uint16
	GetVHost() string
	GetUsername() string
	GetPassword() string
}

type rabbitmqConfigImp struct {
	Host     string `mapstructure:"host" validate:"required,hostname"`
	Port     uint16 `mapstructure:"port" validate:"required,min=1,max=65535"`
	VHost    string `mapstructure:"vhost" validate:"required"`
	Username string `mapstructure:"username" validate:"required"`
	Password string `mapstructure:"password" validate:"required"`
}

func (d rabbitmqConfigImp) GetHost() string {
	return d.Host
}

func (d rabbitmqConfigImp) GetPort() uint16 {
	return d.Port
}

func (d rabbitmqConfigImp) GetVHost() string {
	return d.VHost
}

func (d rabbitmqConfigImp) GetUsername() string {
	return d.Username
}

func (d rabbitmqConfigImp) GetPassword() string {
	return d.Password
}
