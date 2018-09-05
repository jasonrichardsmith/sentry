package config

import (
	"errors"
	"flag"
	"io/ioutil"

	"github.com/jasonrichardsmith/sentry/sentry"
	"github.com/mitchellh/mapstructure"

	yaml "gopkg.in/yaml.v2"
)

func init() {
	flag.StringVar(&configfile, "sentry-config", "config.yaml", "Location of sentry config file.")
}

var (
	ERR_DUPLICATE_MODULENAME = errors.New("Duplicate module names")
	DefaultConfig            = New()
	configfile               string
)

type ModuleConfig interface {
	LoadSentry() sentry.Sentry
	Name() string
}

type CommonConfig struct {
	Enabled           bool     `mapstructure:"enabled"`
	IgnoredNamespaces []string `mapstructure:"ignoredNamespaces"`
}

type Config struct {
	Modules  []ModuleConfig
	decoders map[string]mapstructure.DecodeHookFunc
	ignored  map[string][]string
}

func New() Config {
	return Config{
		make([]ModuleConfig, 0),
		make(map[string]mapstructure.DecodeHookFunc),
		make(map[string][]string),
	}
}

func (c *Config) Load() error {
	if !flag.Parsed() {
		flag.Parse()
	}
	configbytes, err := ioutil.ReadFile(configfile)
	if err != nil {
		return err
	}
	iconfig := make(map[string]interface{})
	err = yaml.Unmarshal(configbytes, &iconfig)
	if err != nil {
		return err
	}
	enabled := make([]ModuleConfig, 0)
	for _, mc := range c.Modules {
		if v, ok := iconfig[mc.Name()]; ok {
			cc := CommonConfig{}
			err := mapstructure.Decode(v, &cc)
			if err != nil {
				return err
			}
			if cc.Enabled {
				c.ignored[mc.Name()] = cc.IgnoredNamespaces
				dc := mapstructure.DecoderConfig{Result: mc}
				if cd, ok := c.decoders[mc.Name()]; ok {
					dc.DecodeHook = cd
				}
				d, err := mapstructure.NewDecoder(&dc)
				if err != nil {
					return err
				}
				err = d.Decode(v)
				if err != nil {
					return err
				}
				enabled = append(enabled, mc)
			}
		}
	}
	c.Modules = enabled
	return nil
}

func Load() error {
	return DefaultConfig.Load()
}

func (c *Config) Register(mc ModuleConfig) error {
	for _, m := range c.Modules {
		if mc.Name() == m.Name() {
			return ERR_DUPLICATE_MODULENAME
		}
	}
	c.Modules = append(c.Modules, mc)
	return nil
}

func (c *Config) Decoder(name string, dhf mapstructure.DecodeHookFunc) {
	c.decoders[name] = dhf
}

func (c *Config) Ignored(name string) []string {
	if v, ok := c.ignored[name]; ok {
		return v
	}
	return make([]string, 0)
}

func Register(mc ModuleConfig) error {
	return DefaultConfig.Register(mc)
}

func Decoder(name string, dhf mapstructure.DecodeHookFunc) {
	DefaultConfig.Decoder(name, dhf)
}
