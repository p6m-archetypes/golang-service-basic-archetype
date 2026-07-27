package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"{{ module_path }}/internal/config"
	"{{ module_path }}/internal/handler"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("config load failed", "error", err)
		os.Exit(1)
	}

	// Structured logging: JSON in production, text in development
	var logHandler slog.Handler
	if cfg.LoggingJSON {
		logHandler = slog.NewJSONHandler(os.Stdout, nil)
	} else {
		logHandler = slog.NewTextHandler(os.Stdout, nil)
	}
	logger := slog.New(logHandler)
	slog.SetDefault(logger)

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Handler: handler.New(cfg.ServiceName),
	}

	mgmtMux := http.NewServeMux()
	mgmtMux.HandleFunc("/health/readiness", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{"status":"ok"}`)
	})
	mgmtMux.HandleFunc("/health/liveness", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{"status":"ok"}`)
	})
	mgmtMux.Handle("/metrics", promhttp.Handler())
	mgmt := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", cfg.Host, cfg.ManagementPort),
		Handler: mgmtMux,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		slog.Info("management server starting", "port", cfg.ManagementPort)
		if err := mgmt.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("management server error", "error", err)
		}
	}()

	go func() {
		slog.Info("service server starting", "port", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("service server error", "error", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	slog.Info("shutting down...")

	shutCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = srv.Shutdown(shutCtx)
	_ = mgmt.Shutdown(shutCtx)
}
