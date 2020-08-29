package istio

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/spaceuptech/helpers"
	"k8s.io/api/core/v1"

	"github.com/spaceuptech/space-cloud/runner/utils"
)

// GetLogs get logs of specified services
func (i *Istio) GetLogs(ctx context.Context, isFollow bool, projectID, taskID, replica string) (io.ReadCloser, error) {
	if taskID == "" {
		arr := strings.Split(replica, "-")
		if len(arr) < 2 {
			return nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), "Invalid replica id", nil, nil)
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
	helpers.Logger.LogDebug(helpers.GetRequestID(ctx), "Sending logs to client", map[string]interface{}{})
	go func() {
		defer utils.CloseTheCloser(b)
		defer utils.CloseTheCloser(pipeWriter)
		// read logs
		rd := bufio.NewReader(b)
		for {
			str, err := rd.ReadString('\n')
			if err != nil {
				if err == io.EOF && !isFollow {
					helpers.Logger.LogDebug(helpers.GetRequestID(ctx), "End of file reached for logs", map[string]interface{}{})
					return
				}
				_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), "Unable to read logs from container", err, nil)
				return
			}
			fmt.Fprint(pipeWriter, str)
		}
	}()
	return pipeReader, nil
}
