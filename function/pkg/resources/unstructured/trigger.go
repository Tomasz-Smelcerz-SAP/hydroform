package unstructured

import (
	"fmt"

	"github.com/kyma-incubator/hydroform/function/pkg/workspace"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

const (
	triggerApiVersion = "eventing.knative.dev/v1alpha1"
	triggerNameFormat = "%s-%s"
)

func NewTriggers(cfg workspace.Cfg) ([]unstructured.Unstructured, error) {
	var list []unstructured.Unstructured

	for _, triggerInfo := range cfg.Triggers {
		trigger := unstructured.Unstructured{
			Object: map[string]interface{}{},
		}
		triggerAttributes := triggerInfo.Attributes()
		subscriberRef := asSubscriberRef(cfg)

		triggerName := fmt.Sprintf(triggerNameFormat, cfg.Name, triggerInfo.Source)
		if triggerInfo.Name != "" {
			triggerName = triggerInfo.Name
		}

		decorators := Decorators{
			decorateWithField(triggerApiVersion, "apiVersion"),
			decorateWithField("Trigger", "kind"),
			decorateWithMetadata(triggerName, cfg.Namespace),
			decorateWithLabels(cfg.Labels),
			decorateWithField("default", "spec", "broker"),
			decorateWithMap(triggerAttributes, "spec", "filter", "attributes"),
			decorateWithMap(subscriberRef, "spec", "subscriber", "ref"),
		}

		if err := decorate(&trigger, decorators); err != nil {
			return list, err
		}
		list = append(list, trigger)
	}
	return list, nil
}

type subscriberRef = map[string]interface{}

func asSubscriberRef(cfg workspace.Cfg) subscriberRef {
	return map[string]interface{}{
		"apiVersion": "v1",
		"kind":       "Service",
		"name":       cfg.Name,
		"namespace":  cfg.Namespace,
	}
}
