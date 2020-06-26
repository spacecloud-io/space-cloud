package istio

import (
	"context"

	kubeErrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (i *Istio) deleteServiceAccountIfExist(ctx context.Context, projectID, serviceID string) error {
	err := i.kube.CoreV1().ServiceAccounts(projectID).Delete(ctx, getServiceAccountName(serviceID), metav1.DeleteOptions{})
	return ignoreErrorIfNotFound(err)
}

func (i *Istio) deleteDeployment(ctx context.Context, projectID, serviceID, version string) error {
	err := i.kube.AppsV1().Deployments(projectID).Delete(ctx, getDeploymentName(serviceID, version), metav1.DeleteOptions{})
	return ignoreErrorIfNotFound(err)
}

func (i *Istio) deleteGeneralService(ctx context.Context, projectID, serviceID string) error {
	err := i.kube.CoreV1().Services(projectID).Delete(ctx, getServiceName(serviceID), metav1.DeleteOptions{})
	return ignoreErrorIfNotFound(err)
}

func (i *Istio) deleteInternalService(ctx context.Context, projectID, serviceID, version string) error {
	err := i.kube.CoreV1().Services(projectID).Delete(ctx, getInternalServiceName(serviceID, version), metav1.DeleteOptions{})
	return ignoreErrorIfNotFound(err)
}

func (i *Istio) deleteVirtualService(ctx context.Context, projectID, serviceID string) error {
	err := i.istio.NetworkingV1alpha3().VirtualServices(projectID).Delete(ctx, getVirtualServiceName(serviceID), metav1.DeleteOptions{})
	return ignoreErrorIfNotFound(err)
}

func (i *Istio) deleteGeneralDestRule(ctx context.Context, projectID, serviceID string) error {
	err := i.istio.NetworkingV1alpha3().DestinationRules(projectID).Delete(ctx, getGeneralDestRuleName(serviceID), metav1.DeleteOptions{})
	return ignoreErrorIfNotFound(err)
}

func (i *Istio) deleteInternalDestRule(ctx context.Context, projectID, serviceID, version string) error {
	err := i.istio.NetworkingV1alpha3().DestinationRules(projectID).Delete(ctx, getInternalDestRuleName(serviceID, version), metav1.DeleteOptions{})
	return ignoreErrorIfNotFound(err)
}

func (i *Istio) deleteAuthorizationPolicy(ctx context.Context, projectID, serviceID, version string) error {
	err := i.istio.SecurityV1beta1().AuthorizationPolicies(projectID).Delete(ctx, getAuthorizationPolicyName(projectID, serviceID, version), metav1.DeleteOptions{})
	return ignoreErrorIfNotFound(err)
}

func (i *Istio) deleteSidecarConfig(ctx context.Context, projectID, serviceID, version string) error {
	err := i.istio.NetworkingV1alpha3().Sidecars(projectID).Delete(ctx, getSidecarName(serviceID, version), metav1.DeleteOptions{})
	return ignoreErrorIfNotFound(err)
}

func ignoreErrorIfNotFound(err error) error {
	if kubeErrors.IsNotFound(err) {
		return nil
	}
	return err
}
