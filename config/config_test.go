package config

import "testing"

func TestNew(t *testing.T) {
	c := New()
	if c.Modules == nil || c.decoders == nil || c.ignored == nil {
		t.Fatal("Not all config structures instantiated")
	}
}
