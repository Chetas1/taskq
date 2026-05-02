// Package config loads runtime configuration from `config/config.yaml`
// and environment variables (via Viper's AutomaticEnv). Environment variables
// override file values; nested keys use `_` separators (e.g. `REDIS_ADDR`).
package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Config is the top-level taskq configuration.
type Config struct {
	Redis Redis
}

// Redis holds the connection details for the backing Redis instance / cluster.
type Redis struct {
	Addr     string
	Password string
	DB       int
}

// LoadConfig reads `config/config.yaml`, applies environment overrides, and
// returns the parsed Config. Returns a wrapped error if the file is missing
// or the YAML cannot be unmarshalled.
func LoadConfig() (*Config, error) {
	v := viper.New()

	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("./config/")
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	return &cfg, nil
}
