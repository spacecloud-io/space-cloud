package k8s

import (
	"context"
	"os"
	"sync"
	"time"

	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	k8sCache "k8s.io/client-go/tools/cache"
	"sigs.k8s.io/controller-runtime/pkg/client/config"

	"github.com/spacecloud-io/space-cloud/managers/configman/adapter"
	"github.com/spacecloud-io/space-cloud/managers/configman/common"
	"github.com/spacecloud-io/space-cloud/managers/source"
	"github.com/spacecloud-io/space-cloud/utils"
)

type K8s struct {
	dc            *dynamic.DynamicClient
	namespace     string
	logger        *zap.Logger
	informers     []k8sCache.SharedIndexInformer
	configuration common.ConfigType
	lock          sync.RWMutex
}

func MakeK8sAdapter() (adapter.Adapter, error) {
	logger, _ := zap.NewDevelopment()
	namespace, ok := os.LookupEnv("K8S_NAMESPACE")
	if !ok {
		namespace = "default"
	}

	restConfig, err := config.GetConfig()
	if err != nil {
		return nil, err
	}

	dc, err := dynamic.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}

	factory := dynamicinformer.NewFilteredDynamicSharedInformerFactory(dc, 5*time.Minute, namespace, nil)

	informers := []k8sCache.SharedIndexInformer{}

	sourcesGVR := source.GetRegisteredSources()
	for _, srcGVR := range sourcesGVR {
		informers = append(informers, factory.ForResource(srcGVR).Informer())
	}

	k := &K8s{
		dc:            dc,
		namespace:     namespace,
		logger:        logger,
		informers:     informers,
		configuration: make(common.ConfigType),
	}

	return k, nil
}

func (k *K8s) GetRawConfig() (common.ConfigType, error) {
	// Load SC config file from cluster
	err := k.loadConfiguration()
	if err != nil {
		return nil, err
	}

	// Load the new caddy config
	return k.configuration, nil
}

func (k *K8s) Run(ctx context.Context) (chan common.ConfigType, error) {
	cfgChan := make(chan common.ConfigType)

	for _, informer := range k.informers {
		informer.AddEventHandler(k8sCache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				u := obj.(*unstructured.Unstructured)

				k.addOrUpdateConfig(u)

				cfgChan <- k.configuration
			},
			UpdateFunc: func(_ interface{}, newObj interface{}) {
				u := newObj.(*unstructured.Unstructured)

				k.addOrUpdateConfig(u)

				cfgChan <- k.configuration
			},
			DeleteFunc: func(obj interface{}) {
				u := obj.(*unstructured.Unstructured)

				gvr := schema.GroupVersionResource{
					Group:    u.GroupVersionKind().Group,
					Version:  u.GroupVersionKind().Version,
					Resource: utils.Pluralize(u.GetKind())}
				key := source.GetModuleName(gvr)
				s := []*unstructured.Unstructured{}
				for _, spec := range k.configuration[key] {
					if spec.GetName() == u.GetName() {
						continue
					}
					s = append(s, spec)
				}
				k.configuration[key] = s

				cfgChan <- k.configuration
			},
		})

		go informer.Run(ctx.Done())
	}

	return cfgChan, nil
}

func (k *K8s) addOrUpdateConfig(u *unstructured.Unstructured) {
	gvr := schema.GroupVersionResource{
		Group:    u.GroupVersionKind().Group,
		Version:  u.GroupVersionKind().Version,
		Resource: utils.Pluralize(u.GetKind())}
	key := source.GetModuleName(gvr)
	found := false
	for i, spec := range k.configuration[key] {
		if spec.GetName() == u.GetName() {
			k.configuration[key][i] = u
			found = true
			break
		}
	}
	if !found {
		k.configuration[key] = append(k.configuration[key], u)
	}
}

// Interface guard
var (
	_ adapter.Adapter = (*K8s)(nil)
)
