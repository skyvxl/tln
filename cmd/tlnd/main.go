// Package main is the entry point for the tln server application.
package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/skyvxl/tln/internal/buildinfo"
	"github.com/skyvxl/tln/internal/config"
	"github.com/skyvxl/tln/internal/logger"
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

	log.Info("tlnd starting",
		"version", buildinfo.Version,
		"commit", buildinfo.Commit,
		"log_format", cfg.LogFormat,
		"log_level", cfg.LogLevel,
		"control_addr", cfg.ControlAddr,
		"tls_cert", cfg.TLSCertFile,
	)
	log.Info("tlnd exiting (stub)")

	return nil
}
