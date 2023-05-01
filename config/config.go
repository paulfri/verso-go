package config

import "fmt"

func GetConfig() *Config {
	env := getEnv()
	conf := loadConfig(env)

	return conf
}

type Config struct {
	Airbrake airbrakeConfig
	Database databaseConfig
	Server   serverConfig
	Worker   workerConfig

	// TODO implement - ported from env
	BaseURL  string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	Debug    bool   `env:"DEBUG" envDefault:"false"`
	Env      string `env:"VERSO_ENV" envDefault:"development"`
	RedisURL string `env:"REDIS_URL"`
}

type airbrakeConfig struct {
	ProjectID  int64  `toml:"project_id"`
	ProjectKey string `toml:"project_key"`
}

type databaseConfig struct {
	Host        string `toml:"host"`
	Port        int    `toml:"port"`
	Database    string `toml:"database"`
	User        string `toml:"user"`
	Password    string `toml:"password"`
	SSLDisabled bool   `toml:"ssl_disabled,omitempty"`
	Migrate     bool   `toml:"migrate, omitempty"`
}

type serverConfig struct {
	Host string `toml:"host"`
	Port int    `toml:"port"`
}

type workerConfig struct {
	Concurrency int `toml:"concurrency"`
}

func (c *databaseConfig) URL() string {
	var sslMode string
	if c.SSLDisabled {
		sslMode = "?sslmode=disable"
	}

	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s%s",
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.Database,
		sslMode,
	)
}
