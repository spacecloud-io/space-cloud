package istio

import (
	"context"

	"istio.io/client-go/pkg/apis/networking/v1alpha3"
	"istio.io/client-go/pkg/apis/security/v1beta1"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	kubeErrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (i *Istio) createServiceAccountIfNotExist(ctx context.Context, ns string, obj *v1.ServiceAccount) error {
	_, err := i.kube.CoreV1().ServiceAccounts(ns).Get(ctx, obj.Name, metav1.GetOptions{})
	if kubeErrors.IsNotFound(err) {
		// Create a service account if it doesn't already exist
		_, err = i.kube.CoreV1().ServiceAccounts(ns).Create(ctx, obj, metav1.CreateOptions{})
		return err
	}
	return err
}

func (i *Istio) applyDeployment(ctx context.Context, ns string, deployment *appsv1.Deployment) error {
	prevDeployment, err := i.kube.AppsV1().Deployments(ns).Get(ctx, deployment.Name, metav1.GetOptions{})
	if kubeErrors.IsNotFound(err) {
		// Create a deployment if it doesn't already exist
		if _, err := i.kube.AppsV1().Deployments(ns).Create(ctx, deployment, metav1.CreateOptions{}); err != nil {
			return err
		}
		return i.cache.setDeployment(ns, deployment.Name, deployment)
	}
	if err != nil {
		return err
	}

	// Update the deployment config
	prevDeployment.Labels = deployment.Labels
	prevDeployment.Annotations = deployment.Annotations
	prevDeployment.Spec = deployment.Spec
	if _, err := i.kube.AppsV1().Deployments(ns).Update(ctx, prevDeployment, metav1.UpdateOptions{}); err != nil {
		return err
	}
	// Update the deployment cache
	if err := i.cache.setDeployment(ns, deployment.Name, prevDeployment); err != nil {
		return err
	}

	return nil
}

func (i *Istio) createServiceIfNotExist(ctx context.Context, ns string, service *v1.Service) error {
	_, err := i.kube.CoreV1().Services(ns).Get(ctx, service.Name, metav1.GetOptions{})
	if kubeErrors.IsNotFound(err) {
		// Create a service if it doesn't already exist
		_, err = i.kube.CoreV1().Services(ns).Create(ctx, service, metav1.CreateOptions{})
		return err
	}
	return err
}

func (i *Istio) applyService(ctx context.Context, ns string, service *v1.Service) error {
	prevService, err := i.kube.CoreV1().Services(ns).Get(ctx, service.Name, metav1.GetOptions{})
	if kubeErrors.IsNotFound(err) {
		// Create a service if it doesn't already exist
		_, err = i.kube.CoreV1().Services(ns).Create(ctx, service, metav1.CreateOptions{})
		return err
	}
	if err != nil {
		return err
	}

	// Update the service
	prevService.Spec.Ports = service.Spec.Ports
	prevService.Labels = service.Labels
	_, err = i.kube.CoreV1().Services(ns).Update(ctx, prevService, metav1.UpdateOptions{})
	return err
}

func (i *Istio) createVirtualServiceIfNotExist(ctx context.Context, ns string, service *v1alpha3.VirtualService) error {
	_, err := i.istio.NetworkingV1alpha3().VirtualServices(ns).Get(ctx, service.Name, metav1.GetOptions{})
	if kubeErrors.IsNotFound(err) {
		// Create a service if it doesn't already exist
		_, err := i.istio.NetworkingV1alpha3().VirtualServices(ns).Create(ctx, service, metav1.CreateOptions{})
		return err
	}
	return err
}

func (i *Istio) applyVirtualService(ctx context.Context, ns string, service *v1alpha3.VirtualService) error {
	prevVirtualService, err := i.istio.NetworkingV1alpha3().VirtualServices(ns).Get(ctx, service.Name, metav1.GetOptions{})
	if kubeErrors.IsNotFound(err) {
		// Create a service if it doesn't already exist
		_, err := i.istio.NetworkingV1alpha3().VirtualServices(ns).Create(ctx, service, metav1.CreateOptions{})
		return err
	}
	if err != nil {
		return err
	}

	// Update the service
	prevVirtualService.Spec = service.Spec
	prevVirtualService.Labels = service.Labels
	_, err = i.istio.NetworkingV1alpha3().VirtualServices(ns).Update(ctx, prevVirtualService, metav1.UpdateOptions{})
	return err
}

func (i *Istio) createDestinationRulesIfNotExist(ctx context.Context, ns string, rule *v1alpha3.DestinationRule) error {
	_, err := i.istio.NetworkingV1alpha3().DestinationRules(ns).Get(ctx, rule.Name, metav1.GetOptions{})
	if kubeErrors.IsNotFound(err) {
		// Create a dest rule if it doesn't already exist
		_, err := i.istio.NetworkingV1alpha3().DestinationRules(ns).Create(ctx, rule, metav1.CreateOptions{})
		return err
	}
	return err
}

func (i *Istio) applyDestinationRules(ctx context.Context, ns string, rule *v1alpha3.DestinationRule) error {
	prevRule, err := i.istio.NetworkingV1alpha3().DestinationRules(ns).Get(ctx, rule.Name, metav1.GetOptions{})
	if kubeErrors.IsNotFound(err) {
		// Create a destination rule if it doesn't already exist
		_, err := i.istio.NetworkingV1alpha3().DestinationRules(ns).Create(ctx, rule, metav1.CreateOptions{})
		return err
	}
	if err != nil {
		return err
	}

	// Update the destination rules
	prevRule.Spec = rule.Spec
	prevRule.Labels = rule.Labels
	_, err = i.istio.NetworkingV1alpha3().DestinationRules(ns).Update(ctx, prevRule, metav1.UpdateOptions{})
	return err
}

func (i *Istio) applyAuthorizationPolicy(ctx context.Context, ns string, policy *v1beta1.AuthorizationPolicy) error {
	prevPolicy, err := i.istio.SecurityV1beta1().AuthorizationPolicies(ns).Get(ctx, policy.Name, metav1.GetOptions{})
	if kubeErrors.IsNotFound(err) {
		// Create a authorization policy if it doesn't already exist
		_, err := i.istio.SecurityV1beta1().AuthorizationPolicies(ns).Create(ctx, policy, metav1.CreateOptions{})
		return err
	}
	if err != nil {
		return err
	}

	// Update the authorization policy
	prevPolicy.Spec = policy.Spec
	prevPolicy.Labels = policy.Labels
	_, err = i.istio.SecurityV1beta1().AuthorizationPolicies(ns).Update(ctx, prevPolicy, metav1.UpdateOptions{})
	return err
}

func (i *Istio) applySidecar(ctx context.Context, ns string, sidecar *v1alpha3.Sidecar) error {
	prevSidecar, err := i.istio.NetworkingV1alpha3().Sidecars(ns).Get(ctx, sidecar.Name, metav1.GetOptions{})
	if kubeErrors.IsNotFound(err) {
		// Create a sidecar config if it doesn't already exist
		_, err := i.istio.NetworkingV1alpha3().Sidecars(ns).Create(ctx, sidecar, metav1.CreateOptions{})
		return err
	}
	if err != nil {
		return err
	}

	// Update the sidecar config
	prevSidecar.Spec = sidecar.Spec
	prevSidecar.Labels = sidecar.Labels
	_, err = i.istio.NetworkingV1alpha3().Sidecars(ns).Update(ctx, prevSidecar, metav1.UpdateOptions{})
	return err
}
