package mux

import (
	log "github.com/Sirupsen/logrus"
	"github.com/jasonrichardsmith/sentry/config"
	"github.com/jasonrichardsmith/sentry/sentry"
	"k8s.io/api/admission/v1beta1"
)

type sentryModule struct {
	sentry.Sentry
	ignored []string
}

type SentryMux struct {
	Sentries []sentryModule
}

func New(c config.Config) SentryMux {
	sm := SentryMux{
		Sentries: make([]sentryModule, 0),
	}
	for _, v := range c.Modules {
		sm.Sentries = append(sm.Sentries,
			sentryModule{
				v.LoadSentry(),
				c.Ignored(v.Name()),
			})
	}
	return sm
}

func (sm sentryModule) Ignore(namespace string) bool {
	log.Infof("Checking to see ignored namespace %v", namespace)
	for _, ignore := range sm.ignored {
		if ignore == namespace {
			return true
			log.Infof("Namespace %v ignored", namespace)

		}
	}
	log.Infof("Namespace %v not ignored", namespace)
	return false
}
func (sm SentryMux) Type() string {
	return "*"
}
func (sm SentryMux) Admit(receivedAdmissionReview v1beta1.AdmissionReview) *v1beta1.AdmissionResponse {
	log.Infof("Received request of kind %v", receivedAdmissionReview.Request.Kind.Kind)
	log.Infof("Itterating over %v sentries.", receivedAdmissionReview.Request.Kind.Kind)
	for _, sm := range sm.Sentries {
		if receivedAdmissionReview.Request.Kind.Kind == sm.Type() {
			if !sm.Ignore(receivedAdmissionReview.Request.Namespace) {
				log.Infof("Running admit for %v", sm.Type())
				ar := sm.Admit(receivedAdmissionReview)
				if !ar.Allowed {
					log.Infof("Not allowed by %v", sm.Type())
					return ar
				}
				log.Infof("Allowed by %v", sm.Type())
			}
		}

	}
	reviewResponse := v1beta1.AdmissionResponse{}
	reviewResponse.Allowed = true
	return &reviewResponse
}
