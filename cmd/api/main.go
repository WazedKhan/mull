package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/WazedKhan/mull/internal/config"
	"github.com/WazedKhan/mull/internal/httpapi"
)

const (
	shutdownTimeout = 10 * time.Second
	readTimeout     = 30 * time.Second
	headerTimeout   = 5 * time.Second
)

func main() {
	os.Exit(run())
}

// run returns the process exit code: 0 on a clean stop, 1 on a server or shutdown error.
func run() int {
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	cfg := config.Load(os.Getenv)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	srv := &http.Server{
		Addr:              cfg.Addr,
		Handler:           httpapi.NewServer(log),
		ReadHeaderTimeout: headerTimeout,
		ReadTimeout:       readTimeout,
	}

	serveErr := make(chan error, 1)
	go func() { serveErr <- srv.ListenAndServe() }()

	log.Info("starting api", "addr", cfg.Addr)
	select {
	case err := <-serveErr:
		if !errors.Is(err, http.ErrServerClosed) {
			log.Error("server stopped", "err", err)
			return 1
		}
	case <-ctx.Done():
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error("graceful shutdown failed", "err", err)
		return 1
	}
	log.Info("api stopped")
	return 0
}
