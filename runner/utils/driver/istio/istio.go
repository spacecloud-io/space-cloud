package istio

import (
	"sync"

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

// deployment stores the deploymentID
type deployments map[string]status

// status stores the value of AvailableReplicas and ReadyReplicas
type status struct {
	AvailableReplicas int32
	ReadyReplicas     int32
}

// Istio manages the istio on kubernetes deployment target
type Istio struct {
	// For internal use
	auth         *auth.Module
	config       *Config
	seviceStatus map[string]deployments
	lock         sync.RWMutex

	// Drivers to talk to k8s and istio
	kube       kubernetes.Interface
	istio      versionedclient.Interface
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

	i := &Istio{auth: auth, config: c, seviceStatus: make(map[string]deployments), kube: kube, istio: istio, keda: kedaClient, kedaScaler: kedaScaler}
	if err := i.WatchDeployments(func(eventType string, availableReplicas, readyReplicas int32, projectID, deploymentID string) {
		i.lock.Lock()
		defer i.lock.Unlock()

		switch eventType {
		case resourceAddEvent, resourceUpdateEvent:
			if i.seviceStatus[projectID] == nil {
				i.seviceStatus[projectID] = deployments{
					deploymentID: {
						AvailableReplicas: availableReplicas,
						ReadyReplicas:     readyReplicas,
					},
				}
			} else {
				i.seviceStatus[projectID][deploymentID] = status{
					AvailableReplicas: availableReplicas,
					ReadyReplicas:     readyReplicas,
				}
			}
		case resourceDeleteEvent:
			deployments, ok := i.seviceStatus[projectID]
			if ok {
				_, found := deployments[deploymentID]
				if found {
					delete(deployments, deploymentID)
				}
				if len(deployments) == 0 {
					delete(i.seviceStatus, projectID)
				}
			}
		}
	}); err != nil {
		return nil, err
	}

	return i, nil
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

func (i *Istio) getStatusOfDeployement(projectID, deployementID string) bool {
	i.lock.RLock()
	if deployments, ok := i.seviceStatus[projectID]; ok {
		if status, ok := deployments[deployementID]; ok {
			i.lock.RUnlock()
			return status.AvailableReplicas >= 1 && status.ReadyReplicas >= 1
		}
	}
	i.lock.RUnlock()
	return false
}
