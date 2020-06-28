package istio

import (
	"sync"

	versionedclient "istio.io/client-go/pkg/clientset/versioned"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/spaceuptech/space-cloud/runner/model"

	"github.com/spaceuptech/space-cloud/runner/utils/auth"
)

// Istio manages the istio on kubernetes deployment target
type Istio struct {
	// For internal use
	auth   *auth.Module
	config *Config

	// For tacking invocations to adjust scale
	adjustScaleLock sync.Map

	// Drivers to talk to k8s and istio
	kube  *kubernetes.Clientset
	istio *versionedclient.Clientset

	// For caching deployments
	cache *cache
}

// NewIstioDriver creates a new instance of the istio driver
func NewIstioDriver(auth *auth.Module, c *Config) (*Istio, error) {
	var restConfig *rest.Config
	var err error

	if c.IsInsideCluster {
		restConfig, err = rest.InClusterConfig()
	} else {
		restConfig, err = clientcmd.BuildConfigFromFlags("", c.KubeConfigPath)
	}
	if err != nil {
		return nil, err
	}

	// Create the kubernetes client
	kube, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}

	// Create the istio client
	istio, err := versionedclient.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}

	// Create a cache
	cache, err := newCache(kube)
	if err != nil {
		return nil, err
	}

	return &Istio{auth: auth, config: c, kube: kube, istio: istio, cache: cache}, nil
}

func checkIfVolumeIsSecret(name string, volumes []v1.Volume) bool {
	for _, v := range volumes {
		if v.Name == name {
			return true
		}
	}
	return false
}

// Type returns the type of the driver
func (i *Istio) Type() model.DriverType {
	return model.TypeIstio
}
