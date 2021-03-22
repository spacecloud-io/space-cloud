package istio

import (
	"context"

	v1 "k8s.io/api/core/v1"
	kubeErrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/spaceuptech/space-cloud/runner/model"
)

// CreateProject creates a new namespace for the client
func (i *Istio) CreateProject(ctx context.Context, project *model.Project) error {
	// Set the kind field if empty
	if project.Kind == "" {
		project.Kind = "project"
	}

	namespace := project.ID
	ns := &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: namespace,
			Labels: map[string]string{
				"istio-injection":              "enabled",
				"app.kubernetes.io/name":       namespace,
				"app.kubernetes.io/managed-by": "space-cloud",
				"space-cloud.io/kind":          project.Kind,
			},
		},
	}
	_, err := i.kube.CoreV1().Namespaces().Create(ctx, ns, metav1.CreateOptions{})
	if kubeErrors.IsAlreadyExists(err) {
		return nil
	}
	return err
}

// DeleteProject deletes a namespace for the client
func (i *Istio) DeleteProject(ctx context.Context, projectID string) error {
	err := i.kube.CoreV1().Namespaces().Delete(ctx, projectID, metav1.DeleteOptions{})
	if kubeErrors.IsNotFound(err) {
		return nil
	}
	return err
}
