package limits

import (
	"github.com/jasonrichardsmith/Sentry/sentry"
	"k8s.io/apimachinery/pkg/api/resource"
)

type Config struct {
	CPU           MinMax `yaml:"cpu"`
	Memory        MinMax `yaml:"memory"`
	sentry.Config `yaml:"-,inline"`
}

type MinMax struct {
	Min string `yaml:"min"`
	Max string `yaml:"max"`
}

func (c *Config) LoadSentry() (sentry.Sentry, error) {
	var ls LimitSentry
	var err error
	ls.MemoryMax, err = resource.ParseQuantity(c.Memory.Max)
	if err != nil {
		return ls, err
	}
	ls.MemoryMin, err = resource.ParseQuantity(c.Memory.Min)
	if err != nil {
		return ls, err
	}
	ls.CPUMin, err = resource.ParseQuantity(c.CPU.Min)
	if err != nil {
		return ls, err
	}
	ls.CPUMax, err = resource.ParseQuantity(c.CPU.Max)
	return ls, err
}
