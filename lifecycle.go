package sen

import (
	"context"
	"sync"

	"golang.org/x/sync/errgroup"
)

// Lifecycle manages the lifecycle of an application.
// An application starts with .Run(ctx) and will be stopped
// when .Shutdown(ctx) is called.
// It also allows to hook into the application via OnRun, OnShutdown and PostRun.
type Lifecycle interface {
	OnRun(h Hook)
	OnShutdown(h Hook)
	PostRun(h Hook)
	Run(ctx context.Context) error
	Shutdown(ctx context.Context) error
}

type defaultLifecycle struct {
	runHooks      []Hook
	shutdownHooks []Hook
	postRunHooks  []Hook
	shutdownOnce  func(ctx context.Context) error
}

// OnRun adds additional logic when the app runs. For a long lasting service
// it should only block the function until the service no longer runs.
func (lc *defaultLifecycle) OnRun(h Hook) {
	lc.runHooks = append(lc.runHooks, h)
}

// PostRun adds additional logic after all services stop running
// and shutdown logic is executed.
// It's useful for syncing logs, etc.
func (lc *defaultLifecycle) PostRun(h Hook) {
	lc.postRunHooks = append(lc.postRunHooks, h)
}

// OnShutdown adds additional logic when the app shuts down.
func (lc *defaultLifecycle) OnShutdown(h Hook) {
	lc.shutdownHooks = append(lc.shutdownHooks, h)
}

// Run runs the application by executing all the registered hooks for this phase.
func (lc *defaultLifecycle) Run(ctx context.Context) (err error) {
	err = executeHooks(ctx, lc.runHooks)
	shutdownErr := lc.shutdownOnce(ctx)
	if shutdownErr != nil && err == nil {
		err = shutdownErr
	}

	afterRunErr := executeHooks(ctx, lc.postRunHooks)
	if afterRunErr != nil && err == nil {
		err = afterRunErr
	}

	return
}

// Shutdown runs the application by executing all the registered hooks for this phase.
func (lc *defaultLifecycle) Shutdown(ctx context.Context) error {
	return lc.shutdownOnce(ctx)
}

// internalShutdown is the internal implementation of the shutdown function.
// It shouldn't be called multiple times so it should be wrapped to run once only.
func (lc *defaultLifecycle) internalShutdown(ctx context.Context) error {
	return executeHooks(ctx, lc.shutdownHooks)
}

func newLifecycle() Lifecycle {
	lc := &defaultLifecycle{}
	lc.shutdownOnce = runOnce(lc.internalShutdown)
	return lc
}

func executeHooks(ctx context.Context, hooks []Hook) error {
	eg, ctx := errgroup.WithContext(ctx)
	for _, h := range hooks {
		h := h
		eg.Go(func() error {
			return h(ctx)
		})
	}

	return eg.Wait()
}

// runOnce allows creates a function that will call fn only once.
// It's different from sync.Once that, all calls will be blocked and returns
// the error from the single call of fn.
func runOnce(fn func(ctx context.Context) error) func(ctx context.Context) error {
	var err error
	once := &sync.Once{}
	done := make(chan struct{})
	return func(ctx context.Context) error {
		once.Do(func() {
			err = fn(ctx)
			close(done)
		})
		select {
		case <-done:
			return err
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
