package nslabels

import (
	log "github.com/Sirupsen/logrus"
	"github.com/jasonrichardsmith/sentry/sentry"
	"k8s.io/api/admission/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	NoLabelsPresent = "NsLabelsSentry: Namespace rejected because of no labels"
)

type Sentry struct{}

func (s Sentry) Type() string {
	return "Namespace"
}

func (s Sentry) Admit(receivedAdmissionReview v1beta1.AdmissionReview) *v1beta1.AdmissionResponse {
	log.Info("Admitting namespace")
	reviewResponse := v1beta1.AdmissionResponse{}
	raw := receivedAdmissionReview.Request.Object.Raw
	ns := corev1.Namespace{}
	if err := sentry.Decode(raw, &ns); err != nil {
		log.Error(err)
		reviewResponse.Result = &metav1.Status{Message: err.Error()}
		return &reviewResponse
	}
	reviewResponse.Allowed = true
	if len(ns.ObjectMeta.Labels) == 0 {
		log.Infof("Rejecting namespace %v because of no label", ns.ObjectMeta.GetName())
		reviewResponse.Allowed = false
		reviewResponse.Result = &metav1.Status{Message: NoLabelsPresent}
		return &reviewResponse
	}
	log.Infof("Namespace %v has labels", ns.ObjectMeta.GetName())
	return &reviewResponse
}
