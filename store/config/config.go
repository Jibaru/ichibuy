package config

import (
	"errors"
	"log"
	"os"
	"reflect"

	"github.com/joho/godotenv"
)

type Config struct {
	APIPort         string `env:"API_PORT"`
	PostgresURI     string `env:"POSTGRES_URI"`
	AuthBaseURL     string `env:"AUTH_BASE_URL"`
	FStorageBaseURL string `end:"FSTORAGE_BASE_URL"`
}

func Load() Config {
	if err := godotenv.Load(); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			log.Println(".env not found, using environment variables as default")
		} else {
			log.Fatal("error loading .env", err)
		}
	}

	var cfg Config
	loadFromEnv(&cfg)
	return cfg
}

func loadFromEnv(cfg *Config) {
	v := reflect.ValueOf(cfg).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		if envTag := fieldType.Tag.Get("env"); envTag != "" {
			if envValue := os.Getenv(envTag); envValue != "" {
				field.SetString(envValue)
			}
		}
	}
}
