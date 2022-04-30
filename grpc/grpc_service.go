package grpc

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"
)

type Config struct {
	Address       string
	ServerOptions []grpc.ServerOption
}

// GRPCService will be used like?
// app.WithService(&grpc.New(opts))
type GRPCService struct {
	*grpc.Server

	Config Config `inject:"grpc-config"`

	lis net.Listener
}

func (s *GRPCService) Init(ctx context.Context) error {
	lis, err := net.Listen("tcp", s.Config.Address)
	if err != nil {
		return fmt.Errorf("failed to listen to %s: %v", s.Config.Address, err)
	}

	s.lis = lis
	s.Server = grpc.NewServer(s.Config.ServerOptions...)
	return nil
}

// Run is called to start the service. The function shouldn't be returned if the service is still running.
func (s GRPCService) Run(ctx context.Context) error {
	return s.Server.Serve(s.lis)
}

// Shutdown is called to graceful shut down the service.
func (s GRPCService) Shutdown(ctx context.Context) error {
	doneCh := make(chan struct{})
	go func() {
		s.Server.GracefulStop()
		close(doneCh)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-doneCh:
		return nil
	}
}
