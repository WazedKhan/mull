package httpapi

import (
	"log/slog"
	"testing"
)

func TestLevelFor(t *testing.T) {
	tests := []struct {
		name   string
		status int
		want   slog.Level
	}{
		{"ok", 200, slog.LevelInfo},
		{"not found", 404, slog.LevelWarn},
		{"server error", 500, slog.LevelError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := levelFor(tt.status); got != tt.want {
				t.Errorf("levelFor(%d) = %v, want %v", tt.status, got, tt.want)
			}
		})
	}
}
