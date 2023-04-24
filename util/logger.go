package util

import (
	"os"

	"github.com/rs/zerolog"
)

func Logger() *zerolog.Logger {
	// TODO: log JSON format in non-dev environments
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	logger = logger.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	return &logger
}
