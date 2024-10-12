package config

import (
	"fmt"
	"github.com/spf13/viper"
)

var AppConfig *config

type config struct {
	General General `mapstructure:"general"`
}

type General struct {
	Listen   string `mapstructure:"listen"`
	LogLevel int8   `mapstructure:"log_level"`
}

type Database struct{}

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
