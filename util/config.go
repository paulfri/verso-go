package util

import (
	"log"

	"github.com/caarlos0/env/v8"
)

type Config struct {
	Host              string `env:"HOST"`
	Port              string `env:"PORT" envDefault:"8080"`
	Env               string `env:"VERSO_ENV" envDefault:"development"`
	DatabaseUrl       string `env:"DATABASE_URL"`
	RedisUrl          string `env:"REDIS_URL"`
	WorkerConcurrency int    `env:"WORKER_CONCURRENCY" envDefault:"10"`
}

func GetConfig() Config {
	config := Config{}
	if err := env.Parse(&config); err != nil {
		log.Fatalf("Error parsing config: %+v\n", err)
	}

	return config
}
