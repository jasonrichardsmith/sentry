package healthz

import (
	"io/ioutil"
	"log"
	"strings"
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
	c := Config{}
	hs, err := c.LoadSentry()
	if err != nil {
		log.Fatal(err)
	}
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
	ar.Request.Object.Raw = podpass[0:5]
	resp = hs.Admit(ar)
	if !strings.Contains(resp.Result.Message, "json parse error") {
		t.Fatal("Expecting json parse error")
	}
}
