// Package config provides configuration structures and loading utilities for the application.
package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

// ClientConfig holds configuration parameters specific to the TLN client.
type ClientConfig struct {
	LogFormat  string `env:"TLN_LOG_FORMAT" envDefault:"text"`
	LogLevel   string `env:"TLN_LOG_LEVEL"  envDefault:"info"`
	ServerAddr string `env:"TLN_SERVER_ADDR" envDefault:"127.0.0.1:7000"`
	ServerName string `env:"TLN_SERVER_NAME" envDefault:"localhost"`
	CAFile     string `env:"TLN_CA_FILE" envDefault:"certs/ca.crt"`
}

// LoadClient parses environment variables and returns a populated ClientConfig.
func LoadClient() (ClientConfig, error) {
	var c ClientConfig
	if err := env.Parse(&c); err != nil {
		return c, fmt.Errorf("parse client config: %w", err)
	}

	return c, nil
}
