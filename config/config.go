package config

import (
	"fmt"
	"github.com/spf13/viper"
	"time"
)

var AppConfig *config

type config struct {
	General  General  `mapstructure:"general"`  // general configs
	Database Database `mapstructure:"database"` // databases configs
	Env      Env      `mapstructure:"env"`
}

type Env struct {
	Env string `mapstructure:"env"`
}

type General struct {
	Listen   string `mapstructure:"listen"` // rest listen port
	LogLevel int8   `mapstructure:"log_level"`
}

type Database struct {
	Postgresql Postgresql `mapstructure:"postgresql"`
}

type Postgresql struct {
	Host         string        `mapstructure:"host"`           // postgres host
	Port         string        `mapstructure:"port"`           // postgres port
	User         string        `mapstructure:"user"`           // postgres user
	Password     string        `mapstructure:"password"`       // postgres pass
	Database     string        `mapstructure:"database"`       // postgres database // postgres test database
	SSLMode      string        `mapstructure:"ssl_mode"`       // postgres ssl mode
	MaxOpenConns int           `mapstructure:"max_open_conns"` // postgres max open connections
	MaxIdleConns int           `mapstructure:"max_idle_conns"` // postgres max idle connections
	MaxIdleTime  time.Duration `mapstructure:"max_idle_time"`
	Timeout      time.Duration `mapstructure:"timeout"`
}

func LoadConfig(path string) {
	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("json")   // REQUIRED if the config file does not have the extension in the name

	if path == "" {
		viper.AddConfigPath("./app/config")             // path to look for the config file in
		viper.AddConfigPath("./config")                 // path to look for the config file in
		viper.AddConfigPath("$HOME/.config/Gophergram") // call multiple times to add many search paths
		viper.AddConfigPath(".")                        // optionally look for config in the working directory
	} else {
		viper.AddConfigPath(path)
	}

	viper.AutomaticEnv() // read in environment variables that match

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	AppConfig = &config{}
	if err = viper.Unmarshal(&AppConfig); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
}
