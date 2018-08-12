package mux

import (
	"flag"
	"github.com/jasonrichardsmith/Sentry/sentry"
	"github.com/jasonrichardsmith/Sentry/sentry/limits"
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
)

var (
	configfile string
)

func init() {
	flag.StringVar(&configfile, "sentry-config", "config.yaml", "Location of sentry config file.")
}

type Config struct {
	Limits SentryConfig `yaml:"limits"`
}

type SentryConfig struct {
	Type              string
	Enabled           bool
	IgnoredNamespaces []string
	Config            sentry.Loader
}

func New() Config {
	l := limits.Config{}
	return Config{
		Limits: SentryConfig{
			Config: &l,
		},
	}
}

func (c Config) LoadSentry() (sentry.Sentry, error) {
	var s SentryMux
	if !flag.Parsed() {
		flag.Parse()
	}
	configbytes, err := ioutil.ReadFile(configfile)
	if err != nil {
		return s, err
	}
	err = c.Unmarshal(configbytes)
	if err != nil {
		return s, err
	}
	return NewFromConfig(c)

}

func (c *Config) Unmarshal(b []byte) error {
	return yaml.Unmarshal(b, c)
}
