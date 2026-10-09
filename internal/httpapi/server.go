package httpapi

import (
	"log/slog"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

// NewServer returns an Echo instance with middleware and routes registered.
func NewServer(log *slog.Logger) *echo.Echo {
	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogMethod:   true,
		LogURI:      true,
		LogStatus:   true,
		HandleError: true,
		LogValuesFunc: func(c *echo.Context, v middleware.RequestLoggerValues) error {
			attrs := []slog.Attr{
				slog.String("method", v.Method),
				slog.String("uri", v.URI),
				slog.Int("status", v.Status),
			}
			if v.Error != nil {
				attrs = append(attrs, slog.String("err", v.Error.Error()))
			}
			log.LogAttrs(c.Request().Context(), levelFor(v.Status), "request", attrs...)
			return nil
		},
	}))

	e.GET("/healthz", health)
	return e
}

// levelFor maps an HTTP status to a log level: 5xx Error, 4xx Warn, else Info.
func levelFor(status int) slog.Level {
	switch {
	case status >= 500:
		return slog.LevelError
	case status >= 400:
		return slog.LevelWarn
	default:
		return slog.LevelInfo
	}
}
