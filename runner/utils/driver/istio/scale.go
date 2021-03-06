package istio

import (
	"context"
	"fmt"
	"time"

	"github.com/spaceuptech/helpers"
	"github.com/spaceuptech/space-cloud/runner/model"
)

//Deployment stores the deploymentID
type Deployment struct {
	DeployemtID map[string]Replicas
}

//Replicas stores the value of AvailableReplicas and ReadyReplicas
type Replicas struct {
	AvailableReplicas int32
	ReadyReplicas     int32
}

// WaitForService adjusts scales, up the service to scale up the number of nodes from zero to one
// TODO: Do one watch per service. Right now its possible to have multiple watches for the same service
func (i *Istio) WaitForService(ctx context.Context, service *model.Service) error {
	ns := service.ProjectID
	helpers.Logger.LogDebug(helpers.GetRequestID(ctx), fmt.Sprintf("Scaling up service (%s:%s:%s) from zero", ns, service.ID, service.Version), nil)

	helpers.Logger.LogDebug(helpers.GetRequestID(ctx), fmt.Sprintf("Watching for service (%s:%s:%s) to scale up and enter ready state", ns, service.ID, service.Version), nil)

	ticker := time.NewTicker(200 * time.Millisecond)
	done := make(chan bool)

	go func() {
		for {
			select {
			case <-done:
				return
			case _ = <-ticker.C:
				if (*i.waitservice)[service.ProjectID].DeployemtID[service.ID].AvailableReplicas >= 1 && (*i.waitservice)[service.ProjectID].DeployemtID[service.ID].ReadyReplicas >= 1 {
					return
				}
			}
		}
	}()
	time.Sleep(5 * time.Second)
	ticker.Stop()
	done <- true

	return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("service (%s:%s) could not be started", ns, service.ID), nil, nil)
}

// ScaleUp is notifies keda to scale up a service
func (i *Istio) ScaleUp(ctx context.Context, projectID, serviceID, version string) error {
	if err := i.kedaScaler.ScaleUp(ctx, projectID, serviceID, version); err != nil {
		return helpers.Logger.LogError(helpers.GetRequestID(ctx), "Unable to scale up service", err, map[string]interface{}{"project": projectID, "service": serviceID, "version": version})
	}

	return nil
}
