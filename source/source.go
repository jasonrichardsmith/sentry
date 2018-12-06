package source

import (
	log "github.com/Sirupsen/logrus"
	"github.com/jasonrichardsmith/sentry/sentry"
	"k8s.io/api/admission/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	SourceUnappoved = "SourceSentry: pod rejected because image is not in allowed list"
)

func Min(a, b int) int {
	if a > b {
		return b
	} else {
		return a
	}
}

type SourceSentry struct {
	allowedSources []string
}

func (ss SourceSentry) Type() string {
	return "Pod"
}

func (ss SourceSentry) Admit(receivedAdmissionReview v1beta1.AdmissionReview) *v1beta1.AdmissionResponse {
	log.Info("Checking source are approved")
	raw := receivedAdmissionReview.Request.Object.Raw
	pod := corev1.Pod{}
	reviewResponse := v1beta1.AdmissionResponse{}
	if err := sentry.Decode(raw, &pod); err != nil {
		log.Error(err)
		reviewResponse.Result = &metav1.Status{Message: err.Error()}
		return &reviewResponse
	}
	reviewResponse.Allowed = true
	if !ss.checkImageDomainAllowed(pod) {
		reviewResponse.Result = &metav1.Status{Message: SourceUnappoved}
		reviewResponse.Allowed = false
		return &reviewResponse
	}
	return &reviewResponse
}

func (ss *SourceSentry) checkImageDomainAllowed(p corev1.Pod) bool {
	if !ss.checkImageDomainAllowedContainer(p.Spec.Containers) {
		log.Info("Checking container source are approved")
		return false
	}
	if !ss.checkImageDomainAllowedContainer(p.Spec.InitContainers) {
		log.Info("Checking initcontainer source are approved")
		return false
	}
	return true
}

func (ss *SourceSentry) checkImageDomainAllowedContainer(cs []corev1.Container) bool {
	for _, c := range cs {
		pass := false
		for _, v := range ss.allowedSources {
			if c.Image[0:Min(len(v), len(c.Image))] == v {
				log.Infof("Found approved source %v for container %v", v, c.Image)
				pass = true
			}
		}
		if !pass {
			log.Infof("%v has no approved source", c.Image)
			return false
		}
	}
	return true
}
