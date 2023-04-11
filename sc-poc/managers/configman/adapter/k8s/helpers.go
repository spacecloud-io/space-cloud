package k8s

import (
	"context"

	"github.com/spacecloud-io/space-cloud/managers/source"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (k *K8s) loadConfiguration() error {

	sourcesGVR := source.GetRegisteredSources()
	for _, srcGVR := range sourcesGVR {
		srcList, err := k.dc.Resource(srcGVR).List(context.Background(), metav1.ListOptions{})
		if err != nil {
			return err
		}

		for _, obj := range srcList.Items {
			key := source.GetModuleName(srcGVR)
			k.configuration[key] = append(k.configuration[key], &obj)
		}
	}

	return nil
}
