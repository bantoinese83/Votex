package config

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/spf13/viper"
)

type Environment string

const (
	Development Environment = "development"
	Staging     Environment = "staging"
	Production  Environment = "production"
)

type Config struct {
	Environment Environment `mapstructure:"ENVIRONMENT"`
	Port        string      `mapstructure:"PORT"`
	DBURL       string      `mapstructure:"DB_URL"`
	RedisURL    string      `mapstructure:"REDIS_URL"`
	JWTSecret   string      `mapstructure:"JWT_SECRET"`
	LogLevel    string      `mapstructure:"LOG_LEVEL"`
	CORSOrigins []string    `mapstructure:"CORS_ORIGINS"`
}

func Load() *Config {
	// Set default config file
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	// Enable environment variable override
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Read config file if it exists
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			slog.Warn("could not read config file", "error", err)
		}
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		slog.Error("could not unmarshal config", "error", err)
		panic(err)
	}

	// Set defaults
	setDefaults(&cfg)

	// Validate configuration
	if err := validateConfig(&cfg); err != nil {
		slog.Error("invalid configuration", "error", err)
		panic(err)
	}

	return &cfg
}

func setDefaults(cfg *Config) {
	if cfg.Environment == "" {
		cfg.Environment = Development
	}
	if cfg.Port == "" {
		cfg.Port = "8080"
	}
	if cfg.DBURL == "" {
		cfg.DBURL = "postgres://user:password@localhost:5432/vortexdb?sslmode=disable"
	}
	if cfg.RedisURL == "" {
		cfg.RedisURL = "redis://localhost:6379"
	}
	if cfg.JWTSecret == "" {
		cfg.JWTSecret = "a-very-secret-key-change-in-production"
	}
	if cfg.LogLevel == "" {
		cfg.LogLevel = "info"
	}
	if len(cfg.CORSOrigins) == 0 {
		cfg.CORSOrigins = []string{"http://localhost:5173", "http://localhost:3000"}
	}
}

func validateConfig(cfg *Config) error {
	if cfg.JWTSecret == "a-very-secret-key-change-in-production" && cfg.Environment == Production {
		return fmt.Errorf("JWT_SECRET must be set in production environment")
	}

	if cfg.DBURL == "" {
		return fmt.Errorf("DB_URL is required")
	}

	return nil
}

func (c *Config) IsDevelopment() bool {
	return c.Environment == Development
}

func (c *Config) IsProduction() bool {
	return c.Environment == Production
}

func (c *Config) IsStaging() bool {
	return c.Environment == Staging
}
