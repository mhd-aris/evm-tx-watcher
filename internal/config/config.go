package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Port     string `mapstructure:"APP_PORT"`
	LogLevel string `mapstructure:"LOG_LEVEL"`
}

// Load reads configuration from environment variables and .env file
func Load() *Config {

	// Default values
	viper.SetDefault("APP_PORT", "8080")
	viper.SetDefault("LOG_LEVEL", "info")

	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Warning: Could not read config file: %v", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Error unmarshalling config: %v", err)
	}

	return &config
}
