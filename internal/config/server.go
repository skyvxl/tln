// Package config provides structures and functions for loading configuration parameters from environment variables.
package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

// ServerConfig holds configuration parameters specific to the TLN server.
type ServerConfig struct {
	LogFormat   string `env:"TLN_LOG_FORMAT" envDefault:"text"`
	LogLevel    string `env:"TLN_LOG_LEVEL"  envDefault:"info"`
	ControlAddr string `env:"TLN_CONTROL_ADDR" envDefault:"127.0.0.1:7000"`
	TLSCertFile string `env:"TLN_TLS_CERT_FILE" envDefault:"certs/server.crt"`
	TLSKeyFile  string `env:"TLN_TLS_KEY_FILE"  envDefault:"certs/server.key"`
}

// LoadServer parses environment variables and returns a populated ServerConfig.
func LoadServer() (ServerConfig, error) {
	var c ServerConfig
	if err := env.Parse(&c); err != nil {
		return c, fmt.Errorf("parse server config: %w", err)
	}

	return c, nil
}
