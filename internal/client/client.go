// Package client implements the TLN client side.
package client

import (
	"context"
	"crypto/tls"
	"fmt"
	"log/slog"

	"github.com/skyvxl/tln/internal/config"
)

// Client represents the TLN client.
type Client struct {
	cfg       config.ClientConfig
	tlsConfig *tls.Config
	log       *slog.Logger
}

// NewClient create a new Client instance with the provided configuration and logger
//
//nolint:gocritic // Constructor is called only once in the application lifecycle.
func NewClient(cfg config.ClientConfig, log *slog.Logger) (*Client, error) {
	tlsCfg, err := LoadTLSConfig(cfg.CAFile, cfg.ServerName)
	if err != nil {
		return nil, fmt.Errorf("failed to load TLS config: %w", err)
	}
	return &Client{
		cfg:       cfg,
		tlsConfig: tlsCfg,
		log:       log,
	}, nil
}

// Run starts the client, establishes a TLS connection to the server, and handles the echo communication.
func (c *Client) Run(ctx context.Context) error {
	dialer := tls.Dialer{Config: c.tlsConfig}
	conn, err := dialer.DialContext(ctx, "tcp", c.cfg.ServerAddr)
	if err != nil {
		return fmt.Errorf("dial tcp: %w", err)
	}
	defer func() { _ = conn.Close() }()
	c.log.Info("connected", "remote_addr", conn.RemoteAddr().String())

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// watchdog
	go func() {
		<-ctx.Done()
		_ = conn.Close()
	}()

	return runEcho(ctx, conn, c.log)
}
