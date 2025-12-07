package config

import (
	"fmt"
	"log"
	"strings"

	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	AppPort   string                   `mapstructure:"APP_PORT"`
	LogLevel  string                   `mapstructure:"LOG_LEVEL"`
	LogFormat string                   `mapstructure:"LOG_FORMAT"`
	DB        DatabaseConfig           `mapstructure:",squash"`
	Redis     RedisConfig              `mapstructure:",squash"`
	Networks  map[string]NetworkConfig `mapstructure:"-"`
}

// NetworkConfig holds network configuration with chain ID
type NetworkConfig struct {
	Name    string
	ChainID int64
	RPC     string
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host     string `mapstructure:"DB_HOST"`
	Port     int    `mapstructure:"DB_PORT"`
	User     string `mapstructure:"DB_USER"`
	Password string `mapstructure:"DB_PASSWORD"`
	Name     string `mapstructure:"DB_NAME"`
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Host     string `mapstructure:"REDIS_HOST"`
	Port     int    `mapstructure:"REDIS_PORT"`
	Password string `mapstructure:"REDIS_PASSWORD"`
	DB       int    `mapstructure:"REDIS_DB"`
}

func Load() (*Config, error) {
	viper.SetDefault("APP_PORT", "8080")
	viper.SetDefault("LOG_LEVEL", "info")
	viper.SetDefault("LOG_FORMAT", "text")
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", 5432)
	viper.SetDefault("REDIS_HOST", "localhost")
	viper.SetDefault("REDIS_PORT", 6379)
	viper.SetDefault("REDIS_DB", 0)

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

	// Setup predefined testnet networks
	config.Networks = getTestnetNetworks()

	// Override RPC URLs from environment if provided
	for name, network := range config.Networks {
		key := fmt.Sprintf("RPC_%s", strings.ToUpper(strings.ReplaceAll(name, "-", "_")))
		if url := viper.GetString(key); url != "" {
			network.RPC = url
			config.Networks[name] = network
		}
	}

	// Validate that all networks have RPC URLs
	for name, network := range config.Networks {
		if network.RPC == "" {
			return nil, fmt.Errorf("RPC URL for network %s is not set", name)
		}
	}

	return &config, nil
}

// getTestnetNetworks returns hardcoded testnet network configurations
func getTestnetNetworks() map[string]NetworkConfig {
	return map[string]NetworkConfig{
		"ethereum-sepolia": {
			Name:    "ethereum-sepolia",
			ChainID: 11155111,
			RPC:     "", // will be set from env
		},
		"base-sepolia": {
			Name:    "base-sepolia",
			ChainID: 84532,
			RPC:     "", // will be set from env
		},
		"arbitrum-sepolia": {
			Name:    "arbitrum-sepolia",
			ChainID: 421614,
			RPC:     "", // will be set from env
		},
	}
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
