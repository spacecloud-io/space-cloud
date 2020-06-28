package istio

import (
	"context"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/spaceuptech/space-cloud/runner/model"
)

// CreateProject creates a new namespace for the client
func (i *Istio) CreateProject(ctx context.Context, project *model.Project) error {
	// Project ID provided here is already in the form `project-env`
	namespace := project.ID
	ns := &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:   namespace,
			Labels: map[string]string{"istio-injection": "enabled"},
		},
	}
	_, err := i.kube.CoreV1().Namespaces().Create(ctx, ns, metav1.CreateOptions{})
	return err
}

// DeleteProject deletes a namespace for the client
func (i *Istio) DeleteProject(ctx context.Context, projectID string) error {
	return i.kube.CoreV1().Namespaces().Delete(ctx, projectID, metav1.DeleteOptions{})
}
