package sen

import (
	"context"
	"sync"

	"golang.org/x/sync/errgroup"
)

// Hook represents a hook to add custom logic in the application life cycle.
type Hook func(ctx context.Context) error

// Plugin represents a plugin in a sen application. It enhances the application
// by providing one or multiple functionalities.
// A plugin can have from zero to many dependencies and they can be injected
// by declaring "inject" tag.
type Plugin interface {
	// Initialize initialises the plugin and installs the plugin into the application.
	Initialize() error
}

// Application represents an application.
// To construct an application from plugins use: app.With. For example:
//
//	app, err := sen.New(plugin1, plugin2)
//	if err != nil {
//		handleError(err)
//	}
type Application struct {
	runHooks      []Hook
	shutdownHooks []Hook
	postRunHooks  []Hook
	shutdownOnce  func(ctx context.Context) error
}

// New creates a new Application from plugins.
func New(plugins ...Plugin) (*Application, error) {
	app := &Application{}
	app.shutdownOnce = runOnce(app.internalShutdown)

	hub := newHub()

	_ = hub.Register("app", app)
	for _, p := range plugins {
		if err := hub.Inject(p); err != nil {
			return nil, err
		}

		if err := p.Initialize(); err != nil {
			return nil, err
		}
	}

	return app, nil
}

// OnRun adds additional logic when the app runs. For a long lasting service
// it should only block the function until the service no longer runs.
func (app *Application) OnRun(h Hook) {
	app.runHooks = append(app.runHooks, h)
}

// PostRun adds additional logic after all services stop running
// and shutdown logic is executed.
// It's useful for syncing logs, etc.
func (app *Application) PostRun(h Hook) {
	app.postRunHooks = append(app.postRunHooks, h)
}

// OnShutdown adds additional logic when the app shuts down.
func (app *Application) OnShutdown(h Hook) {
	app.shutdownHooks = append(app.shutdownHooks, h)
}

// Run runs the application by executing all the registered hooks for this phase.
func (app *Application) Run(ctx context.Context) (err error) {
	err = executeHooks(ctx, app.runHooks)
	shutdownErr := app.shutdownOnce(ctx)
	if shutdownErr != nil && err == nil {
		err = shutdownErr
	}

	afterRunErr := executeHooks(ctx, app.postRunHooks)
	if afterRunErr != nil && err == nil {
		err = afterRunErr
	}

	return
}

// Shutdown runs the application by executing all the registered hooks for this phase.
func (app *Application) Shutdown(ctx context.Context) error {
	return app.shutdownOnce(ctx)
}

// internalShutdown is the internal implementation of the shutdown function.
// It shouldn't be called multiple times so it should be wrapped to run once only.
func (app *Application) internalShutdown(ctx context.Context) error {
	return executeHooks(ctx, app.shutdownHooks)
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
