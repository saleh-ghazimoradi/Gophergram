package config

import (
	"github.com/caarlos0/env/v11"
	"github.com/saleh-ghazimoradi/Gophergram/sLogger"
)

var AppConfig *Config

type Config struct {
	ServerConfig ServerConfig
}

type ServerConfig struct {
	Port    string `env:"SERVER_PORT"`
	Version string `env:"SERVER_VERSION"`
	Env     string `env:"SERVER_ENV"`
}

func LoadConfig() error {
	config := &Config{}

	if err := env.Parse(config); err != nil {
		sLogger.SLogger.Error("error loading config", "error", err)
	}

	AppConfig = config
	return nil
}
