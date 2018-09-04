package mux

import (
	"testing"

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
