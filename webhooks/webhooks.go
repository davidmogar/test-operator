package webhooks

import (
	"github.com/davidmogar/test-operator/webhooks/auth"
	"github.com/go-logr/logr"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

// setupFunctions is a list of register functions to be invoked so all webhooks are added to the Manager
var setupFunctions = []func(manager.Manager, logr.Logger) error{
	auth.SetupWebhook,
}

// SetupWebhooks invoke all SetupWebhook functions defined in setupFunctions, setting all webhooks up and
// adding them to the Manager.
func SetupWebhooks(manager manager.Manager) error {
	log := logf.Log.WithName("webhooks")

	for _, function := range setupFunctions {
		if err := function(manager, log); err != nil {
			return err
		}
	}

	return nil
}
