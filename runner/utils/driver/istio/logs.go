package istio

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"strings"

	"k8s.io/api/core/v1"

	"github.com/spaceuptech/space-cloud/runner/utils"
)

// GetLogs get logs of specified services
func (i *Istio) GetLogs(ctx context.Context, isFollow bool, projectID, taskID, replica string) (io.ReadCloser, error) {
	if taskID == "" {
		arr := strings.Split(replica, "--")
		if len(arr) < 2 {
			return nil, utils.LogError("Invalid replica id", "docker", "get-logs", nil)
		}
		taskID = arr[0]
	}
	// get logs of pods
	req := i.kube.CoreV1().Pods(projectID).GetLogs(replica, &v1.PodLogOptions{
		Container:  taskID,
		Follow:     isFollow,
		Timestamps: true,
	})

	b, err := req.Stream(ctx)
	if err != nil {
		return nil, err
	}

	pipeReader, pipeWriter := io.Pipe()
	utils.LogDebug("Sending logs to client", "docker", "GetLogs", map[string]interface{}{})
	go func() {
		defer utils.CloseTheCloser(b)
		defer utils.CloseTheCloser(pipeWriter)
		// read logs
		rd := bufio.NewReader(b)
		for {
			str, err := rd.ReadString('\n')
			if err != nil {
				if err == io.EOF && !isFollow {
					utils.LogDebug("End of file reached for logs", "docker", "GetLogs", map[string]interface{}{})
					return
				}
				_ = utils.LogError("Unable to read logs from container", "docker", "GetLogs", err)
				return
			}
			// starting 8 bytes of data contains some meta data regarding each log that docker sends
			// ignoring the first 8 bytes, send rest of the data
			fmt.Fprint(pipeWriter, str[8:])
		}
	}()
	return pipeReader, nil
}
