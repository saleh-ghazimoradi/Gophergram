package config

import (
	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
	"log"
	"time"
)

var AppConfig *Config

type Config struct {
	ServerConfig ServerConfig
	DBConfig     DBConfig
	Context      Context
	Pagination   Pagination
}

type ServerConfig struct {
	Port         string        `env:"SERVER_PORT,required"`
	Version      string        `env:"SERVER_VERSION,required"`
	IdleTimeout  time.Duration `env:"SERVER_IDLE_TIMEOUT,required"`
	ReadTimeout  time.Duration `env:"SERVER_READ_TIMEOUT,required"`
	WriteTimeout time.Duration `env:"SERVER_WRITE_TIMEOUT,required"`
	Env          string        `env:"SERVER_ENV,required"`
}

type Context struct {
	ContextTimeout time.Duration `env:"CONTEXT_TIME_OUT,required"`
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
	
	AppConfig = config

	return nil
}
