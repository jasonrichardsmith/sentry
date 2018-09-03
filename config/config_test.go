package config

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/jasonrichardsmith/sentry/limits"
)

func TestLoading(t *testing.T) {
	l := limits.Config{}
	err := Register(&l)
	if err != nil {
		t.Fatal(err)
	}
	err = Load()
	if err != nil {
		t.Fatal(err)
	}
	spew.Dump(l)
}
