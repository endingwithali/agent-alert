package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"agent-monitor/internal/config"
	"agent-monitor/internal/sessionpoller"
)

func main() {
	if err := run(); err != nil {
		slog.Error("service stopped", "error", err)
		os.Exit(1)
	}
}

func run() error {
	serviceConfig, err := config.Load(".env")
	if err != nil {
		return fmt.Errorf("load configuration: %w", err)
	}

	sessionPath := flag.String("session-path", serviceConfig.SessionPath, "Codex session file or directory")
	pollInterval := flag.Duration("poll-interval", serviceConfig.PollInterval, "interval between file polls")
	flag.Parse()

	if *pollInterval <= 0 {
		return fmt.Errorf("poll interval must be greater than zero")
	}

	poller, err := sessionpoller.New(*sessionPath)
	if err != nil {
		return fmt.Errorf("create session poller: %w", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	slog.Info(
		"session reader started",
		"session_path", *sessionPath,
		"poll_interval", pollInterval.String(),
	)

	poll := func() {
		snapshot, err := poller.Poll()
		if err != nil {
			slog.Error("session file poll failed", "error", err)
			return
		}

		slog.Info(
			"session file poll succeeded",
			"path", snapshot.Path,
			"bytes_read", len(snapshot.Data),
			"modified_at", snapshot.ModifiedAt,
		)
	}

	// Poll immediately on startup rather than waiting for the first interval.
	poll()

	ticker := time.NewTicker(*pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			slog.Info("session reader stopped")
			return nil
		case <-ticker.C:
			poll()
		}
	}
}
