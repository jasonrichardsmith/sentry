package nslabels

import (
	"io/ioutil"
	"log"
	"reflect"
	"strings"
	"testing"

	"k8s.io/api/admission/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
)

var (
	nspass []byte
	nsfail []byte
)

func init() {
	var err error
	nspass, err = ioutil.ReadFile("nstest.json.pass")
	if err != nil {
		log.Fatal(err)
	}
	nsfail, err = ioutil.ReadFile("nstest.json.nolabel")
	if err != nil {
		log.Fatal(err)
	}
}

func TestType(t *testing.T) {
	ts := Sentry{}
	if ts.Type() != "Namespace" {
		t.Fatal("Failed type test")
	}
}

func TestName(t *testing.T) {
	c := Config{}
	if c.Name() != "nslabels" {
		t.Fatal("Failed name test")
	}
}

func TestAdmit(t *testing.T) {
	s := Sentry{}
	ar := v1beta1.AdmissionReview{
		Request: &v1beta1.AdmissionRequest{
			Object: runtime.RawExtension{
				Raw: nspass,
			},
		},
	}
	resp := s.Admit(ar)
	if !resp.Allowed {
		t.Fatal("expected passing review")
	}
	ar.Request.Object.Raw = nsfail
	resp = s.Admit(ar)
	if resp.Allowed {
		t.Fatal("expected failed review")
	}
	ar.Request.Object.Raw = nsfail[0:1]
	resp = s.Admit(ar)
	if !strings.Contains(resp.Result.Message, "json parse error") {
		t.Fatal("Expecting json parse error")
	}
}

func TestLoadSentry(t *testing.T) {
	c := Config{}
	s := c.LoadSentry()
	if !reflect.DeepEqual(s, Sentry{}) {
		t.Fatal("Sentry mismatch")
	}
}
