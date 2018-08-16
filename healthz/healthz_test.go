package healthz

import (
	"io/ioutil"
	"log"
	"testing"

	"k8s.io/api/admission/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
)

var (
	podnoliveliness []byte
	podnoreadiness  []byte
	podpass         []byte
)

func init() {
	var err error
	podpass, err = ioutil.ReadFile("podtest.json.pass")
	if err != nil {
		log.Fatal(err)
	}
	podnoliveliness, err = ioutil.ReadFile("podtest.json.noliveliness")
	if err != nil {
		log.Fatal(err)
	}
	podnoreadiness, err = ioutil.ReadFile("podtest.json.noreadiness")
	if err != nil {
		log.Fatal(err)
	}
}

func TestAdmit(t *testing.T) {
	hs := HealthzSentry{}
	ar := v1beta1.AdmissionReview{
		Request: &v1beta1.AdmissionRequest{
			Object: runtime.RawExtension{
				Raw: podpass,
			},
		},
	}
	resp := hs.Admit(ar)
	if !resp.Allowed {
		t.Fatal("expected passing review")
	}
	ar.Request.Object.Raw = podnoliveliness
	resp = hs.Admit(ar)
	if resp.Allowed {
		t.Fatal("Expected no liveliness to fail")
	}
	ar.Request.Object.Raw = podnoreadiness
	resp = hs.Admit(ar)
	if resp.Allowed {
		t.Fatal("Expected no readiness to fail")
	}
}
