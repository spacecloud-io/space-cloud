package istio

import (
	"context"
	"fmt"
	"time"

	"github.com/spaceuptech/helpers"
	"github.com/spaceuptech/space-cloud/runner/model"
)

// WaitForService adjusts scales, up the service to scale up the number of nodes from zero to one
// TODO: Do one watch per service. Right now its possible to have multiple watches for the same service
func (i *Istio) WaitForService(ctx context.Context, service *model.Service) error {

	helpers.Logger.LogDebug(helpers.GetRequestID(ctx), fmt.Sprintf("Scaling up service (%s:%s:%s) from zero", service.ProjectID, service.ID, service.Version), nil)

	helpers.Logger.LogDebug(helpers.GetRequestID(ctx), fmt.Sprintf("Watching for service (%s:%s:%s) to scale up and enter ready state", service.ProjectID, service.ID, service.Version), nil)

	ticker := time.NewTicker(200 * time.Millisecond)

	for {
		select {
		case <-ctx.Done():
			return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Service (%s) could not be started", service.ID), nil, nil)
		case <-ticker.C:
			if i.getStatusOfDeployement(service.ProjectID, service.ID) {
				return nil
			}
		}
	}
}

// ScaleUp is notifies keda to scale up a service
func (i *Istio) ScaleUp(ctx context.Context, projectID, serviceID, version string) error {
	if err := i.kedaScaler.ScaleUp(ctx, projectID, serviceID, version); err != nil {
		return helpers.Logger.LogError(helpers.GetRequestID(ctx), "Unable to scale up service", err, map[string]interface{}{"project": projectID, "service": serviceID, "version": version})
	}

	return nil
}
