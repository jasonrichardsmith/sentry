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
	Sentries map[string]map[string]sentryModule
}

func NewFromConfig(c Config) (SentryMux, error) {
	sm := SentryMux{
		Sentries: make(map[string]map[string]sentryModule),
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
		sm.Sentries[c.Limits.Type] = map[string]sentryModule{"limits": mod}
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
		if v, ok := sm.Sentries[c.Healthz.Type]; ok {
			v["healthz"] = mod
		} else {
			sm.Sentries[c.Healthz.Type] = map[string]sentryModule{"healthz": mod}
		}
	}
	if c.Images.Enabled {
		log.Info("Images enabled loading")
		s, err := c.Images.LoadSentry()
		if err != nil {
			return sm, err
		}
		mod := sentryModule{
			s,
			c.Images.IgnoredNamespaces,
		}
		log.Info("Ignoring Namespaces ", mod.ignored)
		if v, ok := sm.Sentries[c.Images.Type]; ok {
			v["images"] = mod
		} else {
			sm.Sentries[c.Images.Type] = map[string]sentryModule{"images": mod}
		}
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
			v["domains"] = mod
		} else {
			sm.Sentries[c.Domains.Type] = map[string]sentryModule{"domains": mod}
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

func (sm SentryMux) Admit(receivedAdmissionReview v1beta1.AdmissionReview) *v1beta1.AdmissionResponse {
	log.Infof("Received request of kind %v", receivedAdmissionReview.Request.Kind.Kind)
	if sms, ok := sm.Sentries[receivedAdmissionReview.Request.Kind.Kind]; ok {
		log.Infof("Found sentries for kind %v, itterating over %v sentries.", receivedAdmissionReview.Request.Kind.Kind, len(sms))
		for k, sm := range sms {
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
