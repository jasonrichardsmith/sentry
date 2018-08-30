package source

import (
	"io/ioutil"
	"log"
	"strings"
	"testing"

	"k8s.io/api/admission/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
)

var (
	podbadimage []byte
	podpass     []byte
)

func init() {
	var err error
	podpass, err = ioutil.ReadFile("podtest.json.pass")
	if err != nil {
		log.Fatal(err)
	}
	podbadimage, err = ioutil.ReadFile("podtest.json.badimage")
	if err != nil {
		log.Fatal(err)
	}
}

func TestType(t *testing.T) {
	ss := SourceSentry{}
	if ss.Type() != "Pod" {
		t.Fatal("Failed type test")
	}
}

func TestAdmit(t *testing.T) {
	c := Config{
		AllowedSources: []string{
			"this/is/allowed",
		},
	}
	is, err := c.LoadSentry()
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
	resp := is.Admit(ar)
	if !resp.Allowed {
		t.Fatal("expected passing review")
	}
	ar.Request.Object.Raw = podbadimage
	resp = is.Admit(ar)
	if resp.Allowed {
		t.Fatal("Expected badimage to fail")
	}
	ar.Request.Object.Raw = podpass[0:5]
	resp = is.Admit(ar)
	if !strings.Contains(resp.Result.Message, "json parse error") {
		t.Fatal("Expecting json parse error")
	}
}
