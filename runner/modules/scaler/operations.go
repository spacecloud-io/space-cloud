package scaler

import (
	"context"
	"fmt"
	"net"
	"time"

	"google.golang.org/grpc"

	pb "github.com/spaceuptech/space-cloud/runner/modules/scaler/externalscaler"
)

// ScaleUp instructs keda to scale the service from 0 to 1
func (s *Scaler) ScaleUp(ctx context.Context, project, service, version string) error {
	key := generateKey(project, service, version)
	// Check if service is active
	exists, err := s.pubsubClient.CheckAndSet(ctx, key, "service", 10*time.Second)
	if err != nil {
		return err
	}

	// Notify everyone to scale up the service
	if !exists {
		if err := s.pubsubClient.PublishString(ctx, "scale-up", key); err != nil {
			return err
		}
	}

	return nil
}

// Start begins the grpc server
func (s *Scaler) Start() {
	// Start the internal routines
	go s.routineScaleUp()

	// Create a gRPC server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 4060))
	if err != nil {
		panic("Unable to start grpc server: " + err.Error())
	}

	grpcServer := grpc.NewServer()
	pb.RegisterExternalScalerServer(grpcServer, s)
	if err := grpcServer.Serve(lis); err != nil {
		panic("Unable to start grpc server: " + err.Error())
	}
}
