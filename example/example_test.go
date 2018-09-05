package example

import (
	"testing"

	"k8s.io/api/admission/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
)

var (
	faildecode = "{"
)

func TestType(t *testing.T) {
	ts := ExampleSentry{}
	if ts.Type() != "Pod" {
		t.Fatal("Failed type test")
	}
}

func TestName(t *testing.T) {
	c := Config{}
	if c.Name() != "example" {
		t.Fatal("Failed name test")
	}
}

func TestAdmit(t *testing.T) {
	s := ExampleSentry{}
	ar := v1beta1.AdmissionReview{
		Request: &v1beta1.AdmissionRequest{
			Object: runtime.RawExtension{
				Raw: make([]byte, 0),
			},
		},
	}
	resp := s.Admit(ar)
	if !resp.Allowed {
		t.Fatal("expected passing review")
	}
	ar.Request.Object.Raw = []byte(faildecode)
	resp = s.Admit(ar)
	if resp.Allowed {
		t.Fatal("expected passing from failed decode")
	}
}
