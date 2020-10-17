package server

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/spaceuptech/helpers"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/spaceuptech/space-cloud/runner/model"
)

func generateLogRequestFromQueryParams(ctx context.Context, r *url.URL) (*model.LogRequest, error) {
	taskID := r.Query().Get("taskId")
	replicaID := r.Query().Get("replicaId")
	since := r.Query().Get("since")
	sinceTime := r.Query().Get("since-time")
	tail := r.Query().Get("tail")
	_, isFollow := r.Query()["follow"]

	req := &model.LogRequest{
		TaskID:    taskID,
		ReplicaID: replicaID,
		IsFollow:  isFollow,
	}

	if replicaID == "" {
		return nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), "replica id not provided in query param", nil, map[string]interface{}{"since": since, "sinceTime": sinceTime})
	}

	if since != "" && sinceTime != "" {
		return nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), "you cannot only provide since or sinceTime at once", nil, map[string]interface{}{"since": since, "sinceTime": sinceTime})
	}

	if since != "" {
		sinceNum, err := time.ParseDuration(since)
		if err != nil {
			return nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("cannot parse since value (%s) to time", since), err, map[string]interface{}{"since": since})
		}
		temp := int64(sinceNum.Seconds())
		req.Since = &temp
	}

	if sinceTime != "" {
		sinceTimeDuration, err := time.Parse(time.RFC3339, sinceTime)
		if err != nil {
			return nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("cannot parse since time value (%s) to time", sinceTime), err, map[string]interface{}{"sinceTime": sinceTime})
		}
		temp := metav1.NewTime(sinceTimeDuration)
		req.SinceTime = &temp
	}

	if tail != "" {
		tailNum, err := strconv.Atoi(tail)
		if err != nil {
			return nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("incorrect value (%v) provided for tail ", tail), err, map[string]interface{}{"tail": tail})
		}
		if tailNum < 1 {
			// -1 indicates show all logs
			tailNum = -1
		}
		temp := int64(tailNum)
		req.Tail = &temp
	}
	return req, nil
}
