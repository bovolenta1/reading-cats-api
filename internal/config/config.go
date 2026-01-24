package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
}

func Load() Config {
	// Carrega .env.local (se existir) sem quebrar prod
	_ = godotenv.Overload(".env.local")

	return Config{
		DatabaseURL: mustEnv("DATABASE_URL"),
	}
}

func mustEnv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		panic("missing env var: " + k)
	}
	return v
}
