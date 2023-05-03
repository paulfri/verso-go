package config

import (
	"bytes"
	"fmt"
	"log"
	"text/template"

	"github.com/BurntSushi/toml"
)

func loadConfig(env environment) *Config {
	// Find the correct file for the given application environment (VERSO_ENV).
	// e.g., for development, the file is ./config/env/development.toml.
	filepath := configFilePathFromEnv(env.Env)

	// Parse the config file as a template.
	var encoded string
	buf := new(bytes.Buffer)
	tmpl, err := template.ParseFiles(filepath)
	if err != nil {
		log.Fatalf("Error parsing config template: %+v\n", err)
	}

	// Execute the template against the given structured env vars.
	err = tmpl.Execute(buf, env)
	if err != nil {
		log.Fatalf("Error executing config template: %+v\n", err)
	}
	encoded = string(buf.String())

	// Decode the templatized TOML config file into a Config struct.
	var conf Config
	_, err = toml.Decode(encoded, &conf)
	if err != nil {
		log.Fatalf("Error loading config: %+v\n", err)
	}

	// If DATABASE_CONN was provided, parse the connection string into individual
	// chunks. Takes precedence over other config.
	if env.DatabaseConn != "" {
		conf.Database = parseConnectionString(env.DatabaseConn)
	}

	// Set the database URL.
	conf.Database.URL = fmt.Sprintf(
		"postgresql://%s:%s@%s:%d/%s?sslmode=%s",
		conf.Database.User,
		conf.Database.Password,
		conf.Database.Host,
		conf.Database.Port,
		conf.Database.Database,
		conf.Database.SSLMode,
	)

	// Set the connection string, in case the config was parsed as params.
	conf.Database.Conn = fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		conf.Database.Host,
		conf.Database.Port,
		conf.Database.User,
		conf.Database.Password,
		conf.Database.Database,
		conf.Database.SSLMode,
	)

	return &conf
}

func parseConnectionString(conn string) databaseConfig {
	var host string
	var port int
	var user string
	var password string
	var database string
	var sslMode string

	_, err := fmt.Sscanf(
		conn,
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		&host,
		&port,
		&user,
		&password,
		&database,
		&sslMode,
	)

	if err != nil {
		log.Fatalf("Error parsing database URL: %+v\n", err)
	}

	return databaseConfig{
		Host:     host,
		Port:     port,
		Database: database,
		User:     user,
		Password: password,
		SSLMode:  sslMode,
		Conn:     conn,
	}
}
