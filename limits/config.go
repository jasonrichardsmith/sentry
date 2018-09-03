package limits

import (
	"github.com/jasonrichardsmith/sentry/sentry"
	"k8s.io/apimachinery/pkg/api/resource"
)

const (
	NAME = "limits"
)

type Config struct {
	CPU    MinMax `yaml:"cpu"`
	Memory MinMax `yaml:"memory"`
}

type MinMax struct {
	Min string `yaml:"min"`
	Max string `yaml:"max"`
}

func (c *Config) Name() string {
	return NAME
}

func (c *Config) LoadSentry() sentry.Sentry {
	var ls LimitSentry
	var err error
	ls.MemoryMax, err = resource.ParseQuantity(c.Memory.Max)
	if err != nil {
		return ls
	}
	ls.MemoryMin, err = resource.ParseQuantity(c.Memory.Min)
	if err != nil {
		return ls
	}
	ls.CPUMin, err = resource.ParseQuantity(c.CPU.Min)
	if err != nil {
		return ls
	}
	ls.CPUMax, err = resource.ParseQuantity(c.CPU.Max)
	return ls
}
