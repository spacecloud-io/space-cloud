package k8s

import (
	"context"

	"github.com/spacecloud-io/space-cloud/managers/configman/common"
	"github.com/spacecloud-io/space-cloud/managers/source"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (k *K8s) loadConfiguration() error {

	sourcesGVR := source.GetRegisteredSources()
	config := make(common.ConfigType)
	for _, srcGVR := range sourcesGVR {
		srcList, err := k.dc.Resource(srcGVR).List(context.Background(), metav1.ListOptions{})
		if err != nil {
			return err
		}

		for _, obj := range srcList.Items {
			key := source.GetModuleName(srcGVR)
			config[key] = append(config[key], &obj)
		}
	}

	k.setConfig(config)
	return nil
}

func (k *K8s) setConfig(newConfig common.ConfigType) {
	k.lock.Lock()
	defer k.lock.Unlock()

	k.configuration = newConfig
}

func (k *K8s) getConfig() common.ConfigType {
	k.lock.RLock()
	defer k.lock.RUnlock()

	return k.configuration
}

func (k *K8s) copyConfig() common.ConfigType {
	copy := make(common.ConfigType)
	for k, v := range k.configuration {
		copy[k] = v
	}
	return copy
}
