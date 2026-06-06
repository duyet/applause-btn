package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Config holds all application configuration
type Config struct {
	// Server configuration
	Port string

	// Database configuration
	DBLocation string

	// Header configuration for clapper info
	HeaderUserEmail string
	HeaderUserID    string

	// CORS configuration
	AllowedOrigins []string

	// Rate limiting
	RateLimitEnabled bool
	RateLimitPerIP   int // requests per minute
	RateLimitPerURL  int // updates per URL per minute

	// Input limits
	MaxURLsPerRequest int
	MaxClapsPerUpdate int
}

// Load loads configuration from environment variables with sensible defaults
func Load() (*Config, error) {
	cfg := &Config{
		Port:              getEnv("PORT", "3000"),
		DBLocation:        getEnv("DB_LOCATION", "/tmp/badger"),
		HeaderUserEmail:   getEnv("HEADER_USER_EMAIL", "x-authenticated-user-email"),
		HeaderUserID:      getEnv("HEADER_USER_ID", "x-authenticated-uid"),
		AllowedOrigins:    getEnvList("ALLOWED_ORIGINS", []string{"*"}),
		RateLimitEnabled:  getEnvBool("RATE_LIMIT_ENABLED", true),
		RateLimitPerIP:    getEnvInt("RATE_LIMIT_PER_IP", 100),
		RateLimitPerURL:   getEnvInt("RATE_LIMIT_PER_URL", 10),
		MaxURLsPerRequest: getEnvInt("MAX_URLS_PER_REQUEST", 100),
		MaxClapsPerUpdate: getEnvInt("MAX_CLAPS_PER_UPDATE", 10),
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return cfg, nil
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.Port == "" {
		return fmt.Errorf("PORT cannot be empty")
	}

	if c.DBLocation == "" {
		return fmt.Errorf("DB_LOCATION cannot be empty")
	}

	if c.MaxURLsPerRequest <= 0 {
		return fmt.Errorf("MAX_URLS_PER_REQUEST must be positive")
	}

	if c.MaxClapsPerUpdate <= 0 || c.MaxClapsPerUpdate > 50 {
		return fmt.Errorf("MAX_CLAPS_PER_UPDATE must be between 1 and 50")
	}

	return nil
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultVal
}

// getEnvInt gets an environment variable as int or returns a default value
func getEnvInt(key string, defaultVal int) int {
	if val, ok := os.LookupEnv(key); ok {
		if intVal, err := strconv.Atoi(val); err == nil {
			return intVal
		}
	}
	return defaultVal
}

// getEnvBool gets an environment variable as bool or returns a default value
func getEnvBool(key string, defaultVal bool) bool {
	if val, ok := os.LookupEnv(key); ok {
		if boolVal, err := strconv.ParseBool(val); err == nil {
			return boolVal
		}
	}
	return defaultVal
}

// getEnvList gets an environment variable as a comma-separated list or returns a default value
func getEnvList(key string, defaultVal []string) []string {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		parts := strings.Split(val, ",")
		result := make([]string, 0, len(parts))
		for _, part := range parts {
			if trimmed := strings.TrimSpace(part); trimmed != "" {
				result = append(result, trimmed)
			}
		}
		if len(result) > 0 {
			return result
		}
	}
	return defaultVal
}
