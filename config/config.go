package config

import (
	"errors"
	"flag"
	"io/ioutil"

	log "github.com/Sirupsen/logrus"
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

func (c *Config) Load(file string) error {
	log.Infof("Loading config from %v", file)
	configbytes, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	iconfig := make(map[string]interface{})
	err = yaml.Unmarshal(configbytes, &iconfig)
	if err != nil {
		return err
	}
	log.Infof("Found %v configurations", len(iconfig))
	enabled := make([]ModuleConfig, 0)
	for _, mc := range c.Modules {
		log.Infof("Checking for presence of %v", mc.Name())
		if v, ok := iconfig[mc.Name()]; ok {
			cc := CommonConfig{}
			err := mapstructure.Decode(v, &cc)
			if err != nil {
				return err
			}
			log.Infof("Checking if %v is enabled", mc.Name())
			if cc.Enabled {
				log.Infof("%v enabled, loading config", mc.Name())
				c.ignored[mc.Name()] = cc.IgnoredNamespaces
				log.Infof("Ignored namespaces for of %v: %v", mc.Name(), cc.IgnoredNamespaces)
				dc := mapstructure.DecoderConfig{Result: mc}
				if cd, ok := c.decoders[mc.Name()]; ok {
					log.Infof("%v decode hook found", mc.Name())
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
	if !flag.Parsed() {
		flag.Parse()
	}
	return DefaultConfig.Load(configfile)
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
