package configs

import (
	"gorm.io/gorm/logger"
	"time"
)

type DatabaseConfig interface {
	GetHost() string
	GetPort() uint16
	GetName() string
	GetUsername() string
	GetPassword() string
	GetPool() Pool
	GetLogLevel() logger.LogLevel
}

type Pool interface {
	GetMaxIdleConnections() int
	GetMaxOpenConnections() int
	GetConnectionMaxLifetime() time.Duration
	GetConnectionMaxIdleTime() time.Duration
}

type databaseConfigImp struct {
	Host     string          `mapstructure:"host" validate:"required,hostname"`
	Port     uint16          `mapstructure:"port" validate:"required,min=1,max=65535"`
	Name     string          `mapstructure:"name" validate:"required"`
	Username string          `mapstructure:"username" validate:"required"`
	Password string          `mapstructure:"password" validate:"required"`
	Pool     poolImp         `mapstructure:"pool" validate:"required"`
	LogLevel logger.LogLevel `mapstructure:"logLevel" validate:"required,min=1,max=4"`
}

func (d databaseConfigImp) GetHost() string {
	return d.Host
}

func (d databaseConfigImp) GetPort() uint16 {
	return d.Port
}

func (d databaseConfigImp) GetName() string {
	return d.Name
}

func (d databaseConfigImp) GetUsername() string {
	return d.Username
}

func (d databaseConfigImp) GetPassword() string {
	return d.Password
}

func (d databaseConfigImp) GetPool() Pool {
	return d.Pool
}

func (d databaseConfigImp) GetLogLevel() logger.LogLevel {
	return d.LogLevel
}

type poolImp struct {
	MaxIdleConnections    int           `mapstructure:"maxIdleConnections" validate:"required"`
	MaxOpenConnections    int           `mapstructure:"maxOpenConnections" validate:"required"`
	ConnectionMaxLifetime time.Duration `mapstructure:"connectionMaxLifetime" validate:"required"`
	ConnectionMaxIdleTime time.Duration `mapstructure:"connectionMaxIdleTime" validate:"required"`
}

func (p poolImp) GetMaxIdleConnections() int {
	return p.MaxIdleConnections
}

func (p poolImp) GetMaxOpenConnections() int {
	return p.MaxOpenConnections
}

func (p poolImp) GetConnectionMaxLifetime() time.Duration {
	return p.ConnectionMaxLifetime
}

func (p poolImp) GetConnectionMaxIdleTime() time.Duration {
	return p.ConnectionMaxIdleTime
}
