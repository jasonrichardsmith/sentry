package sentry

import (
	"k8s.io/apimachinery/pkg/api/resource"
)

type Config struct {
	CPU    MinMax `yaml:"cpu"`
	Memory MinMax `yaml:"memory"`
}

type MinMax struct {
	Min string `yaml:"Min"`
	Max string `yaml:"Max"`
}

func (c *Config) LoadSentry() (LimitSentry, error) {
	var ls LimitSentry
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
