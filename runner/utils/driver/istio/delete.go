package istio

import (
	"context"

	"github.com/sirupsen/logrus"
)

// DeleteService deletes a service version
func (i *Istio) DeleteService(ctx context.Context, projectID, serviceID, version string) error {
	// Get the count of versions running for this service. This is important to make sure we do not delete shared resources.
	count, err := i.getServiceDeploymentsCount(ctx, projectID, serviceID)
	if err != nil {
		logrus.Debugf("Error in delete service - could not get count of versions for service (%s) - %s", getServiceUniqueID(projectID, serviceID, version), err.Error())
		return err
	}

	// TODO: this could turn out to be a problem when two delete requests come in simultaneously
	if count == 1 {
		if err := i.deleteServiceAccountIfExist(ctx, projectID, serviceID); err != nil {
			logrus.Errorf("Could not delete service - service account could not be deleted - %s", err.Error())
			return err
		}
		if err := i.deleteGeneralService(ctx, projectID, serviceID); err != nil {
			logrus.Errorf("Could not delete service - general service could not be deleted - %s", err.Error())
			return err
		}
		if err := i.deleteGeneralDestRule(ctx, projectID, serviceID); err != nil {
			logrus.Errorf("Could not delete service - general destination rule could not be deleted - %s", err.Error())
			return err
		}
		if err := i.deleteVirtualService(ctx, projectID, serviceID); err != nil {
			logrus.Errorf("Could not delete service - virtual service could not be deleted - %s", err.Error())
			return err
		}
	}

	if err := i.deleteDeployment(ctx, projectID, serviceID, version); err != nil {
		logrus.Errorf("Could not delete service - deployment could not be deleted - %s", err.Error())
		return err
	}
	if err := i.deleteInternalService(ctx, projectID, serviceID, version); err != nil {
		logrus.Errorf("Could not delete service - internal service could not be deleted - %s", err.Error())
		return err
	}
	if err := i.deleteInternalDestRule(ctx, projectID, serviceID, version); err != nil {
		logrus.Errorf("Could not delete service - internal destination rule could not be deleted - %s", err.Error())
		return err
	}
	if err := i.deleteAuthorizationPolicy(ctx, projectID, serviceID, version); err != nil {
		logrus.Errorf("Could not delete service - authorization policy could not be deleted - %s", err.Error())
		return err
	}
	if err := i.deleteSidecarConfig(ctx, projectID, serviceID, version); err != nil {
		logrus.Errorf("Could not delete service - sidecar config could not be deleted - %s", err.Error())
		return err
	}

	return nil
}
