package app

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
// To construct an application from plugins use, Apply. For example:
//
//	app := sen.New()
//	if err := app.Apply(plugin1, plugin2); err != nil {
//	   handleError(err)
//	}
type Application struct {
	injector *defaultInjector

	runHooks      []Hook
	shutdownHooks []Hook
	afterRunHooks []Hook
	shutdownErr   error
	shutdownOnce  sync.Once
}

// New creates a new Application.
func New() *Application {
	app := &Application{
		injector: newInjector(),
	}

	_ = app.Register("app", app)
	return app
}

// Register registers a new component into the application.
func (app *Application) Register(name string, component interface{}) error {
	return app.injector.Register(name, component)
}

// Inject injects dependencies into the given component.
func (app *Application) Inject(component interface{}) error {
	return app.injector.Inject(component)
}

// Retrieve retrives a component via its registered name.
func (app *Application) Retrieve(name string) (interface{}, error) {
	return app.injector.Retrieve(name)
}

// OnRun adds additional logic when the app runs. For a long lasting service
// it should only block the function until the service no longer runs.
func (app *Application) OnRun(h Hook) {
	app.runHooks = append(app.runHooks, h)
}

// AfterRun adds additional logic after all services stop running.
// It's useful for syncing logs, etc.
func (app *Application) AfterRun(h Hook) {
	app.shutdownHooks = append(app.afterRunHooks, h)
}

// OnShutdown adds additional logic when the app shuts down.
func (app *Application) OnShutdown(h Hook) {
	app.shutdownHooks = append(app.shutdownHooks, h)
}

// Run runs the application by executing all the registered hooks for this phase.
func (app *Application) Run(ctx context.Context) error {
	err := executeHooks(ctx, app.runHooks)
	if err != nil {
		return err
	}

	return executeHooks(ctx, app.afterRunHooks)
}

// Shutdown runs the application by executing all the registered hooks for this phase.
func (app *Application) Shutdown(ctx context.Context) error {
	app.shutdownOnce.Do(func() {
		app.shutdownErr = executeHooks(ctx, app.shutdownHooks)
	})

	return app.shutdownErr
}

// With applies a plugin or multiple plugins.
// While applying a plugin, the plugin will be injected
// with dependencies and Init method will be called.
func (app *Application) With(plugins ...Plugin) error {
	for _, p := range plugins {
		if err := app.Inject(p); err != nil {
			return err
		}

		if err := p.Initialize(); err != nil {
			return err
		}
	}

	return nil
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
