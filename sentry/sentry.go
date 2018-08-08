package sentry

import (
	log "github.com/Sirupsen/logrus"
	"k8s.io/api/admission/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

var (
	scheme = runtime.NewScheme()
	codecs = serializer.NewCodecFactory(scheme)
)

type Sentry struct {
	Config
}

func New() Sentry {
	return Sentry{}
}

func NewFromConfig() (Sentry, error) {
	s := New()
	s.Config = NewConfig()
	err := s.Config.Load()
	return s, err
}

func (s *Sentry) Admit(receivedAdmissionReview v1beta1.AdmissionReview) *v1beta1.AdmissionResponse {
	log.Info("Checking limits are present")
	raw := receivedAdmissionReview.Request.Object.Raw
	pod := corev1.Pod{}
	deserializer := codecs.UniversalDeserializer()
	reviewResponse := v1beta1.AdmissionResponse{}
	if _, _, err := deserializer.Decode(raw, nil, &pod); err != nil {
		log.Error(err)
		reviewResponse.Result = &metav1.Status{Message: err.Error()}
		return &reviewResponse
	}

	// TODO: change to validation
	reviewResponse.Allowed = true

	if val, ok := pod.ObjectMeta.Annotations["mwc-example.jasonrichardsmith.com.exclude"]; ok {
		log.Info("annotation exists")
		// if the key is true we will exclude
		if val == "true" {
			log.Info("excluded due to annotation")
			return &reviewResponse
		}
	}
	return &reviewResponse
}

func checkPodLimits(p corev1.Pod) bool {
	return false
}

func checkContainersLimits(containers []corev1.Container) bool {
	for _, c := range containers {
		if c.Resources.Limits.Cpu().IsZero() || c.Resources.Limits.Memory().IsZero() {
			return false
		}
	}
	return true
	//TODO add comparison on max and min
}
