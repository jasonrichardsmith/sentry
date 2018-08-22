package source

import (
	"github.com/jasonrichardsmith/sentry/sentry"
)

type Config struct {
	sentry.Config `yaml:"-,inline"`
	AllowedSource []string `yaml:"allowed"`
}

func (c *Config) LoadSentry() (sentry.Sentry, error) {
	return DomainsSentry{
		allowedSource: c.AllowedSource,
	}, nil
}
