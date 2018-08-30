package source

import (
	"github.com/jasonrichardsmith/sentry/sentry"
)

type Config struct {
	sentry.Config  `yaml:"-,inline"`
	AllowedSources []string `yaml:"allowed"`
}

func (c *Config) LoadSentry() (sentry.Sentry, error) {
	return SourceSentry{
		allowedSources: c.AllowedSources,
	}, nil
}
