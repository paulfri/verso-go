package util

import "github.com/airbrake/gobrake/v5"

func Airbrake(config *Config) *gobrake.Notifier {
	logger := Logger()

	if config.AirbrakeProjectKey == "" || config.AirbrakeProjectID == 0 {
		return nil
	}

	logger.Info().Msg("Airbrake enabled.")

	var airbrake = gobrake.NewNotifierWithOptions(&gobrake.NotifierOptions{
		ProjectId:   config.AirbrakeProjectID,
		ProjectKey:  config.AirbrakeProjectKey,
		Environment: config.Env,
	})

	return airbrake
}
