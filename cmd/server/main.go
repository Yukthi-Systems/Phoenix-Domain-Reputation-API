/*
Copyright (C) 2026 Yukthi Systems Private Limited

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License version 3
as published by the Free Software Foundation.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
version 3 along with this program. If not, see
<https://www.gnu.org/licenses/>.
*/

// Command server runs the domain reputation HTTP API.
package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/Yukthi-Systems/Phoenix-Domain-Reputation-API/internal/config"
	"github.com/Yukthi-Systems/Phoenix-Domain-Reputation-API/internal/httpapi"
	"github.com/Yukthi-Systems/Phoenix-Domain-Reputation-API/internal/ipfire"
	"github.com/Yukthi-Systems/Phoenix-Domain-Reputation-API/internal/reputation"
	"github.com/Yukthi-Systems/Phoenix-Domain-Reputation-API/internal/updater"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// run wires up configuration, the reputation store, the IPFire updater, and
// the HTTP server, then blocks until a shutdown signal is received or the
// server exits with an error. It performs one synchronous IPFire update
// before accepting the shutdown/serve select loop, so the service has the
// best chance of being ready immediately after startup.
func run() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load configuration: %w", err)
	}

	logger := newLogger(cfg.LogLevel)
	slog.SetDefault(logger)
	logger.Info("configuration loaded",
		"server_port", cfg.ServerPort,
		"update_interval", cfg.UpdateInterval,
		"http_timeout", cfg.HTTPTimeout,
		"api_key_configured", cfg.APIKey != "",
	)
	if cfg.APIKey == "" {
		logger.Warn("API_KEY is not set: /v1 reputation endpoints are unauthenticated; set API_KEY before exposing this service outside local development")
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	store := reputation.NewStore()
	client := ipfire.NewClient(cfg.HTTPTimeout)
	u := updater.New(client, store, cfg.Categories, cfg.UpdateInterval, logger)

	handler := httpapi.NewHandler(store, logger)
	server := &http.Server{
		Addr:              ":" + cfg.ServerPort,
		Handler:           httpapi.NewRouter(handler, cfg.APIKey, logger),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	serverErrs := make(chan error, 1)
	go func() {
		logger.Info("http server started", "addr", server.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErrs <- err
			return
		}
		serverErrs <- nil
	}()

	logger.Info("starting initial ipfire update")
	if err := u.UpdateOnce(ctx); err != nil {
		logger.Error("initial ipfire update failed, service will report not-ready until the next successful update", "error", err)
	}

	go u.Run(ctx)

	select {
	case <-ctx.Done():
		logger.Info("shutdown signal received")
	case err := <-serverErrs:
		if err != nil {
			return fmt.Errorf("http server: %w", err)
		}
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("http server shutdown: %w", err)
	}

	logger.Info("shutdown complete")
	return nil
}

// newLogger builds a structured JSON slog.Logger writing to stdout at the
// given level. An unrecognized level falls back to INFO rather than
// failing startup over a logging misconfiguration.
func newLogger(level string) *slog.Logger {
	var lvl slog.Level
	if err := lvl.UnmarshalText([]byte(strings.ToUpper(level))); err != nil {
		lvl = slog.LevelInfo
	}
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: lvl})
	return slog.New(handler)
}
