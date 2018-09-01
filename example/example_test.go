package example

import (
	"testing"

	"k8s.io/api/admission/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
)

func TestLoadSentry(t *testing.T) {
	c := Config{}
	_, err := c.LoadSentry()
	if err != nil {
		t.Fatal(err)
	}
}

func TestType(t *testing.T) {
	ts := ExampleSentry{}
	if ts.Type() != "Pod" {
		t.Fatal("Failed type test")
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
}
