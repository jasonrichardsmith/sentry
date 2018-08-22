package mux

import (
	"flag"
	"io/ioutil"

	"github.com/jasonrichardsmith/sentry/source"
	"github.com/jasonrichardsmith/sentry/healthz"
	"github.com/jasonrichardsmith/sentry/images"
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
	Limits  limits.Config  `yaml:"limits"`
	Healthz healthz.Config `yaml:"healthz"`
	Images  images.Config  `yaml:"images"`
	Domains source.Config `yaml:"source"`
}

func New() *Config {
	l := limits.Config{}
	h := healthz.Config{}
	i := images.Config{}
	d := source.Config{}
	return &Config{
		Limits:  l,
		Healthz: h,
		Images:  i,
		Domains: d,
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
