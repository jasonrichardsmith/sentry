package images

import (
	"io/ioutil"
	"log"
	"testing"

	"k8s.io/api/admission/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
)

var (
	podnotag     []byte
	podlatesttag []byte
	podpass      []byte
)

func init() {
	var err error
	podpass, err = ioutil.ReadFile("podtest.json.pass")
	if err != nil {
		log.Fatal(err)
	}
	podnotag, err = ioutil.ReadFile("podtest.json.notag")
	if err != nil {
		log.Fatal(err)
	}
	podlatesttag, err = ioutil.ReadFile("podtest.json.latesttag")
	if err != nil {
		log.Fatal(err)
	}
}

func TestAdmit(t *testing.T) {
	is := ImagesSentry{}
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
	ar.Request.Object.Raw = podnotag
	resp = is.Admit(ar)
	if resp.Allowed {
		t.Fatal("Expected no tag to fail")
	}
	ar.Request.Object.Raw = podlatesttag
	resp = is.Admit(ar)
	if resp.Allowed {
		t.Fatal("Expected latest tag to fail")
	}
}
