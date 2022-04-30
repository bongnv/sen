package shutdown

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/bongnv/sen/app"
)

// GracefulShutdown shuts down the app gracefully.
type GracefulShutdown struct {
	App        *app.Application `inject:"app"`
	shutdownCh chan struct{}
}

func New() *GracefulShutdown {
	return &GracefulShutdown{}
}

func (s *GracefulShutdown) Apply(ctx context.Context) error {
	s.shutdownCh = make(chan struct{})
	s.App.Register("graceful-shutdown", s)
	s.App.OnShutdown(s.Shutdown)
	s.App.OnRun(s.Run)
	return nil
}

func (s GracefulShutdown) Run(ctx context.Context) error {
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)

	select {
	case <-exit:
		return s.App.Shutdown(ctx)
	case <-ctx.Done():
		// context is cancel only when app.Run returns with error.
		// should shutdown the app gracefully
		return s.App.Shutdown(ctx)
	case <-s.shutdownCh:
		return nil
	}
}

func (s *GracefulShutdown) Shutdown(ctx context.Context) error {
	close(s.shutdownCh)
	return nil
}
