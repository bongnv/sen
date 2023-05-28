package sen

import (
	"context"
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
	hub Hub
	lc  Lifecycle
}

// New creates a new Application.
func New() *Application {
	app := &Application{
		hub: newHub(),
		lc:  newLifecycle(),
	}

	_ = app.hub.Register("app", app)
	_ = app.hub.Register("lifecycle", app.lc)

	return app
}

// With applies a plugin or multiple plugins.
// While applying a plugin, the plugin will be injected
// with dependencies and Init method will be called.
func (app *Application) With(plugins ...Plugin) error {
	for _, p := range plugins {
		if err := app.hub.Inject(p); err != nil {
			return err
		}

		if err := p.Initialize(); err != nil {
			return err
		}
	}

	return nil
}

// Run runs the application.
// It will execute all run hooks, then shutdown hooks and afterRun hooks.
func (app *Application) Run(ctx context.Context) error {
	return app.lc.Run(ctx)
}

// Shutdown runs the application by executing all the registered hooks for this phase.
func (app *Application) Shutdown(ctx context.Context) error {
	return app.lc.Shutdown(ctx)
}
