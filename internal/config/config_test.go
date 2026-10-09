package config_test

import (
	"testing"

	"github.com/WazedKhan/mull/internal/config"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name     string
		env      map[string]string
		wantAddr string
	}{
		{name: "default address", env: map[string]string{}, wantAddr: ":8080"},
		{name: "address from env", env: map[string]string{"MULL_ADDR": ":9000"}, wantAddr: ":9000"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := config.Load(func(k string) string { return tt.env[k] })

			if got.Addr != tt.wantAddr {
				t.Errorf("Addr = %q, want %q", got.Addr, tt.wantAddr)
			}
		})
	}
}
