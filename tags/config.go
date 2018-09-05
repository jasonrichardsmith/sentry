package tags

import (
	"github.com/jasonrichardsmith/sentry/config"
	"github.com/jasonrichardsmith/sentry/sentry"
)

const (
	NAME = "tags"
)

func init() {
	config.Register(&Config{})
}

type Config struct {
}

func (c *Config) Name() string {
	return NAME
}

func (c *Config) LoadSentry() sentry.Sentry {
	return TagsSentry{}
}
