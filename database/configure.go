package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"insider/configs"
	"sync"
)

var (
	instance       *gorm.DB
	connectionOnce sync.Once
)

func Instance() *gorm.DB {
	connectionOnce.Do(func() {
		instance = initiateConnection()
	})

	return instance
}

func initiateConnection() *gorm.DB {
	databaseConfig := configs.Instance().GetDatabase()
	dbDSN := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=UTC",
		databaseConfig.GetHost(),
		databaseConfig.GetUsername(),
		databaseConfig.GetPassword(),
		databaseConfig.GetName(),
		databaseConfig.GetPort(),
		"disable",
	)

	log.Info().Msgf("dbDSN is %s", dbDSN)
	sqlDB, err := sql.Open("postgres", dbDSN)
	if err != nil {
		log.
			Panic().
			Err(err).
			Msg("failed to retrieve pooling instance")
	}

	poolConfig := databaseConfig.GetPool()
	// configuring the connection pooling specific parameters
	sqlDB.SetMaxIdleConns(poolConfig.GetMaxIdleConnections())
	sqlDB.SetMaxOpenConns(poolConfig.GetMaxOpenConnections())
	sqlDB.SetConnMaxLifetime(poolConfig.GetConnectionMaxLifetime())
	sqlDB.SetConnMaxIdleTime(poolConfig.GetConnectionMaxIdleTime())

	connection, err := gorm.Open(
		postgres.New(
			postgres.Config{Conn: sqlDB},
		),
		&gorm.Config{
			Logger: logger.Default.LogMode(databaseConfig.GetLogLevel()),
		},
	)

	if err != nil {
		log.
			Panic().
			Err(err).
			Msg("db initialization failed")
	}
	return connection
}
