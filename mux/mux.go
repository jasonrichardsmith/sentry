package mux

import (
	log "github.com/Sirupsen/logrus"
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

func NewFromConfig(c Config) (SentryMux, error) {
	sm := SentryMux{
		Sentries: make([]sentryModule, 0),
	}
	if c.Limits.Enabled {
		log.Info("Limits enabled loading")
		s, err := c.Limits.LoadSentry()
		if err != nil {
			return sm, err
		}
		mod := sentryModule{
			s,
			c.Limits.IgnoredNamespaces,
		}
		log.Info("Ignoring Namespaces ", mod.ignored)
		sm.Sentries = append(sm.Sentries, mod)
	}
	if c.Healthz.Enabled {
		log.Info("Healthz enabled loading")
		s, err := c.Healthz.LoadSentry()
		if err != nil {
			return sm, err
		}
		mod := sentryModule{
			s,
			c.Healthz.IgnoredNamespaces,
		}
		log.Info("Ignoring Namespaces ", mod.ignored)
		sm.Sentries = append(sm.Sentries, mod)
	}
	if c.Tags.Enabled {
		log.Info("Tags enabled loading")
		s, err := c.Tags.LoadSentry()
		if err != nil {
			return sm, err
		}
		mod := sentryModule{
			s,
			c.Tags.IgnoredNamespaces,
		}
		log.Info("Ignoring Namespaces ", mod.ignored)
		sm.Sentries = append(sm.Sentries, mod)
	}
	if c.Domains.Enabled {
		log.Info("Domains enabled loading")
		s, err := c.Domains.LoadSentry()
		if err != nil {
			return sm, err
		}
		mod := sentryModule{
			s,
			c.Domains.IgnoredNamespaces,
		}
		log.Info("Ignoring Namespaces ", mod.ignored)
		if v, ok := sm.Sentries[c.Domains.Type]; ok {
			v["source"] = mod
		} else {
			sm.Sentries[c.Domains.Type] = map[string]sentryModule{"source": mod}
		}
	}
	return sm, nil
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
	log.Infof("Itterating over %v sentries.", receivedAdmissionReview.Request.Kind.Kind, len(sm.Sentries))
	for k, sm := range sm.Sentries {
		if receivedAdmissionReview.Request.Kind.Kind == sm.Type() {
			if !sm.Ignore(receivedAdmissionReview.Request.Namespace) {
				log.Infof("Running admit for %v", k)
				ar := sm.Admit(receivedAdmissionReview)
				if !ar.Allowed {
					log.Infof("Not allowed by %v", k)
					return ar
				}
				log.Infof("Allowed by %v", k)
			}
		}

	}
	reviewResponse := v1beta1.AdmissionResponse{}
	reviewResponse.Allowed = true
	return &reviewResponse

}
