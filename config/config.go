package config

import (
	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
	"log"
	"time"
)

var AppConfig *Config

type Config struct {
	ServerConfig   ServerConfig
	DBConfig       DBConfig
	Context        Context
	Pagination     Pagination
	Mail           Mail
	Authentication Authentication
	Redis          Redis
	Rate           Rate
}

type ServerConfig struct {
	Port         string        `env:"SERVER_PORT,required"`
	Version      string        `env:"SERVER_VERSION,required"`
	IdleTimeout  time.Duration `env:"SERVER_IDLE_TIMEOUT,required"`
	ReadTimeout  time.Duration `env:"SERVER_READ_TIMEOUT,required"`
	WriteTimeout time.Duration `env:"SERVER_WRITE_TIMEOUT,required"`
	Env          string        `env:"SERVER_ENV,required"`
	APIURL       string        `env:"SERVER_API_URL,required"`
}

type Context struct {
	ContextTimeout time.Duration `env:"CONTEXT_TIME_OUT,required"`
}

type Authentication struct {
	Username string        `env:"USERNAME,required"`
	Password string        `env:"PASSWORD,required"`
	Secret   string        `env:"SECRET,required"`
	Exp      time.Duration `env:"EXP,required"`
	Aud      string        `env:"AUD,required"`
	Iss      string        `env:"ISS,required"`
}

type Rate struct {
	Limit  int           `env:"RATE_LIMIT,required"`
	Window time.Duration `env:"RATE_WINDOW,required"`
}

type Redis struct {
	Addr    string `env:"REDIS_ADDR,required"`
	PW      string `env:"REDIS_PASSWORD,required"`
	DB      int    `env:"REDIS_DB,required"`
	Enabled bool   `env:"REDIS_ENABLED,required"`
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

type DBConfig struct {
	DBDriver     string        `env:"DB_DRIVER,required"`
	DBSource     string        `env:"DB_SOURCE,required"`
	DbHost       string        `env:"DB_HOST,required"`
	DbPort       string        `env:"DB_PORT,required"`
	DbUser       string        `env:"DB_USER,required"`
	DbPassword   string        `env:"DB_PASSWORD,required"`
	DbName       string        `env:"DB_NAME,required"`
	DbSslMode    string        `env:"DB_SSLMODE,required"`
	MaxOpenConns int           `env:"DB_MAX_OPEN_CONNECTIONS,required"`
	MaxIdleConns int           `env:"DB_MAX_IDLE_CONNECTIONS,required"`
	MaxIdleTime  time.Duration `env:"DB_MAX_IDLE_TIME,required"`
	Timeout      time.Duration `env:"DB_TIMEOUT,required"`
}

type Pagination struct {
	Limit  int    `env:"LIMIT,required"`
	Offset int    `env:"OFFSET,required"`
	Sort   string `env:"SORT,required"`
}

func LoadingConfig() error {
	if err := godotenv.Load("app.env"); err != nil {
		log.Fatal("error loading app.env file")
	}

	config := &Config{}

	if err := env.Parse(config); err != nil {
		log.Fatal("error parsing config")
	}

	serverConfig := &ServerConfig{}

	if err := env.Parse(serverConfig); err != nil {
		log.Fatal("error parsing config")
	}

	config.ServerConfig = *serverConfig

	dbConfig := &DBConfig{}
	if err := env.Parse(dbConfig); err != nil {
		log.Fatal("error parsing config")
	}

	contextConfig := &Context{}

	if err := env.Parse(contextConfig); err != nil {
		log.Fatal("error parsing config")
	}

	config.Context = *contextConfig

	config.DBConfig = *dbConfig

	paginationConfig := &Pagination{}

	if err := env.Parse(paginationConfig); err != nil {
		log.Fatal("error parsing pagination config")
	}

	config.Pagination = *paginationConfig

	mailConfig := &Mail{}
	if err := env.Parse(mailConfig); err != nil {
		log.Fatal("error parsing token config")
	}
	config.Mail = *mailConfig

	authConfig := &Authentication{}
	if err := env.Parse(authConfig); err != nil {
		log.Fatal("error parsing token config")
	}
	config.Authentication = *authConfig

	redisConfig := &Redis{}
	if err := env.Parse(redisConfig); err != nil {
		log.Fatal("error parsing token config")
	}
	config.Redis = *redisConfig

	rateConfig := &Rate{}
	if err := env.Parse(rateConfig); err != nil {
		log.Fatal("error parsing token config")
	}

	config.Rate = *rateConfig

	AppConfig = config

	return nil
}
