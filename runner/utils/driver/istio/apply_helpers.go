package istio

import (
	"context"
	"fmt"

	"github.com/kedacore/keda/api/v1alpha1"
	"github.com/spaceuptech/helpers"
	"istio.io/client-go/pkg/apis/networking/v1alpha3"
	"istio.io/client-go/pkg/apis/security/v1beta1"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	v12 "k8s.io/api/rbac/v1"
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
		_, err := i.kube.AppsV1().Deployments(ns).Create(ctx, deployment, metav1.CreateOptions{})
		return err
	}
	if err != nil {
		return err
	}

	// Update the deployment config
	prevReplicas := prevDeployment.Spec.Replicas
	prevDeployment.Labels = deployment.Labels
	prevDeployment.Annotations = deployment.Annotations
	prevDeployment.Spec = deployment.Spec
	prevDeployment.Spec.Replicas = prevReplicas
	_, err = i.kube.AppsV1().Deployments(ns).Update(ctx, prevDeployment, metav1.UpdateOptions{})
	return err
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
	prevService.Annotations = service.Annotations
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
	prevVirtualService.Annotations = service.Annotations
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
	prevRule.Annotations = rule.Annotations
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
	prevPolicy.Annotations = policy.Annotations
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
	prevSidecar.Annotations = sidecar.Annotations
	prevSidecar.Labels = sidecar.Labels
	_, err = i.istio.NetworkingV1alpha3().Sidecars(ns).Update(ctx, prevSidecar, metav1.UpdateOptions{})
	return err
}

func (i *Istio) applyKedaConfig(ctx context.Context, ns string, scaledObj *v1alpha1.ScaledObject, triggerAuths []v1alpha1.TriggerAuthentication) error {
	// First create all the trigger authentication objects
	for _, triggerAuth := range triggerAuths {
		helpers.Logger.LogDebug(helpers.GetRequestID(ctx), fmt.Sprintf("Applying keda trigger auth config (%s) in %s", triggerAuth.Name, ns), nil)

		prevAuth, err := i.keda.KedaV1alpha1().TriggerAuthentications(ns).Get(ctx, triggerAuth.Name, metav1.GetOptions{})
		if kubeErrors.IsNotFound(err) {
			// Create the trigger authentication
			if _, err := i.keda.KedaV1alpha1().TriggerAuthentications(ns).Create(ctx, &triggerAuth, metav1.CreateOptions{}); err != nil {
				return nil
			}
		} else if err != nil {
			return err
		}

		// Update the trigger authentication object if it already exists
		prevAuth.Spec = triggerAuth.Spec
		prevAuth.Labels = triggerAuth.Labels
		if _, err := i.keda.KedaV1alpha1().TriggerAuthentications(ns).Update(ctx, prevAuth, metav1.UpdateOptions{}); err != nil {
			return err
		}
	}

	// Time to apply the keda scaled object
	helpers.Logger.LogDebug(helpers.GetRequestID(ctx), fmt.Sprintf("Applying keda scaled object config (%s) in %s", scaledObj.Name, ns), nil)
	prevScaledObj, err := i.keda.KedaV1alpha1().ScaledObjects(ns).Get(ctx, scaledObj.Name, metav1.GetOptions{})
	if kubeErrors.IsNotFound(err) {
		// Create the keda scaled object
		_, err = i.keda.KedaV1alpha1().ScaledObjects(ns).Create(ctx, scaledObj, metav1.CreateOptions{})
		return err
	} else if err != nil {
		return err
	}

	// Update the keda scaled object if it already exists
	prevScaledObj.Spec = scaledObj.Spec
	prevScaledObj.Labels = scaledObj.Labels
	_, err = i.keda.KedaV1alpha1().ScaledObjects(ns).Update(ctx, prevScaledObj, metav1.UpdateOptions{})
	return err
}

func (i *Istio) applyServiceRole(ctx context.Context, ns string, role *v12.Role) error {
	prevServiceRole, err := i.kube.RbacV1().Roles(ns).Get(ctx, role.Name, metav1.GetOptions{})
	if kubeErrors.IsNotFound(err) {
		// Create a service role if it doesn't already exist
		_, err := i.kube.RbacV1().Roles(ns).Create(ctx, role, metav1.CreateOptions{})
		return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Unable to create service role in project (%s)", ns), err, nil)
	}
	if err != nil {
		return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Unable to get service role in project (%s)", ns), err, nil)
	}

	// Update the service role
	prevServiceRole.Labels = role.Labels
	prevServiceRole.Rules = role.Rules
	_, err = i.kube.RbacV1().Roles(ns).Update(ctx, prevServiceRole, metav1.UpdateOptions{})
	return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Unable to update service role in project (%s)", ns), err, nil)
}

func (i *Istio) applyServiceRoleBinding(ctx context.Context, ns string, rolebinding *v12.RoleBinding) error {
	prevServiceRoleBinding, err := i.kube.RbacV1().RoleBindings(ns).Get(ctx, rolebinding.Name, metav1.GetOptions{})
	if kubeErrors.IsNotFound(err) {
		// Create a service role binding if it doesn't already exist
		_, err := i.kube.RbacV1().RoleBindings(ns).Create(ctx, rolebinding, metav1.CreateOptions{})
		return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Unable to create service role binding in project (%s)", ns), err, nil)
	}
	if err != nil {
		return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Unable to get service role binding in project (%s)", ns), err, nil)
	}

	// Update the service role binding
	prevServiceRoleBinding.Labels = rolebinding.Labels
	prevServiceRoleBinding.Subjects = rolebinding.Subjects
	prevServiceRoleBinding.RoleRef = rolebinding.RoleRef
	_, err = i.kube.RbacV1().RoleBindings(ns).Update(ctx, prevServiceRoleBinding, metav1.UpdateOptions{})
	return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Unable to update service role binding in project (%s)", ns), err, nil)
}

func (i *Istio) applyServiceClusterRole(ctx context.Context, ns string, clusterRole *v12.ClusterRole) error {
	prevServiceClusterRole, err := i.kube.RbacV1().ClusterRoles().Get(ctx, clusterRole.Name, metav1.GetOptions{})
	if kubeErrors.IsNotFound(err) {
		// Create a service clusterRole if it doesn't already exist
		_, err := i.kube.RbacV1().ClusterRoles().Create(ctx, clusterRole, metav1.CreateOptions{})
		return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Unable to create service cluster role in project (%s)", ns), err, nil)
	}
	if err != nil {
		return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Unable to get service cluster role in project (%s)", ns), err, nil)
	}

	// Update the service clusterRole
	prevServiceClusterRole.Labels = clusterRole.Labels
	prevServiceClusterRole.Rules = clusterRole.Rules
	_, err = i.kube.RbacV1().ClusterRoles().Update(ctx, prevServiceClusterRole, metav1.UpdateOptions{})
	return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Unable to update service cluster role in project (%s)", ns), err, nil)
}

func (i *Istio) applyServiceClusterRoleBinding(ctx context.Context, ns string, clusterRoleBinding *v12.ClusterRoleBinding) error {
	prevServiceClusterRoleBinding, err := i.kube.RbacV1().ClusterRoleBindings().Get(ctx, clusterRoleBinding.Name, metav1.GetOptions{})
	if kubeErrors.IsNotFound(err) {
		// Create a service cluster role binding if it doesn't already exist
		_, err := i.kube.RbacV1().ClusterRoleBindings().Create(ctx, clusterRoleBinding, metav1.CreateOptions{})
		return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Unable to create service cluster role binding in project (%s)", ns), err, nil)
	}
	if err != nil {
		return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Unable to get service cluster role binding in project (%s)", ns), err, nil)
	}

	// Update the service cluster role binding
	prevServiceClusterRoleBinding.Labels = clusterRoleBinding.Labels
	prevServiceClusterRoleBinding.Subjects = clusterRoleBinding.Subjects
	prevServiceClusterRoleBinding.RoleRef = clusterRoleBinding.RoleRef
	_, err = i.kube.RbacV1().ClusterRoleBindings().Update(ctx, prevServiceClusterRoleBinding, metav1.UpdateOptions{})
	return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Unable to update service cluster role binding in project (%s)", ns), err, nil)
}
