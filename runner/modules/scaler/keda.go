package scaler

import (
	"context"
	"fmt"
	"strconv"

	"github.com/spaceuptech/helpers"

	pb "github.com/spaceuptech/space-cloud/runner/modules/scaler/externalscaler"
)

func (s *Scaler) GetMetricSpec(ctx context.Context, scaledObject *pb.ScaledObjectRef) (*pb.GetMetricSpecResponse, error) {
	// Extract the target details first
	project, p1 := scaledObject.ScalerMetadata["project"]
	service, p2 := scaledObject.ScalerMetadata["service"]
	version, p3 := scaledObject.ScalerMetadata["version"]
	scalingMode, p4 := scaledObject.ScalerMetadata["type"]
	targetString, p5 := scaledObject.ScalerMetadata["target"]

	// Throw an error if the required fields are not present
	if !p1 || !p2 || !p3 || !p4 || !p5 {
		return nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), "Invalid fields provided in keda scaler config", nil, map[string]interface{}{"project": project, "service": service, "version": version, "target": targetString})
	}

	// Convert target to integer
	target, err := strconv.Atoi(targetString)
	if err != nil {
		return nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), "Invalid target provided in keda scaler config", nil, map[string]interface{}{"project": project, "service": service, "version": version, "target": targetString})
	}

	// Send response
	helpers.Logger.LogDebug(helpers.GetRequestID(ctx), fmt.Sprintf("Keda get metrics spec result: (%s, %s)", scalingMode, targetString), map[string]interface{}{"project": project, "service": service, "version": version})
	return &pb.GetMetricSpecResponse{
		MetricSpecs: []*pb.MetricSpec{{
			MetricName: scalingMode,
			TargetSize: int64(target),
		}},
	}, nil
}

func (s *Scaler) GetMetrics(ctx context.Context, metricRequest *pb.GetMetricsRequest) (*pb.GetMetricsResponse, error) {
	// Extract the service details first
	project, p1 := metricRequest.ScaledObjectRef.ScalerMetadata["project"]
	service, p2 := metricRequest.ScaledObjectRef.ScalerMetadata["service"]
	version, p3 := metricRequest.ScaledObjectRef.ScalerMetadata["version"]

	// Throw an error if the fields are not present
	if !p1 || !p2 || !p3 {
		return nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), "Invalid fields provided in keda scaler config", nil, map[string]interface{}{"project": project, "service": service, "version": version})
	}

	// Get the scaling mode (metric type)
	scalingMode := metricRequest.MetricName

	// Query prometheus to get the metrics
	metric, err := s.queryPrometheus(ctx, project, service, version, scalingMode)
	if err != nil {
		return nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), "Unable to fetch scaling metrics for service", err, map[string]interface{}{"project": project, "service": service, "version": version})
	}

	helpers.Logger.LogDebug(helpers.GetRequestID(ctx), fmt.Sprintf("Keda get metrics result: (%v, %v)", scalingMode, metric), map[string]interface{}{"project": project, "service": service, "version": version})
	return &pb.GetMetricsResponse{
		MetricValues: []*pb.MetricValue{{
			MetricName:  scalingMode,
			MetricValue: metric,
		}},
	}, nil
}

// StreamIsActive is called by keda controller once to allow us to stream isActive events to it.
func (s *Scaler) StreamIsActive(scaledObject *pb.ScaledObjectRef, epsServer pb.ExternalScaler_StreamIsActiveServer) error {
	ctx := epsServer.Context()

	// Extract the service details first
	project, p1 := scaledObject.ScalerMetadata["project"]
	service, p2 := scaledObject.ScalerMetadata["service"]
	version, p3 := scaledObject.ScalerMetadata["version"]

	// Throw an error if the fields are not present
	if !p1 || !p2 || !p3 {
		return helpers.Logger.LogError(helpers.GetRequestID(ctx), "Invalid fields provided in keda scaler config", nil, map[string]interface{}{"project": project, "service": service, "version": version})
	}

	helpers.Logger.LogDebug(helpers.GetRequestID(ctx), "Keda stream is active called", map[string]interface{}{"project": project, "service": service, "version": version})

	// Create a channel which we can use to write data to this stream
	ch := make(chan bool, 5)

	// Add stream to internal map
	s.addIsActiveStream(project, service, version, ch)
	defer s.removeIsActiveStream(project, service, version)

	for {
		select {
		// Exit if the stream has closed
		case <-ctx.Done():
			helpers.Logger.LogDebug(helpers.GetRequestID(ctx), "Exiting from keda is active stream", map[string]interface{}{"project": project, "service": service, "version": version})
			return nil

		// Forward isActive messages to server
		case isActive := <-ch:
			// Quit if isActive is false
			if !isActive {
				return nil
			}

			helpers.Logger.LogDebug(helpers.GetRequestID(ctx), fmt.Sprintf("Streaming is active result: %v", isActive), map[string]interface{}{"project": project, "service": service, "version": version})
			if err := epsServer.Send(&pb.IsActiveResponse{Result: isActive}); err != nil {
				_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), "Unable to forward is stream active resonse to keda", err, map[string]interface{}{"project": project, "service": service, "version": version})
			}
		}
	}
}

// IsActive is called by the keda controller to check if keda should scale the deployment from 0 to 1
func (s *Scaler) IsActive(ctx context.Context, scaledObject *pb.ScaledObjectRef) (*pb.IsActiveResponse, error) {
	// Extract the service details first
	project, p1 := scaledObject.ScalerMetadata["project"]
	service, p2 := scaledObject.ScalerMetadata["service"]
	version, p3 := scaledObject.ScalerMetadata["version"]
	minReplicas, p4 := scaledObject.ScalerMetadata["minReplicas"]

	// Throw an error if the fields are not present
	if !p1 || !p2 || !p3 || !p4 {
		return nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), "Invalid fields provided in keda scaler config", nil, map[string]interface{}{"project": project, "service": service, "version": version})
	}

	// Check if service is active
	isActive := s.isStreamActive(project, service, version, minReplicas)

	helpers.Logger.LogDebug(helpers.GetRequestID(ctx), fmt.Sprintf("Keda is active result: %v", isActive), map[string]interface{}{"project": project, "service": service, "version": version})
	return &pb.IsActiveResponse{Result: isActive}, nil
}
