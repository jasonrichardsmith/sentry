package mux

import (
	"io/ioutil"
	"log"
	"testing"

	"github.com/jasonrichardsmith/Sentry/limits"
	"k8s.io/api/admission/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	testpod []byte
)

func init() {
	var err error
	testpod, err = ioutil.ReadFile("podtest.json")
	if err != nil {
		log.Fatal(err)
	}
}

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

func TestAdmit(t *testing.T) {
	mux := SentryMux{
		Sentries: map[string][]sentryModule{
			"Pod": []sentryModule{
				sentryModule{
					Sentry: limits.LimitSentry{},
					ignored: []string{
						"test1",
						"test2",
					},
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

}
