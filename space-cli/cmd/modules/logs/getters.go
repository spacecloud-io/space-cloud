package logs

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/spaceuptech/space-cloud/space-cli/cmd/utils/transport"
)

// GetServiceLogs gets logs of specified service
func GetServiceLogs(project, taskID, replicaID string, isFollow bool) error {
	u, err := url.Parse(fmt.Sprintf("/v1/runner/%s/services/logs", project))
	if err != nil {
		return err
	}
	params := u.Query()
	if isFollow {
		params.Set("follow", "true")
	}
	params.Set("taskId", taskID)
	params.Set("replicaId", replicaID)
	u.RawQuery = params.Encode()
	if err := transport.Client.GetLogs(u.RequestURI()); err != nil {
		return err
	}
	return nil
}

func getServiceStatus(project, commandName string, params map[string]string) ([]string, error) {
	//ReplicaInfo describes structure of replica info
	type ReplicaInfo struct {
		ID     string `json:"id" yaml:"id"`
		Status string `json:"status" yaml:"status"`
	}
	type ServiceStatus struct {
		ServiceID       string         `json:"serviceId" yaml:"serviceId"`
		Version         string         `json:"version" yaml:"version"`
		DesiredReplicas interface{}    `json:"desiredReplicas" yaml:"desiredReplicas"`
		Replicas        []*ReplicaInfo `json:"replicas" yaml:"replicas"`
	}
	type temp struct {
		Error  string           `json:"error,omitempty"`
		Result []*ServiceStatus `json:"result,omitempty"`
	}
	payload := new(temp)
	if err := transport.Client.Get(http.MethodGet, fmt.Sprintf("/v1/runner/%s/services/status", project), params, payload); err != nil {
		return nil, err
	}
	replicaIDs := make([]string, 0)
	for _, serviceStatus := range payload.Result {
		for _, replica := range serviceStatus.Replicas {
			replicaIDs = append(replicaIDs, replica.ID)
		}
	}
	return replicaIDs, nil
}
