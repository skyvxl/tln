// Package main is the entry point for the tln server application.
package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/skyvxl/tln/internal/buildinfo"
	"github.com/skyvxl/tln/internal/config"
	"github.com/skyvxl/tln/internal/logger"
	"github.com/skyvxl/tln/internal/server"
)

func main() {
	showVersion := flag.Bool("version", false, "print version and exit")
	flag.Parse()

	if *showVersion {
		fmt.Println("tlnd", buildinfo.String())

		return
	}

	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "fatal: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	if err := config.LoadDotenv(); err != nil {
		return err
	}

	cfg, err := config.LoadServer()
	if err != nil {
		return err
	}

	log := logger.New(os.Stderr, logger.Format(cfg.LogFormat), cfg.LogLevel)
	slog.SetDefault(log)

	log.Info(
		"tlnd starting",
		"version", buildinfo.Version,
		"commit", buildinfo.Commit,
		"log_format", cfg.LogFormat,
		"log_level", cfg.LogLevel,
		"control_addr", cfg.ControlAddr,
		"tls_cert", cfg.TLSCertFile,
	)

	srv, err := server.NewServer(cfg, log)
	if err != nil {
		return err
	}
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	err = srv.Run(ctx)
	if err != nil {
		return err
	}
	log.Info("tlnd stopped")
	return nil
}
