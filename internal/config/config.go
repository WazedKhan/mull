package config

const defaultAddr = ":8080"

// Config holds runtime settings read from the environment.
type Config struct {
	Addr string
}

// Load builds a Config from getenv (pass os.Getenv in main).
func Load(getenv func(string) string) Config {
	cfg := Config{Addr: defaultAddr}
	if v := getenv("MULL_ADDR"); v != "" {
		cfg.Addr = v
	}
	return cfg
}
