package source

import (
	"github.com/jasonrichardsmith/sentry/config"
	"github.com/jasonrichardsmith/sentry/sentry"
)

const (
	NAME = "source"
)

func init() {
	config.Register(&Config{})
}

type Config struct {
	Allowed []string `mapstructure:"allowed"`
}

func (c *Config) Name() string {
	return NAME
}

func (c *Config) LoadSentry() sentry.Sentry {
	return SourceSentry{
		allowedSources: c.Allowed,
	}
}
