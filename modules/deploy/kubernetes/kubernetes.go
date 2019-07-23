package kubernetes

import (
	"context"
	"strconv"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	kubeErrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/spaceuptech/space-cloud/config"
	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/utils/projects"
)

// Driver is the main kubernetes driver
type Driver struct {
	client   *kubernetes.Clientset
	registry *config.Registry
}

// New creates a new instance of the kubernetes driver
func New(registry *config.Registry) (*Driver, error) {
	d := &Driver{}
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	d.client = clientset
	d.registry = registry
	return d, nil
}

// Deploy deploys the service in kubernetes
func (d *Driver) Deploy(ctx context.Context, c *model.Deploy, projects *projects.Projects) error {
	// Create deployment and service objects
	deployment := d.generateDeployment(c)
	service := generateService(c)

	// Create deployment and services client
	deploymentsClient := d.client.AppsV1().Deployments(apiv1.NamespaceDefault)
	servicesClient := d.client.CoreV1().Services(apiv1.NamespaceDefault)

	// Attemp to get the deployment
	if prevDeployment, err := deploymentsClient.Get(generateDeploymentName(c.Name), metav1.GetOptions{}); err == nil {
		// Update the count
		count, _ := strconv.Atoi(prevDeployment.Spec.Template.ObjectMeta.Annotations["count"])
		deployment.Spec.Template.ObjectMeta.Annotations["count"] = strconv.Itoa(count + 1)

		// Update the deployment if already exists
		_, err = deploymentsClient.Update(deployment)
		if err != nil {
			return err
		}

		// Create service if ports is present
		if c.Ports != nil {
			_, err = servicesClient.Get(c.Name, metav1.GetOptions{})
			if kubeErrors.IsNotFound(err) {
				servicesClient.Create(service)
			} else {
				servicesClient.Update(service)
			}

			// expose the ports if required
			if err := exposeRoute(c, projects); err != nil {
				return err
			}
		}
	} else if kubeErrors.IsNotFound(err) {

		// Create a new deployment if it does not exist
		if _, err = deploymentsClient.Create(deployment); err == nil {

			// Create a service as well if the ports are defined
			if c.Ports != nil {
				_, err = servicesClient.Create(service)
			}

			if err := exposeRoute(c, projects); err != nil {
				// Delete the deployment and service on error
				deploymentsClient.Delete(c.Name, &metav1.DeleteOptions{})
				servicesClient.Delete(c.Name, &metav1.DeleteOptions{})
				return err
			}
		}
	} else {
		return err
	}

	return nil
}

func exposeRoute(c *model.Deploy, projects *projects.Projects) error {
	// If expose param is present
	if c.Expose != nil {
		p, err := projects.LoadProject(c.Project)
		if err != nil {
			return err
		}

		// Remove all the previouse routes
		p.Static.DeleteRoutesWithID(c.Name)

		// Iterate over all the routes to be exposed
		for _, e := range c.Expose {

			// Add the routes exposed
			p.Static.AddProxyRoute(c.Name, *e.Host, *e.Prefix, *e.Proxy)
		}
	}

	return nil
}
