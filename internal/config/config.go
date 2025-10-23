package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// Config holds all application configuration
type Config struct {
	Server     ServerConfig
	Database   DatabaseConfig
	Schoology  SchoologyConfig
	PowerSchool PowerSchoolConfig
	Sync       SyncConfig
	Cache      CacheConfig
	RateLimits RateLimitsConfig
	Logging    LoggingConfig
}

// ServerConfig holds HTTP server configuration
type ServerConfig struct {
	Port         int
	Env          string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Path           string
	MaxOpenConns   int
	MaxIdleConns   int
	EncryptionKey  string
}

// SchoologyConfig holds Schoology API configuration
type SchoologyConfig struct {
	ConsumerKey    string
	ConsumerSecret string
	BaseURL        string
}

// PowerSchoolConfig holds PowerSchool API configuration
type PowerSchoolConfig struct {
	ClientID     string
	ClientSecret string
	BaseURL      string
}

// SyncConfig holds data synchronization configuration
type SyncConfig struct {
	Interval       time.Duration
	RetryAttempts  int
	RetryBackoff   string
}

// CacheConfig holds cache configuration
type CacheConfig struct {
	Enabled         bool
	TTLAssignments  time.Duration
	TTLGrades       time.Duration
	TTLEvents       time.Duration
	TTLMessages     time.Duration
}

// RateLimitsConfig holds rate limiting configuration
type RateLimitsConfig struct {
	Schoology    int // requests per minute
	PowerSchool  int // requests per minute
}

// LoggingConfig holds logging configuration
type LoggingConfig struct {
	Level  string
	Format string // "json" or "text"
}

// Load reads configuration from environment and config file
func Load() (*Config, error) {
	v := viper.New()

	// Set defaults
	setDefaults(v)

	// Read from environment variables
	v.AutomaticEnv()
	v.SetEnvPrefix("TRUNCHBULL")

	// Try to read config file
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("./config")
	v.AddConfigPath("/config")
	v.AddConfigPath(".")

	if err := v.ReadInConfig(); err != nil {
		// Config file not required, environment variables take precedence
		fmt.Printf("Warning: config file not found, using environment variables\n")
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Validate configuration
	if err := validate(&cfg); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &cfg, nil
}

func setDefaults(v *viper.Viper) {
	// Server defaults
	v.SetDefault("server.port", 8080)
	v.SetDefault("server.env", "development")
	v.SetDefault("server.read_timeout", "10s")
	v.SetDefault("server.write_timeout", "10s")

	// Database defaults
	v.SetDefault("database.path", "/data/trunchbull.db")
	v.SetDefault("database.max_open_conns", 10)
	v.SetDefault("database.max_idle_conns", 5)

	// Schoology defaults
	v.SetDefault("schoology.base_url", "https://api.schoology.com/v1")

	// PowerSchool defaults - will be set by district
	v.SetDefault("powerschool.base_url", "https://powerschool.seattleschools.org")

	// Sync defaults
	v.SetDefault("sync.interval", "30m")
	v.SetDefault("sync.retry_attempts", 3)
	v.SetDefault("sync.retry_backoff", "exponential")

	// Cache defaults
	v.SetDefault("cache.enabled", true)
	v.SetDefault("cache.ttl_assignments", "15m")
	v.SetDefault("cache.ttl_grades", "30m")
	v.SetDefault("cache.ttl_events", "24h")
	v.SetDefault("cache.ttl_messages", "10m")

	// Rate limits
	v.SetDefault("rate_limits.schoology", 60)
	v.SetDefault("rate_limits.powerschool", 30)

	// Logging defaults
	v.SetDefault("logging.level", "info")
	v.SetDefault("logging.format", "text")
}

func validate(cfg *Config) error {
	// Validate required fields
	if cfg.Server.Port <= 0 || cfg.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", cfg.Server.Port)
	}

	if cfg.Database.Path == "" {
		return fmt.Errorf("database path is required")
	}

	// Note: API credentials can be empty initially, will be set via OAuth flow
	// but warn if they're missing
	if cfg.Schoology.ConsumerKey == "" {
		fmt.Println("Warning: Schoology consumer key not set")
	}
	if cfg.PowerSchool.ClientID == "" {
		fmt.Println("Warning: PowerSchool client ID not set")
	}

	return nil
}
