package config

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/spf13/viper"
)

type Environment string
type DatabaseType string

const (
	Development Environment = "development"
	Staging     Environment = "staging"
	Production  Environment = "production"
)

const (
	PostgreSQL DatabaseType = "postgres"
	SQLite     DatabaseType = "sqlite"
)

type Config struct {
	Environment Environment  `mapstructure:"ENVIRONMENT"`
	Port        string       `mapstructure:"PORT"`
	DBURL       string       `mapstructure:"DB_URL"`
	DBType      DatabaseType `mapstructure:"DB_TYPE"`
	SQLitePath  string       `mapstructure:"SQLITE_PATH"`
	RedisURL    string       `mapstructure:"REDIS_URL"`
	JWTSecret   string       `mapstructure:"JWT_SECRET"`
	LogLevel    string       `mapstructure:"LOG_LEVEL"`
	CORSOrigins []string     `mapstructure:"CORS_ORIGINS"`

	// Email configuration
	SMTPHost     string `mapstructure:"SMTP_HOST"`
	SMTPPort     int    `mapstructure:"SMTP_PORT"`
	SMTPUsername string `mapstructure:"SMTP_USERNAME"`
	SMTPPassword string `mapstructure:"SMTP_PASSWORD"`
	SMTPFrom     string `mapstructure:"SMTP_FROM"`
	SMTPTLS      bool   `mapstructure:"SMTP_TLS"`

	// Password reset configuration
	PasswordResetTokenExpiry int    `mapstructure:"PASSWORD_RESET_TOKEN_EXPIRY"` // in hours
	AppURL                   string `mapstructure:"APP_URL"`

	// Rate limiting
	RateLimitRequests int `mapstructure:"RATE_LIMIT_REQUESTS"` // requests per minute
	RateLimitBurst    int `mapstructure:"RATE_LIMIT_BURST"`    // burst size
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
	if cfg.DBType == "" {
		cfg.DBType = PostgreSQL
	}
	if cfg.SQLitePath == "" {
		cfg.SQLitePath = "./data/votex.db"
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

	// Email defaults
	if cfg.SMTPHost == "" {
		cfg.SMTPHost = "localhost"
	}
	if cfg.SMTPPort == 0 {
		cfg.SMTPPort = 587
	}
	if cfg.SMTPFrom == "" {
		cfg.SMTPFrom = "noreply@vortex.com"
	}
	if !cfg.SMTPTLS {
		cfg.SMTPTLS = true // Default to TLS for security
	}

	// Password reset defaults
	if cfg.PasswordResetTokenExpiry == 0 {
		cfg.PasswordResetTokenExpiry = 24 // 24 hours
	}
	if cfg.AppURL == "" {
		cfg.AppURL = "http://localhost:5173"
	}

	// Rate limiting defaults
	if cfg.RateLimitRequests == 0 {
		cfg.RateLimitRequests = 100 // 100 requests per minute
	}
	if cfg.RateLimitBurst == 0 {
		cfg.RateLimitBurst = 20 // burst of 20 requests
	}
}

func validateConfig(cfg *Config) error {
	if cfg.JWTSecret == "a-very-secret-key-change-in-production" && cfg.Environment == Production {
		return fmt.Errorf("JWT_SECRET must be set in production environment")
	}

	if cfg.DBURL == "" && cfg.DBType == PostgreSQL {
		return fmt.Errorf("DB_URL is required for PostgreSQL")
	}

	if cfg.SQLitePath == "" && cfg.DBType == SQLite {
		return fmt.Errorf("SQLITE_PATH is required for SQLite")
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

func (c *Config) IsSQLite() bool {
	return c.DBType == SQLite
}

func (c *Config) IsPostgreSQL() bool {
	return c.DBType == PostgreSQL
}
