package istio

import (
	kubeErrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (i *Istio) deleteServiceAccountIfExist(projectID, serviceID string) error {
	err := i.kube.CoreV1().ServiceAccounts(projectID).Delete(getServiceAccountName(serviceID), &metav1.DeleteOptions{})
	return ignoreErrorIfNotFound(err)
}

func (i *Istio) deleteDeployment(projectID, serviceID, version string) error {
	err := i.kube.AppsV1().Deployments(projectID).Delete(getDeploymentName(serviceID, version), &metav1.DeleteOptions{})
	return ignoreErrorIfNotFound(err)
}

func (i *Istio) deleteGeneralService(projectID, serviceID string) error {
	err := i.kube.CoreV1().Services(projectID).Delete(getServiceName(serviceID), &metav1.DeleteOptions{})
	return ignoreErrorIfNotFound(err)
}

func (i *Istio) deleteInternalService(projectID, serviceID, version string) error {
	err := i.kube.CoreV1().Services(projectID).Delete(getInternalServiceName(serviceID, version), &metav1.DeleteOptions{})
	return ignoreErrorIfNotFound(err)
}

func (i *Istio) deleteVirtualService(projectID, serviceID string) error {
	err := i.istio.NetworkingV1alpha3().VirtualServices(projectID).Delete(getVirtualServiceName(serviceID), &metav1.DeleteOptions{})
	return ignoreErrorIfNotFound(err)
}

func (i *Istio) deleteGeneralDestRule(projectID, serviceID string) error {
	err := i.istio.NetworkingV1alpha3().DestinationRules(projectID).Delete(getGeneralDestRuleName(serviceID), &metav1.DeleteOptions{})
	return ignoreErrorIfNotFound(err)
}

func (i *Istio) deleteInternalDestRule(projectID, serviceID, version string) error {
	err := i.istio.NetworkingV1alpha3().DestinationRules(projectID).Delete(getInternalDestRuleName(serviceID, version), &metav1.DeleteOptions{})
	return ignoreErrorIfNotFound(err)
}

func (i *Istio) deleteAuthorizationPolicy(projectID, serviceID, version string) error {
	err := i.istio.SecurityV1beta1().AuthorizationPolicies(projectID).Delete(getAuthorizationPolicyName(projectID, serviceID, version), &metav1.DeleteOptions{})
	return ignoreErrorIfNotFound(err)
}

func (i *Istio) deleteSidecarConfig(projectID, serviceID, version string) error {
	err := i.istio.NetworkingV1alpha3().Sidecars(projectID).Delete(getSidecarName(serviceID, version), &metav1.DeleteOptions{})
	return ignoreErrorIfNotFound(err)
}

func ignoreErrorIfNotFound(err error) error {
	if kubeErrors.IsNotFound(err) {
		return nil
	}
	return err
}
