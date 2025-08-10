package config

import (
	"errors"
	"io/fs"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DB   DBConfig
	Port string
}

type DBConfig struct {
	DSN string
}

func New() (Config, error) {
	err := godotenv.Load()
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return Config{}, err
	}

	return Config{
		Port: getEnv("WORDDY_PORT", "8080"),
		DB: DBConfig{
			DSN: os.Getenv("WORDDY_DB_DSN"),
		},
	}, nil

}

func getEnv(key, fallback string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	return value
}
