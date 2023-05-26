package sen

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

// GracefulShutdown creates a new GracefulShutdownPlugin.
// The plugin will allow the application calling its Shutdown
// when an interrupt signal (Ctrl+C) is received.
func GracefulShutdown() Plugin {
	return &gracefulShutdownPlugin{}
}

// gracefulShutdownPlugin is a plugin to allow the application to receive SIGTERM signal
// and shuts down the application gracefully.
type gracefulShutdownPlugin struct {
	Injector  Injector  `inject:"injector"`
	Lifecycle Lifecycle `inject:"lifecycle"`
}

func (s *gracefulShutdownPlugin) Initialize() error {
	shutdownCh := make(chan struct{})

	if err := s.Injector.Register("graceful-shutdown", s); err != nil {
		return err
	}

	s.Lifecycle.OnShutdown(func(ctx context.Context) error {
		close(shutdownCh)
		return nil
	})

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)

	s.Lifecycle.OnRun(func(ctx context.Context) error {
		go func() {
			select {
			case <-exit:
				_ = s.Lifecycle.Shutdown(ctx)
			case <-ctx.Done():
			case <-shutdownCh:
			}
		}()
		return nil
	})

	return nil
}
