package database

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	pg "github.com/lib/pq"
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

var pgDriver = pg.Driver{}

type customDriver struct {
}

func (d *customDriver) Open(name string) (driver.Conn, error) {
	// Open a new connection
	conn, err := pgDriver.Open(name)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func Instance() *gorm.DB {
	connectionOnce.Do(func() {
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

		sql.Register("custom_postgres", &customDriver{})

		log.Info().Msgf("dbDSN is %s", dbDSN)
		sqlDB, err := sql.Open("custom_postgres", dbDSN)
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

		instance = connection
	})

	return instance
}
