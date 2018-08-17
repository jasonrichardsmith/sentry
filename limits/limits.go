package limits

import (
	log "github.com/Sirupsen/logrus"
	"k8s.io/api/admission/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

var (
	scheme = runtime.NewScheme()
	codecs = serializer.NewCodecFactory(scheme)
)

const (
	LimitsNotPresent    = "LimitSentry: pod rejected because of missing limits"
	LimitsOutsideMemory = "LimitSentry: pod rejected because some containers are outside the memory limits"
	LimitsOutsideCPU    = "LimitSentry: pod rejected because some containers are outside the cpu limits"
)

type LimitSentry struct {
	MemoryMin resource.Quantity
	MemoryMax resource.Quantity
	CPUMin    resource.Quantity
	CPUMax    resource.Quantity
}

func (ls LimitSentry) BetweenCPU(q resource.Quantity) bool {
	if ls.CPUMax.Cmp(q) >= 0 && ls.CPUMin.Cmp(q) <= 0 {
		return true
	}
	return false
}

func (ls LimitSentry) BetweenMemory(q resource.Quantity) bool {
	if ls.MemoryMax.Cmp(q) >= 0 && ls.MemoryMin.Cmp(q) <= 0 {
		return true
	}
	return false
}
func (ls LimitSentry) Admit(receivedAdmissionReview v1beta1.AdmissionReview) *v1beta1.AdmissionResponse {
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
	if !ls.checkPodLimitsExist(pod) {
		log.Info("Pod limits are not present")
		reviewResponse.Result = &metav1.Status{Message: LimitsNotPresent}
		reviewResponse.Allowed = false
		return &reviewResponse
	}
	if !ls.checkPodLimitsMemInRange(pod) {
		log.Info("Pod memory limits are not in range")
		reviewResponse.Result = &metav1.Status{Message: LimitsOutsideMemory}
		reviewResponse.Allowed = false
		return &reviewResponse
	}
	if !ls.checkPodLimitsCPUInRange(pod) {
		log.Info("Pod cpu limits are not in range")
		reviewResponse.Result = &metav1.Status{Message: LimitsOutsideCPU}
		reviewResponse.Allowed = false
		return &reviewResponse
	}
	return &reviewResponse
}

func (ls *LimitSentry) checkPodLimitsExist(p corev1.Pod) bool {
	if !ls.checkContainerLimitsExist(p.Spec.InitContainers) {
		return false
	}
	return ls.checkContainerLimitsExist(p.Spec.Containers)
}

func (ls *LimitSentry) checkContainerLimitsExist(containers []corev1.Container) bool {
	for _, c := range containers {
		if c.Resources.Limits.Cpu().IsZero() || c.Resources.Limits.Memory().IsZero() {
			log.Infof("%v is missing limits", c.Name)
			return false
		}

	}
	return true
}

func (ls *LimitSentry) checkPodLimitsMemInRange(p corev1.Pod) bool {
	if !ls.checkContainerLimitsMemInRange(p.Spec.InitContainers) {
		return false
	}
	return ls.checkContainerLimitsMemInRange(p.Spec.Containers)
}

func (ls *LimitSentry) checkContainerLimitsMemInRange(containers []corev1.Container) bool {
	for _, c := range containers {
		if !ls.BetweenMemory(c.Resources.Limits[corev1.ResourceMemory]) {
			log.Infof("%v memory not in range", c.Name)
			return false
		}
	}
	return true
}

func (ls *LimitSentry) checkPodLimitsCPUInRange(p corev1.Pod) bool {
	if !ls.checkContainerLimitsCPUInRange(p.Spec.InitContainers) {
		return false
	}
	return ls.checkContainerLimitsCPUInRange(p.Spec.Containers)
}

func (ls *LimitSentry) checkContainerLimitsCPUInRange(containers []corev1.Container) bool {
	for _, c := range containers {
		if !ls.BetweenCPU(c.Resources.Limits[corev1.ResourceCPU]) {
			log.Infof("%v cpu not in range", c.Name)
			return false
		}
	}
	return true
}
