package app

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
)

// using function or interface?
// using interface + component allows other to interact with the hook
// we don't need as it can be a service
type Hook func(context.Context) error

type Lifecycle interface {
	OnRun(Hook)
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

// Run runs the application and its services.
func (l *defaultLifecycle) Run(ctx context.Context) error {
	return executeHooks(ctx, l.runHooks)
}

func (l *defaultLifecycle) Shutdown(ctx context.Context) error {
	if swapped := atomic.CompareAndSwapInt32(&l.shutdownOnce, 0, 1); !swapped {
		return errors.New("app: the app is already shutdown")
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	if err := executeHooks(ctx, l.shutdownHooks); err != nil {
		return err
	}

	return executeHooks(ctx, l.postShutdownHooks)
}

func executeHooks(ctx context.Context, hooks []Hook) error {
	errCh := make(chan error)
	wg := sync.WaitGroup{}

	for _, h := range hooks {
		wg.Add(1)
		h := h
		go func() {
			defer wg.Done()
			errCh <- h(ctx)
		}()
	}

	go func() {
		wg.Wait()
		close(errCh)
	}()

	return <-errCh
}
