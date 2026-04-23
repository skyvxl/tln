package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type ClientConfig struct {
	LogFormat string `env:"TLN_LOG_FORMAT" envDefault:"text"`
	LogLevel  string `env:"TLN_LOG_LEVEL"  envDefault:"info"`
}

func LoadClient() (ClientConfig, error) {
	var c ClientConfig
	if err := env.Parse(&c); err != nil {
		return c, fmt.Errorf("parse client config: %w", err)
	}

	return c, nil
}
