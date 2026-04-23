package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type ServerConfig struct {
	LogFormat string `env:"TLN_LOG_FORMAT" envDefault:"text"`
	LogLevel  string `env:"TLN_LOG_LEVEL"  envDefault:"info"`
}

func LoadServer() (ServerConfig, error) {
	var c ServerConfig
	if err := env.Parse(&c); err != nil {
		return c, fmt.Errorf("parse server config: %w", err)
	}

	return c, nil
}
