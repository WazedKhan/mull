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
	"github.com/labstack/echo/v5"
)

const shutdownTimeout = 10 * time.Second

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	cfg := config.Load(os.Getenv)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	sc := echo.StartConfig{
		Address:         cfg.Addr,
		HideBanner:      true,
		HidePort:        true,
		GracefulTimeout: shutdownTimeout,
		OnShutdownError: func(err error) {
			log.Error("graceful shutdown failed", "err", err)
		},
	}

	log.Info("starting api", "addr", cfg.Addr)
	if err := sc.Start(ctx, httpapi.NewServer(log)); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Error("server stopped", "err", err)
		os.Exit(1)
	}
	log.Info("api stopped")
}
