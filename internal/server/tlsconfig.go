// Package server provides functions for setting up and configuring the server.
package server

import (
	"crypto/tls"
	"fmt"
)

// LoadTLSConfig loads TLS configuration from the specified certificate and key files.
func LoadTLSConfig(certFile, keyFile string) (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, fmt.Errorf("load TLS keypair from %s and %s: %w", certFile, keyFile, err)
	}
	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
	}, nil
}
