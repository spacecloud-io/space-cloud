package istio

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/spaceuptech/helpers"
	v1 "k8s.io/api/core/v1"

	"github.com/spaceuptech/space-cloud/runner/model"
	"github.com/spaceuptech/space-cloud/runner/utils"
)

// GetLogs get logs of specified services
func (i *Istio) GetLogs(ctx context.Context, projectID string, info *model.LogRequest) (io.ReadCloser, error) {
	if info.TaskID == "" {
		arr := strings.Split(info.ReplicaID, "-")
		if len(arr) < 2 {
			return nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), "Invalid replica id", nil, nil)
		}
		info.TaskID = arr[0]
	}
	// get logs of pods
	req := i.kube.CoreV1().Pods(projectID).GetLogs(info.ReplicaID, &v1.PodLogOptions{
		Container:    info.TaskID,
		Follow:       info.IsFollow,
		Timestamps:   true,
		SinceSeconds: info.Since,
		SinceTime:    info.SinceTime,
		TailLines:    info.Tail,
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
				if err == io.EOF && !info.IsFollow {
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
