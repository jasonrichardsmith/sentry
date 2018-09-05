package config

import (
	"flag"
	"reflect"
	"testing"

	"github.com/jasonrichardsmith/sentry/sentry"
	"github.com/mitchellh/mapstructure"
	"k8s.io/api/admission/v1beta1"
)

var (
	ignored = map[string][]string{
		"source": {"test1"},
	}
	decoders = map[string]mapstructure.DecodeHookFunc{"test": FakeHookFunc}
)

func TestNew(t *testing.T) {
	c := New()
	if c.Modules == nil || c.decoders == nil || c.ignored == nil {
		t.Fatal("Not all config structures instantiated")
	}
}

type FakeSentry struct{}

func (fs FakeSentry) Admit(v1beta1.AdmissionReview) (ar *v1beta1.AdmissionResponse) { return }

func (fs FakeSentry) Type() string { return "test" }

type FakeModuleConfig struct {
	Allowed []string
}

func (f FakeModuleConfig) LoadSentry() sentry.Sentry { return FakeSentry{} }
func (f FakeModuleConfig) Name() string              { return "source" }

type DisabledModuleConfig struct{}

func (d DisabledModuleConfig) LoadSentry() sentry.Sentry { return FakeSentry{} }
func (d DisabledModuleConfig) Name() string              { return "healthz" }

func TestLoad(t *testing.T) {
	c := New()
	err := c.Load()
	if err != nil {
		t.Fatal(err)
	}
	if !flag.Parsed() {
		t.Fatal("Expecting parsed flags")
	}
	c = New()
	c.Register(FakeModuleConfig{})
	err = c.Load()
	if err == nil {
		t.Fatal("expect pointer error")
	}
	c = New()
	c.Register(&FakeModuleConfig{})
	err = c.Load()
	if err != nil {
		t.Fatal(err)
	}
	if len(c.Modules) != 1 {
		t.Fatal("Expected one module loaded")
	}
	if !reflect.DeepEqual(c.Ignored("source"), []string{"test1"}) {
		t.Fatal("Expected one ignore value")
	}
	if !reflect.DeepEqual(c.ignored, ignored) {
		t.Fatal("mismatched ignored")
	}
	c = New()
	c.Register(&DisabledModuleConfig{})
	err = c.Load()
	if err != nil {
		t.Fatal(err)
	}
	if len(c.Modules) != 0 {
		t.Fatal("Expected no module loaded")
	}

	configfile = "config.yaml.fail"
	err = c.Load()
	if err == nil {
		t.Fatal("Expected yaml Unmarshal error")
	}
}

func TestDefaultLoad(t *testing.T) {
	configfile = "config.yaml"
	err := Load()
	if err != nil {
		t.Fatal(err)
	}
	c := New()
	err = c.Load()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(DefaultConfig, c) {
		t.Fatal("expected matching configs")
	}
}

func TestRegister(t *testing.T) {
	c := New()
	err := c.Register(&FakeModuleConfig{})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(c.Modules[0], &FakeModuleConfig{}) {
		t.Fatal("expected matching moduleconfigs")
	}
	err = c.Register(&FakeModuleConfig{})
	if err != ERR_DUPLICATE_MODULENAME {
		t.Fatal("Expected matching error")
	}
	err = c.Register(&DisabledModuleConfig{})
	if err != nil {
		t.Fatal(err)
	}
	if len(c.Modules) != 2 {
		t.Fatal("Expected two modules registered")
	}
}

func FakeHookFunc(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
	return data, nil
}

func TestDecoder(t *testing.T) {
	Decoder("test", FakeHookFunc)
	if len(DefaultConfig.decoders) != 1 {
		t.Fatal("Expected one decoder")
	}
}

func TestIgnored(t *testing.T) {

}
