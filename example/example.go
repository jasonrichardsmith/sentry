package example

import "k8s.io/api/admission/v1beta1"

var (
// Uncomment to deserialize
/*
	scheme = runtime.NewScheme()
	codecs = serializer.NewCodecFactory(scheme)
*/
)

// ExampleSentry satisfies sentry.Sentry interface
type ExampleSentry struct{}

// Type returns the type of object you are checking
func (es ExampleSentry) Type() string {
	// Return your type here.
	return "Pod"
}

// This is your validating function
func (es ExampleSentry) Admit(receivedAdmissionReview v1beta1.AdmissionReview) *v1beta1.AdmissionResponse {

	reviewResponse := v1beta1.AdmissionResponse{}
	// This is an en example of how to deserialize an object.
	/*
		raw := receivedAdmissionReview.Request.Object.Raw
		pod := corev1.Pod{}
		deserializer := codecs.UniversalDeserializer()
		reviewResponse := v1beta1.AdmissionResponse{}
		if _, _, err := deserializer.Decode(raw, nil, &pod); err != nil {
			log.Error(err)
			reviewResponse.Result = &metav1.Status{Message: err.Error()}
			return &reviewResponse
		}
	*/
	// Here you would validate your pod
	// We are just returning true here.
	reviewResponse.Allowed = true
	return &reviewResponse
}
