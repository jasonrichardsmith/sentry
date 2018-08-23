package sentry

import (
	"crypto/tls"
	"encoding/json"
	"flag"
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

func init() {
	corev1.AddToScheme(scheme)
	admissionregistrationv1beta1.AddToScheme(scheme)
	flag.StringVar(&tlscert, "tlscert", "/etc/webhook/certs/cert.pem", "Location of TLS Cert file.")
	flag.StringVar(&tlskey, "tlskey", "/etc/webhook/certs/key.pem", "Location of TLS key file.")
}

var (
	scheme          = runtime.NewScheme()
	codecs          = serializer.NewCodecFactory(scheme)
	tlscert, tlskey string
)

type Config struct {
	Type              string   `yaml:"type"`
	Enabled           bool     `yaml:"enabled"`
	IgnoredNamespaces []string `yaml:"ignoredNamespaces"`
}

type Sentry interface {
	Admit(v1beta1.AdmissionReview) *v1beta1.AdmissionResponse
}

type Loader interface {
	LoadSentry() (Sentry, error)
}

func admissionResponseError(err error) *v1beta1.AdmissionResponse {
	return &v1beta1.AdmissionResponse{
		Result: &metav1.Status{
			Message: err.Error(),
		},
	}
}

type SentryHandler struct {
	Sentry
}

func (sh SentryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var body []byte
	if r.Body != nil {
		if data, err := ioutil.ReadAll(r.Body); err == nil {
			body = data
		}
	}
	log.Info("Received request")
	if r.URL.Path == "/healthz" {
		log.Info("Received health check")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("200 - Healthy"))
		return
	}
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		log.Errorf("contentType=%s, expect application/json", contentType)
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte("415 - Wrong Content Type"))
		return
	}
	log.Info("Correct ContentType")
	var admissionResponse *v1beta1.AdmissionResponse
	receivedAdmissionReview := v1beta1.AdmissionReview{}
	deserializer := codecs.UniversalDeserializer()
	if _, _, err := deserializer.Decode(body, nil, &receivedAdmissionReview); err != nil {
		log.Error(err)
		admissionResponse = admissionResponseError(err)
	} else {
		admissionResponse = sh.Sentry.Admit(receivedAdmissionReview)
		log.Infof("Received response of %v from sentry", admissionResponse.Allowed)

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

func NewSentryServer(s Sentry) (*http.Server, error) {
	if !flag.Parsed() {
		flag.Parse()
	}
	server := NewSentryServerNoSSL(s)
	sCert, err := tls.LoadX509KeyPair(tlscert, tlskey)
	if err != nil {
		return server, err
	}
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{sCert},
	}
	server.TLSConfig = tlsConfig
	server.Addr = ":443"
	return server, nil
}

func NewSentryServerNoSSL(s Sentry) *http.Server {
	return &http.Server{
		Handler: SentryHandler{
			Sentry: s,
		},
		Addr: ":8080",
	}
}
