package auth

import (
	"context"
	"encoding/json"
	"github.com/davidmogar/test-operator/api/v1alpha1"
	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	"net/http"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

type auth struct {
	Log     logr.Logger
	decoder *admission.Decoder
}

// Handle yields a response to an AdmissionRequest.
//
// The supplied context is extracted from the received http.Request, allowing wrapping
// http.Handlers to inject values into and control cancelation of downstream request processing.
func (a *auth) Handle(ctx context.Context, req admission.Request) admission.Response {
	a.Log.Info("Handling resource authorization", "Operation", req.Operation, "User", req.UserInfo)

	test := &v1alpha1.Test{}

	err := a.decoder.Decode(req, test)
	if err != nil {
		return admission.Errored(http.StatusInternalServerError, errors.Wrap(err, "error decoding request"))
	}

	test.Annotations["author"] = req.UserInfo.Username

	marshaledTest, err := json.Marshal(test)
	if err != nil {
		return admission.Errored(http.StatusInternalServerError, errors.Wrap(err, "error encoding response object"))
	}

	return admission.PatchResponseFromRaw(req.Object.Raw, marshaledTest)
}

// InjectDecoder is a function that will be called automatically to inject a decoder. This is a consequence of auth
// implementing admission.DecoderInjector.
func (a *auth) InjectDecoder(decoder *admission.Decoder) error {
	a.decoder = decoder

	return nil
}

// +kubebuilder:webhook:path=/mutate-damoreno-redhat-com-v1alpha1-test,mutating=true,failurePolicy=fail,sideEffects=None,groups=damoreno.redhat.com,resources=tests,verbs=create;update,versions=v1alpha1,name=mtest.kb.io,admissionReviewVersions=v1

// SetupWebhook registers a new mutating webhook for the resources in the kubebuilder annotation above.
func SetupWebhook(mgr ctrl.Manager, log logr.Logger) error {
	mgr.GetWebhookServer().Register("/mutate-damoreno-redhat-com-v1alpha1-test", &webhook.Admission{
		Handler: &auth{
			Log: log,
		},
	})

	return nil
}
