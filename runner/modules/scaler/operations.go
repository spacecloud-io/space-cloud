package scaler

import (
	"fmt"
	"net"

	"google.golang.org/grpc"

	pb "github.com/spaceuptech/space-cloud/runner/modules/scaler/externalscaler"
)

// ScaleUp instructs keda to scale the service from 0 to 1
func (s *Scaler) ScaleUp(project, service, version string) error {
	s.lock.RLock()
	defer s.lock.RUnlock()

	// Check if the stream exists
	stream, p := s.isActiveStreams[generateKey(project, service, version)]
	if !p {
		return fmt.Errorf("provided service cannot be scaled")
	}

	// Send signal to scale up
	stream.Notify()

	return nil
}

// Start begins the grpc server
func (s *Scaler) Start() {
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
