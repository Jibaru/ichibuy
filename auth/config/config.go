package config

import (
	"errors"
	"log"
	"os"
	"reflect"

	"github.com/joho/godotenv"
)

type Config struct {
	APIPort            string `env:"API_PORT"`
	APIBaseURI         string `env:"API_BASE_URI"`
	WebBaseURI         string `env:"WEB_BASE_URI"`
	PostgresURI        string `env:"POSTGRES_URI"`
	JWTPrivateKey      string `env:"JWT_PRIVATE_KEY"`
	GoogleClientID     string `env:"GOOGLE_CLIENT_ID"`
	GoogleClientSecret string `env:"GOOGLE_CLIENT_SECRET"`
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
