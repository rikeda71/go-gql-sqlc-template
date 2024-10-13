package internal

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

// Config is the configuration for the API server.
type Config struct {
	Port            int  `envconfig:"PORT" default:"8000"`
	GracefulTimeout int  `envconfig:"GRACEFUL_TIMEOUT" default:"30"`
	DebugMode       bool `envconfig:"DEBUG_MODE" default:"false"`
	/// DB
	DatabaseUser     string `envconfig:"DATABASE_USER" required:"true"`
	DatabasePassword string `envconfig:"DATABASE_PASSWORD" required:"true"`
	DatabaseHost     string `envconfig:"DATABASE_HOST" required:"true"`
	DatabaseName     string `envconfig:"DATABASE_NAME" required:"true"`
	DatabasePort     int    `envconfig:"DATABASE_PORT" default:"5432"`
}

func (cnf *Config) DataSource() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		cnf.DatabaseUser, cnf.DatabasePassword, cnf.DatabaseHost, cnf.DatabasePort, cnf.DatabaseName)
}

func NewConfig() (*Config, error) {
	conf := &Config{}
	if err := envconfig.Process("", conf); err != nil {
		return nil, err
	}
	return conf, nil
}
