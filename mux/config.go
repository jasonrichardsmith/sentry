package mux

import (
	"flag"
	"io/ioutil"

	"github.com/jasonrichardsmith/sentry/limits"
	"github.com/jasonrichardsmith/sentry/sentry"
	yaml "gopkg.in/yaml.v2"
)

var (
	configfile string
)

func init() {
	flag.StringVar(&configfile, "sentry-config", "config.yaml", "Location of sentry config file.")
}

type Config struct {
	Limits limits.Config `yaml:"limits"`
}

func New() *Config {
	l := limits.Config{}
	return &Config{
		Limits: l,
	}
}

func (c *Config) LoadFromFile() error {
	if !flag.Parsed() {
		flag.Parse()
	}
	configbytes, err := ioutil.ReadFile(configfile)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(configbytes, &c)
}

func (c *Config) LoadSentry() (sentry.Sentry, error) {
	return NewFromConfig(*c)
}
