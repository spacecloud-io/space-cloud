package istio

import (
	kedaVersionedClient "github.com/kedacore/keda/pkg/generated/clientset/versioned"
	versionedclient "istio.io/client-go/pkg/clientset/versioned"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/spaceuptech/space-cloud/runner/model"
	"github.com/spaceuptech/space-cloud/runner/modules/scaler"

	"github.com/spaceuptech/space-cloud/runner/utils/auth"
)

// Istio manages the istio on kubernetes deployment target
type Istio struct {
	// For internal use
	auth   *auth.Module
	config *Config

	// Drivers to talk to k8s and istio
	kube       kubernetes.Interface
	istio      *versionedclient.Clientset
	keda       *kedaVersionedClient.Clientset
	kedaScaler *scaler.Scaler
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

	kedaClient, err := kedaVersionedClient.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}

	kedaScaler, err := scaler.New(c.PrometheusAddr)
	if err != nil {
		return nil, err
	}

	// Start the keda external scaler
	go kedaScaler.Start()

	return &Istio{auth: auth, config: c, kube: kube, istio: istio, keda: kedaClient, kedaScaler: kedaScaler}, nil
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
