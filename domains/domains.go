package domains

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

const (
	DomainsUnappovedDomain = "DomainsSentry: pod rejected because image is not in allowed domains"
)

type DomainsSentry struct {
	allowedDomains []string
}

func (ds DomainsSentry) Admit(receivedAdmissionReview v1beta1.AdmissionReview) *v1beta1.AdmissionResponse {
	log.Info("Checking domains are approved")
	raw := receivedAdmissionReview.Request.Object.Raw
	pod := corev1.Pod{}
	deserializer := codecs.UniversalDeserializer()
	reviewResponse := v1beta1.AdmissionResponse{}
	if _, _, err := deserializer.Decode(raw, nil, &pod); err != nil {
		log.Error(err)
		reviewResponse.Result = &metav1.Status{Message: err.Error()}
		return &reviewResponse
	}
	reviewResponse.Allowed = true
	if !ds.checkImageDomainAllowed(pod) {
		reviewResponse.Result = &metav1.Status{Message: DomainsUnappovedDomain}
		reviewResponse.Allowed = false
		return &reviewResponse
	}
	return &reviewResponse
}

func (ds *DomainsSentry) checkImageDomainAllowed(p corev1.Pod) bool {
	if !ds.checkImageDomainAllowedContainer(p.Spec.Containers) {
		log.Info("Checking container domains are approved")
		return false
	}
	if !ds.checkImageDomainAllowedContainer(p.Spec.InitContainers) {
		log.Info("Checking initcontainer domains are approved")
		return false
	}
	return true
}

func (ds *DomainsSentry) checkImageDomainAllowedContainer(cs []corev1.Container) bool {
	for _, c := range cs {
		pass := false
		for _, v := range ds.allowedDomains {
			if c.Image[0:len(v)] == v {
				log.Infof("Found approved domain %v for container %v", v, c.Image)
				pass = true
			}
		}
		if !pass {
			log.Infof("%v has no approved domains", c.Image)
			return false
		}
	}
	return true
}
