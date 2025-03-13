package config

import (
	"github.com/caarlos0/env/v11"
	"github.com/saleh-ghazimoradi/Gophergram/sLogger"
	"time"
)

var AppConfig *Config

type Config struct {
	ServerConfig   ServerConfig
	Database       Database
	Authentication Authentication
}

type Database struct {
	DatabaseHost     string        `env:"DATABASE_HOST,required"`
	DatabasePort     string        `env:"DATABASE_PORT,required"`
	DatabaseUser     string        `env:"DATABASE_USER,required"`
	DatabasePassword string        `env:"DATABASE_PASSWORD,required"`
	DatabaseName     string        `env:"DATABASE_NAME,required"`
	DatabaseSSLMode  string        `env:"DATABASE_SSLMODE,required"`
	MaxOpenConn      int           `env:"DB_MAX_OPEN_CONNECTIONS,required"`
	MaxIdleConn      int           `env:"DB_MAX_IDLE_CONNECTIONS,required"`
	MaxIdleTime      time.Duration `env:"DB_MAX_IDLE_TIME,required"`
	Timeout          time.Duration `env:"DB_TIMEOUT,required"`
}

type ServerConfig struct {
	Port    string `env:"SERVER_PORT"`
	Version string `env:"SERVER_VERSION"`
	Env     string `env:"SERVER_ENV"`
}

type Authentication struct {
	Secret string `env:"AUTHENTICATION_SECRET"`
}

func LoadConfig() error {
	config := &Config{}

	if err := env.Parse(config); err != nil {
		sLogger.SLogger.Error("error loading config", "error", err)
	}

	AppConfig = config
	return nil
}
