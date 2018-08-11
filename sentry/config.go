package sentry

import (
	"flag"
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
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
	Config            *Loader
}

type Loader interface {
	LoadSentry() (Sentry, error)
}

func New() Config {
	c := Config{
		Limits: SentryConfig{
			Config: &limits.Config{},
		},
	}
}

func (c *Config) LoadSentry() (SentryMux, error) {
	var s SentryMux
	if !flag.Parsed() {
		flag.Parse()
	}
	configbytes, err := ioutil.ReadFile(configfile)
	if err != nil {
		return s, err
	}
	err := c.Unmarshal(configbytes)
	if err != nil {
		return s, err
	}
	return NewFromConfig(c)

}

func (c *Config) Unmarshal(b []byte) error {
	return yaml.Unmarshal(b, c)
}
