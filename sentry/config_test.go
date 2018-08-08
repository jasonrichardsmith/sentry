package sentry

import (
	"testing"
)

func TestLoad(t *testing.T) {
	c := NewConfig()
	err := c.Load()
	if err != nil {
		t.Fatal(err)
	}
}
