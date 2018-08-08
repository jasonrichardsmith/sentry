package main

import (
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"k8s.io/api/admission/v1beta1"
	admissionregistrationv1beta1 "k8s.io/api/admissionregistration/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

type Config struct {
	CertFile string
	KeyFile  string
}

var (
	config = Config{
		CertFile: "/etc/webhook/certs/cert.pem",
		KeyFile:  "/etc/webhook/certs/key.pem",
	}
	scheme = runtime.NewScheme()
	codecs = serializer.NewCodecFactory(scheme)
)

type AdmissionController interface {
	Admit(v1beta1.AdmissionReview) *v1beta1.AdmissionResponse
}

type Handler struct {
	ac AdmissionController
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var body []byte
	if r.Body != nil {
		if data, err := ioutil.ReadAll(r.Body); err == nil {
			body = data
		}
	}
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		log.Errorf("contentType=%s, expect application/json", contentType)
		return
	}
	var admissionResponse *v1beta1.AdmissionResponse
	receivedAdmissionReview := v1beta1.AdmissionReview{}
	deserializer := codecs.UniversalDeserializer()
	if _, _, err := deserializer.Decode(body, nil, &receivedAdmissionReview); err != nil {
		log.Error(err)
		admissionResponse = admissionResponseError(err)
	} else {
		admissionResponse = h.ac.Admit(receivedAdmissionReview)

	}
	returnedAdmissionReview := v1beta1.AdmissionReview{}
	if admissionResponse != nil {
		returnedAdmissionReview.Response = admissionResponse
		returnedAdmissionReview.Response.UID = receivedAdmissionReview.Request.UID
	}
	responseInBytes, err := json.Marshal(returnedAdmissionReview)
	if err != nil {
		log.Error(err)
		return
	}
	log.Info("Writing response")
	if _, err := w.Write(responseInBytes); err != nil {
		log.Error(err)
	}
}

func admissionResponseError(err error) *v1beta1.AdmissionResponse {
	return &v1beta1.AdmissionResponse{
		Result: &metav1.Status{
			Message: err.Error(),
		},
	}
}

func init() {
	corev1.AddToScheme(scheme)
	admissionregistrationv1beta1.AddToScheme(scheme)
}

func main() {
	sCert, err := tls.LoadX509KeyPair(config.CertFile, config.KeyFile)
	if err != nil {
		log.Fatal(err)
	}
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{sCert},
	}

	server := &http.Server{
		Handler: Handler{
			ac: sentry.New(),
		},
		Addr:      ":443",
		TLSConfig: tlsConfig,
	}
	server.ListenAndServeTLS("", "")
}
