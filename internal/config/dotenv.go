// Package config provides structures and functions for loading configuration parameters from environment variables.
package config

import (
	"errors"
	"fmt"
	"io/fs"

	"github.com/joho/godotenv"
)

// LoadDotenv loads environment variables from a .env file.
//
// If the file does not exist, it returns nil. Broken .env files will return an error.
//
// If environment variables are already set in the system, they will not be overridden by values in the .env file.
func LoadDotenv() error {
	if err := godotenv.Load(); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return fmt.Errorf("load .env: %w", err)
	}
	return nil
}
