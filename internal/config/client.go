// Package config provides configuration structures and loading utilities for the application.
package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

// ClientConfig holds configuration parameters specific to the TLN client.
type ClientConfig struct {
	LogFormat string `env:"TLN_LOG_FORMAT" envDefault:"text"`
	LogLevel  string `env:"TLN_LOG_LEVEL"  envDefault:"info"`
}

// LoadClient parses environment variables and returns a populated ClientConfig.
func LoadClient() (ClientConfig, error) {
	var c ClientConfig
	if err := env.Parse(&c); err != nil {
		return c, fmt.Errorf("parse client config: %w", err)
	}

	return c, nil
}
