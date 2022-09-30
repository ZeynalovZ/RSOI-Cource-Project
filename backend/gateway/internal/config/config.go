package config

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

type (
	Config struct {
		HTTP     HTTPConfig
		Services ServicesConfig
	}

	HTTPConfig struct {
		Host string `mapstructure:"host"`
		Port string `mapstructure:"port"`
	}

	ServicesConfig struct {
		UserService         string `mapstructure:"user_service"`
		MusicService        string `mapstructure:"music_service"`
		NotificationService string `mapstructure:"notification_service"`
		PaymentService      string `mapstructure:"payment_service"`
		SessionService      string `mapstructure:"session_service"`
	}
)

func Init(path string, logger *log.Logger) (*Config, error) {
	if err := parseConfigFile(path); err != nil {
		logger.Printf("failed to parse path to config file: %s", err)
		return nil, err
	}

	var cfg Config
	if err := unmarshal(&cfg, logger); err != nil {
		logger.Printf("failed to unmarshal config: %s", err)
		return nil, err
	}

	return &cfg, nil
}

func unmarshal(cfg *Config, logger *log.Logger) error {
	if err := viper.UnmarshalKey("http", &cfg.HTTP); err != nil {
		logger.Printf("failed to unmarshal http key in config: %s", err)
		return err
	}

	if err := viper.UnmarshalKey("services", &cfg.Services); err != nil {
		logger.Printf("failed to unmarshal http key in config: %s", err)
		return err
	}

	return nil
}

func parseConfigFile(filepath string) error {
	path := strings.Split(filepath, "/")

	viper.AddConfigPath(path[0])
	viper.SetConfigName(path[1])

	return viper.ReadInConfig()
}
