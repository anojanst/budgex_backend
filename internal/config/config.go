package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port int

	DatabaseURL string
}

func Load() (Config, error) {
	// Load .env if present (env vars still override)
	_ = godotenv.Load()

	cfg := Config{
		Port: envInt("PORT", 8080),

		DatabaseURL: envStr("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/budgex?sslmode=disable"),
	}
	return cfg, nil
}

func (c Config) PostgresDSN() string {
	// Return the database URL directly
	return c.DatabaseURL
}

func envStr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func envInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		var out int
		fmt.Sscanf(v, "%d", &out)
		if out != 0 {
			return out
		}
	}
	return def
}
