package kubernetes

import (
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/spaceuptech/space-cloud/model"
)

func generateService(c *model.Deploy) *apiv1.Service {
	return &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: generateServiceName(c.Name),
		},
		Spec: apiv1.ServiceSpec{
			Selector: map[string]string{"app": c.Name},
			Ports:    generateServicePorts(c.Ports),
		},
	}
}

func (d *Driver) generateDeployment(c *model.Deploy) *appsv1.Deployment {
	// The revision history limit
	revisionHistoryLimit := int32(1)

	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:   generateDeploymentName(c.Name),
			Labels: map[string]string{"app": c.Name},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &c.Constraints.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": c.Name},
			},
			RevisionHistoryLimit: &revisionHistoryLimit,
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      map[string]string{"app": c.Name},
					Annotations: map[string]string{"count": "1"},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:            c.Name,
							Image:           c.Runtime.Name,
							Ports:           generatePorts(c.Ports),
							ImagePullPolicy: apiv1.PullAlways,
							Resources:       *generateResourceRequirements(c.Constraints),
							Env:             d.generateEnv(c),
						},
					},
				},
			},
		},
	}
}
func generateServiceName(name string) string {
	return fmt.Sprintf("%s", name)
}

func generateDeploymentName(name string) string {
	return fmt.Sprintf("%s", name)
}

func generateServicePorts(ports []*model.Port) []apiv1.ServicePort {
	if ports == nil {
		return nil
	}
	array := make([]apiv1.ServicePort, len(ports))
	for i, v := range ports {
		// Create a container port object
		p := apiv1.ServicePort{Port: int32(v.Port)}

		// Add the protocol if it exists
		if v.Protocol != nil {
			p.Protocol = apiv1.Protocol(*v.Protocol)
		}

		// Add name if it exists
		if v.Name != nil {
			p.Name = *v.Name
		}

		array[i] = p
	}

	return array
}

func generatePorts(ports []*model.Port) []apiv1.ContainerPort {
	if ports == nil {
		return []apiv1.ContainerPort{}
	}

	array := make([]apiv1.ContainerPort, len(ports))
	for i, v := range ports {
		// Create a container port object
		p := apiv1.ContainerPort{ContainerPort: int32(v.Port)}

		// Add the protocol if it exists
		if v.Protocol != nil {
			p.Protocol = apiv1.Protocol(*v.Protocol)
		}

		// Add name if it exists
		if v.Name != nil {
			p.Name = *v.Name
		}

		array[i] = p
	}

	return array
}

func generateResourceRequirements(c *model.Constraints) *apiv1.ResourceRequirements {
	if c == nil {
		return nil
	}

	resources := apiv1.ResourceRequirements{Limits: apiv1.ResourceList{}}

	// Set the cpu contraint
	if c.CPU != nil {
		resources.Limits[apiv1.ResourceCPU] = *resource.NewMilliQuantity(int64(*c.CPU*1000), resource.DecimalSI)
	}

	// Set the memory contraint
	if c.Memory != nil {
		resources.Limits[apiv1.ResourceMemory] = *resource.NewQuantity(*c.Memory*1024*1024, resource.BinarySI)
	}

	return &resources
}

func (d *Driver) generateEnv(c *model.Deploy) []apiv1.EnvVar {

	// Create basic env
	array := make([]apiv1.EnvVar, 6)
	array[0] = apiv1.EnvVar{Name: "PROJECT", Value: c.Project}
	array[1] = apiv1.EnvVar{Name: "NAME", Value: c.Name}
	array[2] = apiv1.EnvVar{Name: "REGISTRY_URL", Value: d.registry.URL}
	array[3] = apiv1.EnvVar{Name: "REGISTRY_TOKEN", Value: *d.registry.Token}
	array[4] = apiv1.EnvVar{Name: "CMD_INSTALL", Value: c.Runtime.Install}
	array[5] = apiv1.EnvVar{Name: "CMD_RUN", Value: c.Runtime.Run}

	// Return if there are no use provided env
	if c.Env == nil {
		return array
	}

	for k, v := range c.Env {
		array = append(array, apiv1.EnvVar{Name: k, Value: v})
	}

	return array
}
