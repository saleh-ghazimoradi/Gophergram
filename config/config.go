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
	Mail           Mail
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

type Mail struct {
	Exp                 time.Duration `env:"TOKEN_EXPIRATION,required"`
	FromName            string        `env:"MAIL_FROM_NAME,required"`
	MaxRetries          uint          `env:"MAIL_MAX_RETRIES,required"`
	UserWelcomeTemplate string        `env:"TOKEN_USER_WELCOME_TEMPLATE,required"`
	FromEmail           string        `env:"FROM_EMAIL,required"`
	ApiKey              string        `env:"API_KEY,required"`
	FrontendURL         string        `env:"FRONTEND_URL,required"`
}

type Authentication struct {
	Secret   string `env:"AUTHENTICATION_SECRET"`
	Password string `env:"AUTHENTICATION_PASSWORD"`
	Username string `env:"AUTHENTICATION_USERNAME"`
}

func LoadConfig() error {
	config := &Config{}

	if err := env.Parse(config); err != nil {
		sLogger.SLogger.Error("error loading config", "error", err)
	}

	AppConfig = config
	return nil
}
