package istio

import (
	"context"
	"fmt"

	"github.com/spaceuptech/helpers"
)

// DeleteService deletes a service version
func (i *Istio) DeleteService(ctx context.Context, projectID, serviceID, version string) error {
	// Get the count of versions running for this service. This is important to make sure we do not delete shared resources.
	count, err := i.getServiceDeploymentsCount(ctx, projectID, serviceID)
	if err != nil {
		return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Error in delete service - could not get count of versions for service (%s)", getServiceUniqueID(projectID, serviceID, version)), err, nil)
	}

	// TODO: this could turn out to be a problem when two delete requests for the same service come in simultaneously
	if count == 1 {
		if err := i.DeleteServiceRole(ctx, projectID, serviceID, "*"); err != nil {
			return helpers.Logger.LogError(helpers.GetRequestID(ctx), "Could not delete service - service role could not be deleted", err, nil)
		}
		if err := i.deleteServiceAccountIfExist(ctx, projectID, serviceID); err != nil {
			return helpers.Logger.LogError(helpers.GetRequestID(ctx), "Could not delete service - service account could not be deleted", err, nil)
		}
		if err := i.deleteGeneralService(ctx, projectID, serviceID); err != nil {
			return helpers.Logger.LogError(helpers.GetRequestID(ctx), "Could not delete service - general service could not be deleted", err, nil)
		}
		if err := i.deleteGeneralDestRule(ctx, projectID, serviceID); err != nil {
			return helpers.Logger.LogError(helpers.GetRequestID(ctx), "Could not delete service - general destination rule could not be deleted", err, nil)
		}
		if err := i.deleteVirtualService(ctx, projectID, serviceID); err != nil {
			return helpers.Logger.LogError(helpers.GetRequestID(ctx), "Could not delete service - virtual service could not be deleted", err, nil)
		}
	}

	if err := i.deleteDeployment(ctx, projectID, serviceID, version); err != nil {
		return helpers.Logger.LogError(helpers.GetRequestID(ctx), "Could not delete service - deployment could not be deleted", err, nil)
	}
	if err := i.deleteInternalService(ctx, projectID, serviceID, version); err != nil {
		return helpers.Logger.LogError(helpers.GetRequestID(ctx), "Could not delete service - internal service could not be deleted", err, nil)
	}
	if err := i.deleteInternalDestRule(ctx, projectID, serviceID, version); err != nil {
		return helpers.Logger.LogError(helpers.GetRequestID(ctx), "Could not delete service - internal destination rule could not be deleted", err, nil)
	}
	if err := i.deleteAuthorizationPolicy(ctx, projectID, serviceID, version); err != nil {
		return helpers.Logger.LogError(helpers.GetRequestID(ctx), "Could not delete service - authorization policy could not be deleted", err, nil)
	}
	if err := i.deleteSidecarConfig(ctx, projectID, serviceID, version); err != nil {
		return helpers.Logger.LogError(helpers.GetRequestID(ctx), "Could not delete service - sidecar config could not be deleted", err, nil)
	}
	if err := i.deleteKedaConfig(ctx, projectID, serviceID, version); err != nil {
		return helpers.Logger.LogError(helpers.GetRequestID(ctx), "Could not delete service - autoscaler config could not be deleted", err, nil)
	}

	return nil
}

// DeleteServiceRole deletes a service role
func (i *Istio) DeleteServiceRole(ctx context.Context, projectID, serviceID, id string) error {
	if err := i.deleteServiceRoleIfExist(ctx, projectID, serviceID, id); err != nil {
		return err
	}
	return nil
}
