package mux

import (
	"flag"
	"io/ioutil"

	"github.com/jasonrichardsmith/sentry/example"
	"github.com/jasonrichardsmith/sentry/healthz"
	"github.com/jasonrichardsmith/sentry/limits"
	"github.com/jasonrichardsmith/sentry/sentry"
	"github.com/jasonrichardsmith/sentry/source"
	"github.com/jasonrichardsmith/sentry/tags"
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
	Source  source.Config  `yaml:"source"`
	Tags    tags.Config    `yaml:"tags"`
	Example example.Config `yaml:"example"`
}

func New() *Config {
	l := limits.Config{}
	h := healthz.Config{}
	i := tags.Config{}
	s := source.Config{}
	e := example.Config{}
	return &Config{
		Limits:  l,
		Healthz: h,
		Tags:    i,
		Source:  s,
		Example: e,
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
