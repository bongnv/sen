package app

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
)

// Hook represents a hook which allows to customise the application life cycle.
type Hook func(context.Context) error

// Lifecycle represents an application life cycle.
type Lifecycle interface {
	// OnRun adds additional logic when the app runs. For a long lasting service
	// it should only exit the function when the service no longer runs.
	OnRun(Hook)
	// OnShutdown adds additional logic when the app shuts down.
	OnShutdown(Hook)
}

type defaultLifecycle struct {
	shutdownOnce      int32
	runHooks          []Hook
	shutdownHooks     []Hook
	postShutdownHooks []Hook
}

func (l *defaultLifecycle) OnRun(h Hook) {
	l.runHooks = append(l.runHooks, h)
}

func (l *defaultLifecycle) OnShutdown(h Hook) {
	l.shutdownHooks = append(l.shutdownHooks, h)
}

// Run runs the application by executing all the registered hooks for this phase.
func (l *defaultLifecycle) Run(ctx context.Context) error {
	return executeHooks(ctx, l.runHooks)
}

// Shutdown runs the application by executing all the registered hooks for this phase.
// It will include OnShutdown & AfterShutdown.
func (l *defaultLifecycle) Shutdown(ctx context.Context) error {
	if swapped := atomic.CompareAndSwapInt32(&l.shutdownOnce, 0, 1); !swapped {
		return errors.New("app: Shutdown has been called")
	}

	if err := executeHooks(ctx, l.shutdownHooks); err != nil {
		return err
	}

	return executeHooks(ctx, l.postShutdownHooks)
}

func executeHooks(ctx context.Context, hooks []Hook) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	errCh := make(chan error)
	wg := sync.WaitGroup{}

	for _, h := range hooks {
		wg.Add(1)
		h := h
		go func() {
			if err := h(ctx); err != nil {
				select {
				case errCh <- err:
				case <-ctx.Done():
				}
			}
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(errCh)
	}()

	return <-errCh
}
