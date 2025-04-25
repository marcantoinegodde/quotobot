package config

import (
	"quotobot/pkg/logger"

	"github.com/spf13/viper"
)

type Config struct {
	Bot    Bot    `mapstructure:"bot"`
	Server Server `mapstructure:"server"`
}

type Bot struct {
	Token      string `mapstructure:"token"`
	ChatID     int64  `mapstructure:"chat_id"`
	BaseURL    string `mapstructure:"base_url"`
	HMACSecret string `mapstructure:"hmac_secret"`
}

type Server struct {
	SessionSecret string `mapstructure:"session_secret"`
	HMACSecret    string `mapstructure:"hmac_secret"`
	ProviderURL   string `mapstructure:"provider_url"`
	ClientID      string `mapstructure:"client_id"`
	ClientSecret  string `mapstructure:"client_secret"`
	RedirectURL   string `mapstructure:"redirect_url"`
}

func LoadConfig(logger *logger.Logger) *Config {
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
