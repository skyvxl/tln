// Package server implements the tln server side.
package server

import (
	"context"
	"crypto/tls"
	"fmt"
	"log/slog"

	"github.com/skyvxl/tln/internal/config"
)

// Server represents the TLN server.
type Server struct {
	cfg       config.ServerConfig
	log       *slog.Logger
	tlsConfig *tls.Config
}

// NewServer creates a new Server instance with the provided configuration and logger.
//
//nolint:gocritic // Constructor is called only once in the application lifecycle.
func NewServer(cfg config.ServerConfig, log *slog.Logger) (*Server, error) {
	tlsCfg, err := LoadTLSConfig(cfg.TLSCertFile, cfg.TLSKeyFile)
	if err != nil {
		return nil, fmt.Errorf("init server: load tls config: %w", err)
	}
	server := &Server{
		cfg:       cfg,
		log:       log,
		tlsConfig: tlsCfg,
	}
	return server, nil
}

// Run starts the server and blocks until the context is done.
func (s *Server) Run(ctx context.Context) error {
	s.log.Info("server started", "control_addr", s.cfg.ControlAddr)
	defer s.log.Info("server stopped")
	return s.runControl(ctx)
}
