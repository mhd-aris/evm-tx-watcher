package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	AppPort  string         `mapstructure:"APP_PORT"`
	LogLevel string         `mapstructure:"LOG_LEVEL"`
	DB       DatabaseConfig `mapstructure:",squash"`
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host     string `mapstructure:"DB_HOST"`
	Port     int    `mapstructure:"DB_PORT"`
	User     string `mapstructure:"DB_USER"`
	Password string `mapstructure:"DB_PASSWORD"`
	Name     string `mapstructure:"DB_NAME"`
}

// Load reads configuration from environment variables and .env file
func Load() *Config {
	// Set defaults
	viper.SetDefault("APP_PORT", "8080")
	viper.SetDefault("LOG_LEVEL", "info")
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", 5432)

	// Read from .env file
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	// Try to read config file, but don't fail if it doesn't exist
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Warning: Could not read config file: %v", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Error unmarshalling config: %v", err)
	}

	// Validate required fields
	if err := validateConfig(&config); err != nil {
		log.Fatalf("Config validation failed: %v", err)
	}

	return &config
}

// validateConfig validates the loaded configuration
func validateConfig(cfg *Config) error {
	if cfg.AppPort == "" {
		return fmt.Errorf("APP_PORT is required")
	}
	if cfg.LogLevel == "" {
		return fmt.Errorf("LOG_LEVEL is required")
	}
	if cfg.DB.Host == "" {
		return fmt.Errorf("DB_HOST is required")
	}
	return nil
}
