package config

import (
	"log"

	e "github.com/caarlos0/env/v8"
	"github.com/joho/godotenv"
)

// This is only for parsing environment-configurable configuration for
// templatized TOML config files.
type environment struct {
	Env string `env:"VERSO_ENV" envDefault:"development"`

	// Only includes options that are configurable by the environment, i.e.,
	// secrets.
	AirbrakeProjectID  int64  `env:"AIRBRAKE_PROJECT_ID"`
	AirbrakeProjectKey string `env:"AIRBRAKE_PROJECT_KEY"`
	DatabaseConn       string `env:"DATABASE_CONN"`
	RedisURL           string `env:"REDIS_URL"`
}

func getEnv() environment {
	// Ignore error because the .env file is optional.
	_ = godotenv.Load()

	env := environment{}
	if err := e.Parse(&env); err != nil {
		log.Fatalf("Error parsing environment: %+v\n", err)
	}

	return env
}
