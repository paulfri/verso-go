package util

import (
	"log"

	"github.com/caarlos0/env/v8"
)

type Config struct {
	BaseURL           string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	DatabaseMigrate   bool   `env:"DATABASE_MIGRATE" envDefault:"false"`
	DatabaseURL       string `env:"DATABASE_URL"`
	Debug             bool   `env:"DEBUG" envDefault:"false"`
	Env               string `env:"VERSO_ENV" envDefault:"development"`
	Host              string `env:"HOST"`
	Port              string `env:"PORT" envDefault:"8080"`
	RedisURL          string `env:"REDIS_URL"`
	WorkerConcurrency int    `env:"WORKER_CONCURRENCY" envDefault:"10"`
}

func GetConfig() Config {
	config := Config{}
	if err := env.Parse(&config); err != nil {
		log.Fatalf("Error parsing config: %+v\n", err)
	}

	return config
}
