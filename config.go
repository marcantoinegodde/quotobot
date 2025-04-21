package main

import (
	"github.com/spf13/viper"
)

type Config struct {
	Token  string `mapstructure:"token"`
	ChatID int64  `mapstructure:"chat_id"`
}

func loadConfig(logger *Logger) *Config {
	var config Config

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/quotobot/")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		logger.Error.Fatalln("Fatal error reading config file:", err)
	}

	if err := viper.Unmarshal(&config); err != nil {
		logger.Error.Fatalln("Fatal error unmarshal config:", err)
	}

	return &config
}
