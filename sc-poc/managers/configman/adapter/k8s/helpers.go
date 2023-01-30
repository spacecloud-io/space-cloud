package k8s

import (
	"context"

	"github.com/spacecloud-io/space-cloud/managers/source"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (k *K8s) loadConfiguration() error {

	sourcesGVR := source.GetSourcesGVR()
	for _, srcGVR := range sourcesGVR {
		srcList, err := k.dc.Resource(srcGVR).List(context.Background(), metav1.ListOptions{})
		if err != nil {
			return err
		}

		for _, obj := range srcList.Items {
			kind := obj.GetKind()
			key := source.GetModuleName(srcGVR)

			k.configuration[kind] = append(k.configuration[kind], &obj)
			k.configurationN[key] = append(k.configurationN[key], &obj)
		}
	}

	return nil
}
