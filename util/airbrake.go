package util

import (
	"github.com/airbrake/gobrake/v5"
	"github.com/versolabs/verso/config"
)

func Airbrake(config *config.Config) *gobrake.Notifier {
	logger := Logger()

	if config.Airbrake.ProjectKey == "" || config.Airbrake.ProjectID == 0 {
		return nil
	}

	logger.Info().Msg("Airbrake enabled.")

	var airbrake = gobrake.NewNotifierWithOptions(&gobrake.NotifierOptions{
		ProjectId:   config.Airbrake.ProjectID,
		ProjectKey:  config.Airbrake.ProjectKey,
		Environment: config.Env,
	})

	return airbrake
}
