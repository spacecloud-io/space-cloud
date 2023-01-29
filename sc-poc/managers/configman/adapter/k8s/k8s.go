package k8s

import (
	"context"
	"encoding/json"
	"os"
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
)

var compiledgraphqlsourcesResource = schema.GroupVersionResource{Group: "core.space-cloud.io", Version: "v1alpha1", Resource: "compiledgraphqlsources"}
var graphqlsourcesResource = schema.GroupVersionResource{Group: "core.space-cloud.io", Version: "v1alpha1", Resource: "graphqlsources"}
var jwthsasecretsResource = schema.GroupVersionResource{Group: "core.space-cloud.io", Version: "v1alpha1", Resource: "jwthsasecrets"}
var opapoliciesResource = schema.GroupVersionResource{Group: "core.space-cloud.io", Version: "v1alpha1", Resource: "opapolicies"}

type K8s struct {
	dc             *dynamic.DynamicClient
	namespace      string
	logger         *zap.Logger
	informers      []k8sCache.SharedIndexInformer
	configuration  map[string][]*unstructured.Unstructured
	configurationN map[string][]*unstructured.Unstructured
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
	informers = append(informers,
		factory.ForResource(compiledgraphqlsourcesResource).Informer(),
		factory.ForResource(graphqlsourcesResource).Informer(),
		factory.ForResource(jwthsasecretsResource).Informer(),
		factory.ForResource(opapoliciesResource).Informer())

	k := &K8s{
		dc:             dc,
		namespace:      namespace,
		logger:         logger,
		informers:      informers,
		configuration:  make(map[string][]*unstructured.Unstructured),
		configurationN: make(map[string][]*unstructured.Unstructured),
	}

	return k, nil
}

func (k *K8s) GetRawConfig() ([]byte, error) {
	// Load SC config file from cluster
	err := k.loadConfiguration()
	if err != nil {
		return nil, err
	}

	// Load the new caddy config
	return k.getCaddyConfig()
}

func (k *K8s) Run(ctx context.Context) (chan []byte, error) {
	cfgChan := make(chan []byte)

	for _, informer := range k.informers {
		informer.AddEventHandler(k8sCache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				u := obj.(*unstructured.Unstructured)

				k.addOrUpdateConfig(u)

				resp, err := k.getCaddyConfig()
				if err != nil {
					k.logger.Error("error reloading config: ", zap.Error(err))
					return
				}

				cfgChan <- resp
			},
			UpdateFunc: func(_ interface{}, newObj interface{}) {
				u := newObj.(*unstructured.Unstructured)

				k.addOrUpdateConfig(u)

				resp, err := k.getCaddyConfig()
				if err != nil {
					k.logger.Error("error reloading config: ", zap.Error(err))
					return
				}

				cfgChan <- resp
			},
			DeleteFunc: func(obj interface{}) {
				u := obj.(*unstructured.Unstructured)

				// new configuration
				key := source.GetModuleName(u.GetAPIVersion(), u.GetKind())
				s := []*unstructured.Unstructured{}
				for _, spec := range k.configurationN[key] {
					if spec.GetName() == u.GetName() {
						continue
					}
					s = append(s, spec)
				}
				k.configurationN[key] = s

				// old configuration
				kind := u.GetKind()
				s = []*unstructured.Unstructured{}
				for _, spec := range k.configuration[kind] {
					if spec.GetName() == u.GetName() {
						continue
					}
					s = append(s, spec)
				}
				k.configuration[kind] = s

				resp, err := k.getCaddyConfig()
				if err != nil {
					k.logger.Error("error reloading config: ", zap.Error(err))
					return
				}

				cfgChan <- resp
			},
		})

		go informer.Run(ctx.Done())
	}

	return cfgChan, nil
}

func (k *K8s) getCaddyConfig() ([]byte, error) {
	// Load the new caddy config
	config, err := common.PrepareConfig(k.configuration, k.configurationN)
	if err != nil {
		return nil, err
	}

	return json.MarshalIndent(config, "", "  ")
}

func (k *K8s) addOrUpdateConfig(u *unstructured.Unstructured) {
	// new configuration
	key := source.GetModuleName(u.GetAPIVersion(), u.GetKind())
	found := false
	for i, spec := range k.configurationN[key] {
		if spec.GetName() == u.GetName() {
			k.configurationN[key][i] = u
			found = true
			break
		}
	}
	if !found {
		k.configurationN[key] = append(k.configurationN[key], u)
	}

	// old configuration
	kind := u.GetKind()
	found = false
	for i, spec := range k.configuration[kind] {
		if spec.GetName() == u.GetName() {
			k.configuration[kind][i] = u
			found = true
			break
		}
	}
	if !found {
		k.configuration[kind] = append(k.configuration[kind], u)
	}
}

// Interface guard
var (
	_ adapter.Adapter = (*K8s)(nil)
)
