package config

import (
	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
	"log"
)

var AppConfig *Config

type Config struct {
	ServerConfig ServerConfig
	DBConfig     DBConfig
}

type ServerConfig struct {
	Port    string `env:"SERVER_PORT,required"`
	Version string `env:"SERVER_VERSION,required"`
}

type DBConfig struct {
	DbHost     string `env:"DB_HOST,required"`
	DbPort     string `env:"DB_PORT,required"`
	DbUser     string `env:"DB_USER,required"`
	DbPassword string `env:"DB_PASSWORD,required"`
	DbName     string `env:"DB_NAME,required"`
	DbSslMode  string `env:"DB_SSLMODE,required"`
}

func LoadingConfig() error {
	if err := godotenv.Load(); err != nil {
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

	config.DBConfig = *dbConfig

	AppConfig = config

	return nil
}
