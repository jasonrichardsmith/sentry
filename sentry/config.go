package sentry

import (
	"flag"
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
	"k8s.io/apimachinery/pkg/api/resource"
)

type Config struct {
	CPU struct {
		Min    string            `yaml:"Min"`
		Max    string            `yaml:"Max"`
		qtyMin resource.Quantity `yaml:"-"`
		qtyMax resource.Quantity `yaml:"-"`
	} `yaml:"CPU"`
	Memory struct {
		Min    string            `yaml:"Min"`
		Max    string            `yaml:"Max"`
		qtyMin resource.Quantity `yaml:"-"`
		qtyMax resource.Quantity `yaml:"-"`
	} `yaml:"Memory"`
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
	c.CPU.qtyMax, err = resource.ParseQuantity(c.CPU.Max)
	if err != nil {
		return err
	}
	c.CPU.qtyMin, err = resource.ParseQuantity(c.CPU.Min)
	if err != nil {
		return err
	}
	c.Memory.qtyMax, err = resource.ParseQuantity(c.Memory.Max)
	if err != nil {
		return err
	}
	c.Memory.qtyMin, err = resource.ParseQuantity(c.Memory.Min)
	return err
}

func (c *Config) Unmarshal(b []byte) error {
	return yaml.Unmarshal(b, c)
}

func (c *Config) ValidMemory(q resource.Quantity) bool {
	if c.Memory.qtyMax.Cmp(q) >= 0 && c.Memory.qtyMin.Cmp(q) <= 0 {
		return true
	}
	return false
}

func (c *Config) ValidCPU(q resource.Quantity) bool {
	if c.CPU.qtyMax.Cmp(q) >= 0 && c.CPU.qtyMin.Cmp(q) <= 0 {
		return true
	}
	return false
}
