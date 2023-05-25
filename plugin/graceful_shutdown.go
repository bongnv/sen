package plugin

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/bongnv/sen/app"
)

// GracefulShutdown creates a new GracefulShutdownPlugin.
func GracefulShutdown() app.Plugin {
	return &gracefulShutdownPlugin{}
}

// gracefulShutdownPlugin is a plugin to allow the application to receive SIGTERM signal
// and shuts down the application gracefully.
type gracefulShutdownPlugin struct {
	App *app.Application `inject:"*"`
}

func (s *gracefulShutdownPlugin) Initialize() error {
	shutdownCh := make(chan struct{})

	if err := s.App.Register("graceful-shutdown", s); err != nil {
		return err
	}

	s.App.OnShutdown(func(ctx context.Context) error {
		close(shutdownCh)
		return nil
	})

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)

	s.App.OnRun(func(ctx context.Context) error {
		go func() {
			select {
			case <-exit:
				_ = s.App.Shutdown(ctx)
			case <-ctx.Done():
			case <-shutdownCh:
			}
		}()
		return nil
	})

	return nil
}
