package configs

type ServerConfig interface {
	GetPort() uint16
}

type serverConfigImp struct {
	Port uint16 `mapstructure:"port" validate:"required,min=1,max=65535"`
}

func (s serverConfigImp) GetPort() uint16 {
	return s.Port
}
