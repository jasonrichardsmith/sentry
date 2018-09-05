package sentry

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"k8s.io/api/admission/v1beta1"
)

type FakeSentry struct{}

func (fs FakeSentry) Type() string {
	return "Pod"
}
func (fs FakeSentry) Admit(ar v1beta1.AdmissionReview) *v1beta1.AdmissionResponse {

	reviewResponse := v1beta1.AdmissionResponse{}
	reviewResponse.Allowed = true
	return &reviewResponse
}

func TestAdmissionResponseError(t *testing.T) {
	err := errors.New("test error")
	ar := admissionResponseError(err)
	if ar.Result.Message != err.Error() {
		t.Fatal("Expected error and admission response message to match")
	}
}

func TestSentryHandler(t *testing.T) {
	s := SentryHandler{FakeSentry{}}
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/healthz", nil)
	s.ServeHTTP(w, r)
	if w.Code != http.StatusOK {
		t.Fatal("Expecting 200 for healthz check")
	}
	if w.Body == nil {
		t.Fatal("Expecting body for healthz check")
	}
	data, err := ioutil.ReadAll(w.Body)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(data, healthResponse) {
		t.Fatal("Expected health response body")
	}
	r = httptest.NewRequest("POST", "/", nil)
	w = httptest.NewRecorder()
	s.ServeHTTP(w, r)
	if w.Code != http.StatusUnsupportedMediaType {
		t.Fatal("Expecting 415 for wrong content type")
	}
	if w.Body == nil {
		t.Fatal("Expecting body for wrong content type")
	}
	data, err = ioutil.ReadAll(w.Body)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(data, wrongContentResponse) {
		t.Fatal("Expected wrong content type response body")
	}
	r.Header.Add("Content-Type", "application/json")
	w = httptest.NewRecorder()
	s.ServeHTTP(w, r)
	rar := v1beta1.AdmissionReview{}
	data, err = ioutil.ReadAll(w.Body)
	if err != nil {
		t.Fatal(err)
	}
	err = Decode(data, &rar)
	if err != nil {
		t.Fatal(err)
	}
	if rar.Response.Result.Message != ErrNoUID.Error() {
		t.Fatal("Expected no uid error")
	}
	sar := v1beta1.AdmissionReview{
		Request: &v1beta1.AdmissionRequest{
			UID: "test",
		},
	}
	sarbytes, err := json.Marshal(sar)
	if err != nil {
		t.Fatal(err)
	}
	sarbytesreader := bytes.NewReader(sarbytes)
	r = httptest.NewRequest("POST", "/", sarbytesreader)
	r.Header.Add("Content-Type", "application/json")
	s.ServeHTTP(w, r)
	data, err = ioutil.ReadAll(w.Body)
	if err != nil {
		t.Fatal(err)
	}
	err = Decode(data, &rar)
	if err != nil {
		t.Fatal(err)
	}
	if rar.Response.Allowed != true {
		t.Fatal("Expected allowed response")
	}
}

func TestNewSentryServer(t *testing.T) {
	_, err := NewSentryServer(FakeSentry{})
	if err == nil {
		t.Fatal("expected pem error")
	}
}
