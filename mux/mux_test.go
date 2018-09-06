package mux

import (
	"testing"

	"github.com/jasonrichardsmith/sentry/config"
	"github.com/jasonrichardsmith/sentry/sentry"

	"k8s.io/api/admission/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestIgnore(t *testing.T) {
	sm := sentryModule{
		ignored: []string{
			"ignoreme",
		},
	}
	if sm.Ignore("donotignore") {
		t.Fatal("ignored unlisted namespace")
	}
	if !sm.Ignore("ignoreme") {
		t.Fatal("did not ignore listed namespace")
	}
}

func TestType(t *testing.T) {
	is := SentryMux{}
	if is.Type() != "*" {
		t.Fatal("Failed type test")
	}
}

type FakeModuleConfig struct {
}

func (f FakeModuleConfig) LoadSentry() sentry.Sentry { return FakeSentry{} }
func (f FakeModuleConfig) Name() string              { return "one" }

type FakeModuleConfig2 struct {
}

func (f FakeModuleConfig2) LoadSentry() sentry.Sentry { return FakeSentry{} }
func (f FakeModuleConfig2) Name() string              { return "two" }

type FakeSentry struct {
	admit bool
}

func (fs FakeSentry) Type() string {
	return "Pod"

}
func (fs FakeSentry) Admit(ar v1beta1.AdmissionReview) *v1beta1.AdmissionResponse {

	reviewResponse := v1beta1.AdmissionResponse{}
	reviewResponse.Allowed = fs.admit
	return &reviewResponse
}

func TestAdmit(t *testing.T) {
	mux := SentryMux{
		Sentries: []sentryModule{
			sentryModule{
				Sentry: FakeSentry{true},
				ignored: []string{
					"test1",
					"test2",
				},
			},
		},
	}
	ar := v1beta1.AdmissionReview{
		Request: &v1beta1.AdmissionRequest{
			Namespace: "test1",
			Kind: metav1.GroupVersionKind{
				Kind: "test",
			},
		},
	}
	resp := mux.Admit(ar)
	if resp.Allowed != true {
		t.Fatal("Return false on unmatched kind")
	}
	ar.Request.Kind.Kind = "Pod"
	resp = mux.Admit(ar)
	if resp.Allowed != true {
		t.Fatal("Return false on ignored namespace")
	}
	ar.Request.Namespace = "test3"
	resp = mux.Admit(ar)
	if resp.Allowed != true {
		t.Fatal("Return false expected true")
	}
	mux.Sentries = []sentryModule{
		sentryModule{
			Sentry: FakeSentry{false},
		},
	}
	resp = mux.Admit(ar)
	if resp.Allowed != false {
		t.Fatal("Return true expected false")
	}

}

func TestNew(t *testing.T) {
	c := config.New()
	sm := New(c)
	if sm.Sentries == nil || len(sm.Sentries) > 0 {
		t.Fatal("expected non nil slice of 0 length")
	}
	c.Register(FakeModuleConfig{})
	c.Register(FakeModuleConfig2{})
	sm = New(c)
	if len(sm.Sentries) != 2 {
		t.Fatal("expected 2 modules")
	}

}
