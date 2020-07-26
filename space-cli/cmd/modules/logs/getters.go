package logs

import (
	"fmt"
	"github.com/spaceuptech/space-cloud/space-cli/cmd/utils/transport"
	"net/http"

	"github.com/spaceuptech/space-cloud/space-cli/cmd/utils"
)

// GetServiceLogs gets logs of specified service
func GetServiceLogs(project, serviceID, taskID, replicaID string) error {
	if serviceID == "" {
		_ = utils.LogError("Service id not specified in flag", nil)
		return nil
	}

	if taskID == "" {
		_ = utils.LogError("Task id not specified in flag", nil)
		return nil
	}

	if replicaID == "" {
		_ = utils.LogError("Replica id not specified in flag", nil)
		return nil
	}

	url := fmt.Sprintf("/v1/runner/%s/services/logs/%s/%s/%s", project, serviceID, taskID, replicaID)
	if err := transport.Client.GetLogs(http.MethodGet, url); err != nil {
		return err
	}
	return nil
}
