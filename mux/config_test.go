package mux

import (
	"reflect"
	"testing"

	"k8s.io/apimachinery/pkg/api/resource"

	"github.com/jasonrichardsmith/sentry/limits"
	"github.com/jasonrichardsmith/sentry/sentry"
)

func TestLoadFromFile(t *testing.T) {
	match := &Config{
		Limits: limits.Config{
			CPU: limits.MinMax{
				Min: "1G",
				Max: "1G",
			},
			Memory: limits.MinMax{
				Min: "1G",
				Max: "1G",
			},
			Config: sentry.Config{
				Type:    "Pod",
				Enabled: true,
				IgnoredNamespaces: []string{
					"test1",
					"test2",
				},
			},
		},
	}
	c := New()
	err := c.LoadFromFile()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(match, c) {
		t.Fatal("Deep Equal for marshaling from file not equal")
	}
}

func TestLoadSentry(t *testing.T) {
	c := &Config{
		Limits: limits.Config{
			CPU: limits.MinMax{
				Min: "1G",
				Max: "1G",
			},
			Memory: limits.MinMax{
				Min: "1G",
				Max: "1G",
			},
			Config: sentry.Config{
				Type:    "Pod",
				Enabled: true,
				IgnoredNamespaces: []string{
					"test1",
					"test2",
				},
			},
		},
	}
	qty, err := resource.ParseQuantity("1G")
	if err != nil {
		t.Fatal(err)
	}
	match := SentryMux{
		Sentries: map[string][]sentryModule{
			"Pod": []sentryModule{
				sentryModule{
					Sentry: limits.LimitSentry{
						MemoryMin: qty,
						MemoryMax: qty,
						CPUMin:    qty,
						CPUMax:    qty,
					},
					ignored: []string{
						"test1",
						"test2",
					},
				},
			},
		},
	}

	s, err := c.LoadSentry()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(match, s) {
		t.Fatal("Deep Equal for LoadSentry not equal")
	}
}
