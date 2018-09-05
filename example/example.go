package example

import (
	log "github.com/Sirupsen/logrus"
	"github.com/jasonrichardsmith/sentry/sentry"
	"k8s.io/api/admission/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ExampleSentry satisfies sentry.Sentry interface
type ExampleSentry struct{}

// Type returns the type of object you are checking
func (es ExampleSentry) Type() string {
	// Return your type here.
	return "Pod"
}

// This is your validating function
func (es ExampleSentry) Admit(receivedAdmissionReview v1beta1.AdmissionReview) *v1beta1.AdmissionResponse {

	reviewResponse := v1beta1.AdmissionResponse{}
	raw := receivedAdmissionReview.Request.Object.Raw
	pod := corev1.Pod{}
	if err := sentry.Decode(raw, &pod); err != nil {
		log.Error(err)
		reviewResponse.Result = &metav1.Status{Message: err.Error()}
		return &reviewResponse
	}
	// Here you would validate your pod
	// We are just returning true here.
	reviewResponse.Allowed = true
	return &reviewResponse
}
