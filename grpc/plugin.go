package grpc

import (
	"context"
	"fmt"
	"net"
	"os"

	"google.golang.org/grpc"

	"github.com/bongnv/sen"
)

// Config contains configurations for starting a GRPC server.
type Config struct {
	Address       string
	ServerOptions []grpc.ServerOption
}

// EnvConfigPlugin creates a plugin to load GRPC address from
// the environment variable GRPC_ADDR.
func EnvConfigPlugin(serverOptions ...grpc.ServerOption) sen.Plugin {
	return &envConfigPlugin{
		serverOptions: serverOptions,
	}
}

type envConfigPlugin struct {
	App           *sen.Application `inject:"app"`
	serverOptions []grpc.ServerOption
}

// Init initialises Config.
func (p *envConfigPlugin) Init() error {
	address := os.Getenv("GRPC_ADDR")
	cfg := &Config{
		Address:       address,
		ServerOptions: p.serverOptions,
	}

	return p.App.Register("grpc.config", cfg)
}

// Plugin creates a new sen plugin for running a GRPC server.
func Plugin() sen.Plugin {
	return &basePlugin{}
}

type basePlugin struct {
	App    *sen.Application `inject:"app"`
	Config *Config          `inject:"grpc.config"`
}

func (p *basePlugin) Init() error {
	lis, err := net.Listen("tcp", p.Config.Address)
	if err != nil {
		return fmt.Errorf("grpc: failed to listen to %s: %v", p.Config.Address, err)
	}

	grpcServer := grpc.NewServer(p.Config.ServerOptions...)

	p.App.OnRun(func(ctx context.Context) error {
		return grpcServer.Serve(lis)
	})

	p.App.OnShutdown(func(ctx context.Context) error {
		doneCh := make(chan struct{})
		go func() {
			grpcServer.GracefulStop()
			close(doneCh)
		}()

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-doneCh:
			return nil
		}
	})

	return p.App.Register("grpc.server", grpcServer)
}
