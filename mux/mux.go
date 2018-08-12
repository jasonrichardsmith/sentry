package mux

import (
	"github.com/jasonrichardsmith/Sentry/sentry"
	"k8s.io/api/admission/v1beta1"
	"reflect"
)

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
		sc := v.Field(i).Interface().(SentryConfig)
		if !sc.Enabled {
			continue
		}
		s, err := sc.Config.LoadSentry()
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
			sm.Sentries[sc.Type] = []sentryModule{mod}
		}
	}
	return sm, nil
}

func (sm sentryModule) Ignore(namespace string) bool {
	for _, ignore := range sm.ignored {
		if ignore == namespace {
			return true

		}
	}
	return false
}

func (sm SentryMux) Admit(receivedAdmissionReview v1beta1.AdmissionReview) *v1beta1.AdmissionResponse {
	if sms, ok := sm.Sentries[receivedAdmissionReview.Request.Kind.Kind]; ok {
		for _, sm := range sms {
			if !sm.Ignore(receivedAdmissionReview.Request.Namespace) {
				ar := sm.Admit(receivedAdmissionReview)
				if !ar.Allowed {
					return ar
				}
			}
		}

	}
	reviewResponse := v1beta1.AdmissionResponse{}
	reviewResponse.Allowed = true
	return &reviewResponse

}
