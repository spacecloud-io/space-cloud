package istio

import (
	"context"
	"fmt"

	"github.com/spaceuptech/helpers"
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

func (i *Istio) deleteKedaConfig(ctx context.Context, projectID, serviceID, version string) error {
	// Delete the scaled object first
	if err := i.keda.KedaV1alpha1().ScaledObjects(projectID).Delete(ctx, getKedaScaledObjectName(serviceID, version), metav1.DeleteOptions{}); !kubeErrors.IsNotFound(err) {
		return err
	}

	// Fetch all the keda trigger auth objects for this version
	triggers, err := i.getKedaTriggerAuthsForVersion(ctx, projectID, serviceID, version)
	if err != nil {
		return err
	}

	// Delete each keda trigger auth object
	for _, trigger := range triggers.Items {
		if err := i.keda.KedaV1alpha1().TriggerAuthentications(projectID).Delete(ctx, trigger.Name, metav1.DeleteOptions{}); !kubeErrors.IsNotFound(err) {
			return err
		}
	}

	return nil
}

func (i *Istio) deleteServiceRoleIfExist(ctx context.Context, projectID, serviceID, id string) error {
	if id == "*" {
		rolelist, err := i.kube.RbacV1().Roles(projectID).List(ctx, metav1.ListOptions{LabelSelector: fmt.Sprintf("app.kubernetes.io/managed-by=space-cloud,app=%s", serviceID)})
		if err != nil {
			return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Unable to list roles in project (%s)", projectID), err, nil)
		}
		for _, role := range rolelist.Items {
			err := i.kube.RbacV1().RoleBindings(projectID).Delete(ctx, role.Name, metav1.DeleteOptions{})
			if err != nil {
				return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Unable to delete service role binding with id (%s)", role.Name), nil, nil)
			}
			err = i.kube.RbacV1().Roles(projectID).Delete(ctx, role.Name, metav1.DeleteOptions{})
			if err != nil {
				return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Unable to delete service role with id (%s)", role.Name), nil, nil)
			}
		}

		clusterRoleList, err := i.kube.RbacV1().ClusterRoles().List(ctx, metav1.ListOptions{LabelSelector: fmt.Sprintf("app.kubernetes.io/managed-by=space-cloud,app=%s", serviceID)})
		if err != nil {
			return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Unable to list cluster roles in project (%s)", projectID), err, nil)
		}
		for _, role := range clusterRoleList.Items {
			err := i.kube.RbacV1().ClusterRoleBindings().Delete(ctx, role.Name, metav1.DeleteOptions{})
			if err != nil {
				return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Unable to delete service role binding with id (%s)", role.Name), nil, nil)
			}
			err = i.kube.RbacV1().ClusterRoles().Delete(ctx, role.Name, metav1.DeleteOptions{})
			if err != nil {
				return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Unable to delete service role with id (%s)", role.Name), nil, nil)
			}
		}
		return nil
	}
	_, err1 := i.kube.RbacV1().Roles(projectID).Get(ctx, id, metav1.GetOptions{})
	if !kubeErrors.IsNotFound(err1) {
		err := i.kube.RbacV1().RoleBindings(projectID).Delete(ctx, id, metav1.DeleteOptions{})
		if err != nil {
			return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Unable to delete service role binding with id (%s)", id), nil, nil)
		}
		err = i.kube.RbacV1().Roles(projectID).Delete(ctx, id, metav1.DeleteOptions{})
		if err != nil {
			return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Unable to delete service role with id (%s)", id), nil, nil)
		}
		return nil
	}
	_, err2 := i.kube.RbacV1().ClusterRoles().Get(ctx, id, metav1.GetOptions{})
	if !kubeErrors.IsNotFound(err2) {
		err := i.kube.RbacV1().ClusterRoleBindings().Delete(ctx, id, metav1.DeleteOptions{})
		if err != nil {
			return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Unable to delete service role binding with id (%s)", id), nil, nil)
		}
		err = i.kube.RbacV1().ClusterRoles().Delete(ctx, id, metav1.DeleteOptions{})
		if err != nil {
			return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Unable to delete service role with id (%s)", id), nil, nil)
		}
		return nil
	}

	return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Specified role id (%s) was not present", id), nil, map[string]interface{}{"roleError": err1, "clusterRoleError": err2})
}

func ignoreErrorIfNotFound(err error) error {
	if kubeErrors.IsNotFound(err) {
		return nil
	}
	return err
}
