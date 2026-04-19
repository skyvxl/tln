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
		fmt.Println("tlnc", buildinfo.String())
		return
	}

	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "fatal: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.LoadClient()
	if err != nil {
		return err
	}

	log := logger.New(os.Stderr, logger.Format(cfg.LogFormat), cfg.LogLevel)
	slog.SetDefault(log)

	log.Info("tlnc starting",
		"version", buildinfo.Version,
		"commit", buildinfo.Commit,
	)

	log.Info("tlnc exiting (phase 0 stub)")
	return nil
}
