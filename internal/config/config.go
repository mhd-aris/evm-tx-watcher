package config

import (
	"fmt"
	"log"
	"strings"

	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	AppPort   string            `mapstructure:"APP_PORT"`
	LogLevel  string            `mapstructure:"LOG_LEVEL"`
	LogFormat string            `mapstructure:"LOG_FORMAT"`
	DB        DatabaseConfig    `mapstructure:",squash"`
	Networks  map[string]string `mapstructure:"-"`
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host     string `mapstructure:"DB_HOST"`
	Port     int    `mapstructure:"DB_PORT"`
	User     string `mapstructure:"DB_USER"`
	Password string `mapstructure:"DB_PASSWORD"`
	Name     string `mapstructure:"DB_NAME"`
}

func Load() (*Config, error) {
	viper.SetDefault("APP_PORT", "8080")
	viper.SetDefault("LOG_LEVEL", "info")
	viper.SetDefault("LOG_FORMAT", "text")
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", 5432)

	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Warning: Could not read config file: %v", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshalling config: %w", err)
	}

	if err := validateConfig(&config); err != nil {
		return nil, err
	}

	// Parse NETWORKS
	networksRaw := viper.GetString("NETWORKS")
	networkNames := strings.Split(networksRaw, ",")
	networkMap := make(map[string]string)

	for _, name := range networkNames {
		name = strings.TrimSpace(name)
		if name == "" {
			continue
		}
		key := fmt.Sprintf("RPC_%s", strings.ToUpper(name))
		url := viper.GetString(key)
		if url == "" {
			return nil, fmt.Errorf("RPC URL for network %s is not set (env: %s)", name, key)
		}
		networkMap[name] = url
	}
	config.Networks = networkMap

	return &config, nil
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
