package configs

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-playground/validator/v10"

	"github.com/spf13/viper"
)

// Load initializes the configuration from .env file and environment variables
// Environment variables take precedence over .env file values
func Load() (*Config, error) {
	cfg := &Config{}
	v := viper.New()
	// Determine the config path
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		// Try to find the project root by looking for go.mod file
		configPath = findProjectRoot()
		if configPath == "" {
			// Fallback to current directory or parent directory
			configPath = "."
			// Check if we're in a subdirectory (like cmd/server)
			if _, err := os.Stat(".env"); os.IsNotExist(err) {
				if _, err := os.Stat("../../.env"); err == nil {
					configPath = "../.."
				} else if _, err := os.Stat("../.env"); err == nil {
					configPath = ".."
				}
			}
		}
	}

	v.AddConfigPath(configPath)
	v.SetConfigType("env")
	v.SetConfigName(".env")

	// Set config type and name for .env file
	env := os.Getenv("APP.ENV")
	if env == "" {
		env = os.Getenv("APP_ENV")
	}

	// For non-local environments, try .env.example as fallback
	if env != "local" && env != "" {
		v.SetConfigName(".env.example")
	}

	// Enable automatic environment variable binding
	v.AutomaticEnv()

	// Try to read .env file (optional for local development)
	if err := v.ReadInConfig(); err != nil {
		// It's okay if .env doesn't exist (production uses env vars)
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return cfg, fmt.Errorf("error reading config file: %w", err)
		}
	}

	// Set default values
	// setDefaults(v)

	// Unmarshal config into struct using mapstructure tags
	cfg = &Config{}
	if err := v.Unmarshal(cfg); err != nil {
		return cfg, fmt.Errorf("unable to decode config: %w", err)
	}

	// Validate required fields
	validate := validator.New()
	if err := validate.Struct(cfg); err != nil {
		return cfg, fmt.Errorf("config validation failed: %w", err)
	}

	return cfg, nil
}

// findProjectRoot attempts to find the project root by looking for go.mod file
func findProjectRoot() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}

	// Start from current directory and walk up to find go.mod
	for {
		goModPath := filepath.Join(dir, "go.mod")
		if _, err := os.Stat(goModPath); err == nil {
			return dir
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			// Reached filesystem root without finding go.mod
			break
		}
		dir = parent
	}

	// Alternative: look for specific files that indicate project root
	dir, _ = os.Getwd()
	for {
		envPath := filepath.Join(dir, ".env")
		makefilePath := filepath.Join(dir, "Makefile")
		if _, err := os.Stat(envPath); err == nil {
			return dir
		}
		if _, err := os.Stat(makefilePath); err == nil {
			return dir
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return ""
}
