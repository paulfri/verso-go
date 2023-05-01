package config

import (
	"bytes"
	"log"
	"text/template"

	"github.com/BurntSushi/toml"
)

func loadConfig(env environment) *Config {
	// Find the correct file for the given application environment (VERSO_ENV).
	// e.g., for development, the file is ./config/env/development.toml.
	filepath := configFilePathFromEnv(env.Env)

	// Parse the config file as a template and execute the template against the
	// loaded environment variables
	var encoded string
	buf := new(bytes.Buffer)
	tmpl, err := template.ParseFiles(filepath)
	if err != nil {
		log.Fatalf("Error parsing config template: %+v\n", err)
	}
	_ = tmpl.Execute(buf, env)
	encoded = string(buf.String())

	// Decode the templatized TOML config file into a Config struct.
	var conf Config
	_, err = toml.Decode(encoded, &conf)
	if err != nil {
		log.Fatalf("Error loading config: %+v\n", err)
	}

	return &conf
}
