// Package main is the entry point for the tln client application.
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
	"github.com/skyvxl/tln/internal/client"
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
	if err := config.LoadDotenv(); err != nil {
		return err
	}

	cfg, err := config.LoadClient()
	if err != nil {
		return err
	}

	log := logger.New(os.Stderr, logger.Format(cfg.LogFormat), cfg.LogLevel)
	slog.SetDefault(log)

	log.Info(
		"tlnc starting",
		"version", buildinfo.Version,
		"commit", buildinfo.Commit,
	)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cli, err := client.NewClient(cfg, log)
	if err != nil {
		return err
	}

	err = cli.Run(ctx)
	if err != nil {
		return err
	}

	return nil
}
