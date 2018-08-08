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

	reviewResponse.Allowed = true
	if s.skipNamespace(pod) {
		return &reviewResponse
	}
	if !s.checkPodLimitsExist(pod) {
		reviewResponse.Result = &metav1.Status{Message: "LimitSentry: pod rejected because of missing limits"}
		reviewResponse.Allowed = false
		return &reviewResponse
	}
	if s.Config.Memory != nil {
		if !s.checkPodLimitsMemInRange(pod) {
			reviewResponse.Result = &metav1.Status{Message: "LimitSentry: pod rejected because some containers are outside the memory limits"}
			reviewResponse.Allowed = false
			return &reviewResponse
		}
	}
	if s.Config.CPU != nil {
		if !s.checkPodLimitsCPUInRange(pod) {
			reviewResponse.Result = &metav1.Status{Message: "LimitSentry: pod rejected because some containers are outside the cpu limits"}
			reviewResponse.Allowed = false
			return &reviewResponse
		}
	}
	return &reviewResponse
}

func (s *Sentry) skipNamespace(p corev1.Pod) bool {
	for _, v := range s.Config.Ignored {
		if v == p.ObjectMeta.Namespace {
			return true
		}
	}
	return false
}

func (s *Sentry) checkPodLimitsExist(p corev1.Pod) bool {
	if !s.checkContainerLimitsExist(p.Spec.InitContainers) {
		return false
	}
	return s.checkContainerLimitsExist(p.Spec.Containers)
}

func (s *Sentry) checkContainerLimitsExist(containers []corev1.Container) bool {
	for _, c := range containers {
		if c.Resources.Limits.Cpu().IsZero() || c.Resources.Limits.Memory().IsZero() {
			return false
		}

	}
	return true
}

func (s *Sentry) checkPodLimitsMemInRange(p corev1.Pod) bool {
	if !s.checkContainerLimitsMemInRange(p.Spec.InitContainers) {
		return false
	}
	return s.checkContainerLimitsMemInRange(p.Spec.Containers)
}

func (s *Sentry) checkContainerLimitsMemInRange(containers []corev1.Container) bool {
	for _, c := range containers {
		if !s.Config.Memory.Between(c.Resources.Limits[corev1.ResourceMemory]) {
			return false
		}
	}
	return true
}

func (s *Sentry) checkPodLimitsCPUInRange(p corev1.Pod) bool {
	if !s.checkContainerLimitsCPUInRange(p.Spec.InitContainers) {
		return false
	}
	return s.checkContainerLimitsCPUInRange(p.Spec.Containers)
}

func (s *Sentry) checkContainerLimitsCPUInRange(containers []corev1.Container) bool {
	for _, c := range containers {
		if !s.Config.CPU.Between(c.Resources.Limits[corev1.ResourceCPU]) {
			return false
		}
	}
	return true
}
