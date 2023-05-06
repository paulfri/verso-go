package config

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
	BaseURL string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	Debug   bool   `env:"DEBUG" envDefault:"false"`
	Env     string `env:"VERSO_ENV" envDefault:"development"`
}

type airbrakeConfig struct {
	ProjectID  int64  `toml:"project_id"`
	ProjectKey string `toml:"project_key"`
}

type databaseConfig struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	Database string `toml:"database"`
	User     string `toml:"user"`
	Password string `toml:"password"`
	SSLMode  string `toml:"ssl_mode"`
	Migrate  bool   `toml:"migrate, omitempty"`
	Conn     string `toml:"conn"`
	URL      string `toml:"url"`
}

type serverConfig struct {
	Host string `toml:"host"`
	Port int    `toml:"port"`
}

type workerConfig struct {
	Concurrency int    `toml:"concurrency"`
	RedisURL    string `toml:"redis_url"`
}
