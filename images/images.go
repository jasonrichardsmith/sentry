package images

import (
	"strings"

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
	ImagesNoTag = "ImagesSentry: pod rejected because of missing image tag"
)

type ImagesSentry struct{}

func (is ImagesSentry) Admit(receivedAdmissionReview v1beta1.AdmissionReview) *v1beta1.AdmissionResponse {
	log.Info("Checking image tags are present")
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
	if !is.checkImageTagsExist(pod) {
		reviewResponse.Result = &metav1.Status{Message: ImagesNoTag}
		reviewResponse.Allowed = false
		return &reviewResponse
	}
	return &reviewResponse
}

func (is *ImagesSentry) checkImageTagsExist(p corev1.Pod) bool {
	for _, c := range p.Spec.Containers {
		split := strings.Split(c.Image, ":")
		if len(split) == 1 || split[1] == "latest" {
			log.Infof("%c has no image tag or tag of latest", c.Name)
			return false
		}
	}
	for _, c := range p.Spec.InitContainers {
		split := strings.Split(c.Image, ":")
		if len(split) == 1 || split[1] == "latest" {
			log.Infof("%c has no image tag or tag of latest", c.Name)
			return false
		}
	}
	return true
}
