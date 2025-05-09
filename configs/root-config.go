package configs

type Config interface {
	GetServer() ServerConfig
	GetDatabase() DatabaseConfig
}

type configImp struct {
	Server   serverConfigImp   `mapstructure:"server" validate:"required"`
	Database databaseConfigImp `mapstructure:"database" validate:"required"`
}

func (c configImp) GetServer() ServerConfig {
	return c.Server
}

func (c configImp) GetDatabase() DatabaseConfig {
	return c.Database
}
