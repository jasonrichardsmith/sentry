package sentry

import (
	"k8s.io/api/admission/v1beta1"
)

type Sentry interface {
	Admit(v1beta1.AdmissionReview) *v1beta1.AdmissionResponse
}

type Loader interface {
	LoadSentry() (Sentry, error)
}
