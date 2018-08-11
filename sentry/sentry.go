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

type Sentry interface {
	Admit(v1beta1.AdmissionReview) *v1beta1.AdmissionResponse
}

type sentryModule struct {
	sentry.Sentry
	ignored []string
}

type SentryMux struct {
	Sentries map[string][]sentryModule
}

func NewFromConfig(c Config) (SentryMux, error) {
	sm := SentryMux{
		Sentries: make(map[string][]sentryModule),
	}
	v := reflect.ValueOf(c)
	for i := 0; i < v.NumField(); i++ {
		sc = v.Field(i).Interface().(SentryConfig)
		if !sc.Enabled {
			continue
		}
		s, err := sc.LoadSentry()
		if err != nil {
			return sm, err
		}
		mod := sentryModule{
			s,
			sc.IgnoredNamespaces,
		}
		if val, ok := sm.Sentries[sc.Type]; ok {
			val = append(val, mod)
		} else {
			sm.Sentries[sc.Type] = []SentryModule{mod}
		}
	}
	return sm, nil
}

func (s *SentryMux) Admit(receivedAdmissionReview v1beta1.AdmissionReview) *v1beta1.AdmissionResponse {
	if sms, ok := sm.Sentries[receivedAdmissionReview.Request.Kind]; ok {
		for _, sm := range sms {
			ar := sm.Admit(receivedAdmissionReview)
			if !ar.Allowed {
				return ar
			}
		}

	}
	reviewResponse := v1beta1.AdmissionResponse{}
	reviewResponse.Allowed = true
	return reviewResponse

}
