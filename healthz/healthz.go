package healthz

import (
	log "github.com/Sirupsen/logrus"
	"github.com/jasonrichardsmith/sentry/sentry"
	"k8s.io/api/admission/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	HealthzNotPresent = "HealthzSentry: pod rejected because of missing health checks"
)

type HealthzSentry struct{}

func (hs HealthzSentry) Type() string {
	return "Pod"
}

func (hs HealthzSentry) Admit(receivedAdmissionReview v1beta1.AdmissionReview) *v1beta1.AdmissionResponse {
	log.Info("Checking health checks are present")
	raw := receivedAdmissionReview.Request.Object.Raw
	pod := corev1.Pod{}
	reviewResponse := v1beta1.AdmissionResponse{}
	if err := sentry.Decode(raw, &pod); err != nil {
		log.Error(err)
		reviewResponse.Result = &metav1.Status{Message: err.Error()}
		return &reviewResponse
	}
	reviewResponse.Allowed = true
	if !hs.checkHealthChecksExist(pod) {
		reviewResponse.Result = &metav1.Status{Message: HealthzNotPresent}
		reviewResponse.Allowed = false
		return &reviewResponse
	}
	return &reviewResponse
}

func (hs *HealthzSentry) checkHealthChecksExist(p corev1.Pod) bool {
	for _, c := range p.Spec.Containers {
		if c.LivenessProbe == nil || c.ReadinessProbe == nil {
			log.Infof("%c missing health or readiness", c.Name)
			return false
		}
	}
	return true
}
