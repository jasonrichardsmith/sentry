package sentry

import (
	"flag"
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
	"k8s.io/apimachinery/pkg/api/resource"
)

type Config struct {
	CPU     *MinMax  `yaml:"CPU"`
	Memory  *MinMax  `yaml:"Memory"`
	Audit   bool     `yaml:"Audit"`
	Ignored []string `yaml:"Ignored"`
}

var (
	configfile string
)

func init() {
	flag.StringVar(&configfile, "sentry-config", "config.yaml", "Location of sentry config file.")
}

func NewConfig() Config {
	return Config{}
}

func (c *Config) Load() error {
	if !flag.Parsed() {
		flag.Parse()
	}
	configbytes, err := ioutil.ReadFile(configfile)
	if err != nil {
		return err
	}
	err = c.Unmarshal(configbytes)
	if err != nil {
		return err
	}
	return c.SetResources()
}

func (c *Config) SetResources() (err error) {
	if c.Memory != nil {
		err = c.Memory.SetResources()
		if err != nil {
			return err
		}
	}
	if c.CPU != nil {
		return c.CPU.SetResources()
	}
	return nil
}

func (c *Config) Unmarshal(b []byte) error {
	return yaml.Unmarshal(b, c)
}

type MinMax struct {
	Min    string            `yaml:"Min"`
	Max    string            `yaml:"Max"`
	qtyMin resource.Quantity `yaml:"-"`
	qtyMax resource.Quantity `yaml:"-"`
}

func (mm *MinMax) SetResources() (err error) {
	mm.qtyMax, err = resource.ParseQuantity(mm.Max)
	if err != nil {
		return err
	}
	mm.qtyMin, err = resource.ParseQuantity(mm.Min)
	return err
}

func (mm *MinMax) Between(q resource.Quantity) bool {
	if mm.qtyMax.Cmp(q) >= 0 && mm.qtyMin.Cmp(q) <= 0 {
		return true
	}
	return false
}
