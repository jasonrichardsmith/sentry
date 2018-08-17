package domains

import (
	"github.com/jasonrichardsmith/sentry/sentry"
)

type Config struct {
	sentry.Config  `yaml:"-,inline"`
	AllowedDomains []string `yaml:"allowedDomains"`
}

func (c *Config) LoadSentry() (sentry.Sentry, error) {
	return DomainsSentry{
		allowedDomains: c.AllowedDomains,
	}, nil
}
