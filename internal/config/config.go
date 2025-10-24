package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// Config represents the application configuration
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Auth     AuthConfig     `mapstructure:"auth"`
	NATS     NATSConfig     `mapstructure:"nats"`
	Logger   LoggerConfig   `mapstructure:"logger"`
}

// ServerConfig contains server-related configuration
type ServerConfig struct {
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	Mode            string        `mapstructure:"mode"` // debug, release
	ReadTimeout     time.Duration `mapstructure:"read_timeout"`
	WriteTimeout    time.Duration `mapstructure:"write_timeout"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"`
}

// DatabaseConfig contains database-related configuration
type DatabaseConfig struct {
	URI            string        `mapstructure:"uri"`
	Database       string        `mapstructure:"database"`
	MaxPoolSize    uint64        `mapstructure:"max_pool_size"`
	MinPoolSize    uint64        `mapstructure:"min_pool_size"`
	ConnectTimeout time.Duration `mapstructure:"connect_timeout"`
}

// AuthConfig contains authentication-related configuration
type AuthConfig struct {
	JWTSecret         string        `mapstructure:"jwt_secret"`
	JWTExpiration     time.Duration `mapstructure:"jwt_expiration"`
	RefreshExpiration time.Duration `mapstructure:"refresh_expiration"`
}

// NATSConfig contains NATS-related configuration
type NATSConfig struct {
	URL     string `mapstructure:"url"`
	Enabled bool   `mapstructure:"enabled"`
}

// LoggerConfig contains logging-related configuration
type LoggerConfig struct {
	Level       string `mapstructure:"level"`       // debug, info, warn, error
	Development bool   `mapstructure:"development"` // true for development mode
	Encoding    string `mapstructure:"encoding"`    // json, console
}

// Load loads configuration from file and environment variables
func Load(configPath string) (*Config, error) {
	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")

	// Enable reading from environment variables
	viper.AutomaticEnv()

	// Set default values
	setDefaults()

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Validate configuration
	if err := validate(&config); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &config, nil
}

// setDefaults sets default configuration values
func setDefaults() {
	// Server defaults
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.mode", "debug")
	viper.SetDefault("server.read_timeout", 10)
	viper.SetDefault("server.write_timeout", 10)
	viper.SetDefault("server.shutdown_timeout", 30)

	// Database defaults
	viper.SetDefault("database.uri", "mongodb://localhost:27017")
	viper.SetDefault("database.database", "nutrient_db")
	viper.SetDefault("database.max_pool_size", 100)
	viper.SetDefault("database.min_pool_size", 10)
	viper.SetDefault("database.connect_timeout", 10)

	// Auth defaults
	viper.SetDefault("auth.jwt_expiration", 3600)
	viper.SetDefault("auth.refresh_expiration", 604800)

	// NATS defaults
	viper.SetDefault("nats.url", "nats://localhost:4222")
	viper.SetDefault("nats.enabled", false)

	// Logger defaults
	viper.SetDefault("logger.level", "debug")
	viper.SetDefault("logger.development", true)
	viper.SetDefault("logger.encoding", "console")
}

// validate validates the configuration
func validate(config *Config) error {
	if config.Server.Port <= 0 || config.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", config.Server.Port)
	}

	if config.Server.Mode != "debug" && config.Server.Mode != "release" {
		return fmt.Errorf("invalid server mode: %s", config.Server.Mode)
	}

	if config.Database.URI == "" {
		return fmt.Errorf("database URI is required")
	}

	if config.Database.Database == "" {
		return fmt.Errorf("database name is required")
	}

	if config.Auth.JWTSecret == "" {
		return fmt.Errorf("JWT secret is required")
	}

	if config.Logger.Level != "debug" && config.Logger.Level != "info" &&
		config.Logger.Level != "warn" && config.Logger.Level != "error" {
		return fmt.Errorf("invalid logger level: %s", config.Logger.Level)
	}

	if config.Logger.Encoding != "json" && config.Logger.Encoding != "console" {
		return fmt.Errorf("invalid logger encoding: %s", config.Logger.Encoding)
	}

	return nil
}
