package app

import (
	"context"
)

// there are 2 ways it can be implemented:
// Using factory: nicer, isolated
// Injecting immediately, how to handle error?
// we need to inject immediately so application can inject further components while doing it
// how do we handle error? returning error always is annoying

type Application struct {
	Injector
	defaultLifecycle
	plugins Plugin
}

func New(plugins ...Plugin) *Application {
	app := &Application{
		plugins:  Module(plugins...),
		Injector: newInjector(),
	}

	return app
}

func (app *Application) Run(ctx context.Context) error {
	if err := app.Register("app", app); err != nil {
		return err
	}

	if err := app.ApplyPlugin(ctx, app.plugins); err != nil {
		return err
	}

	return app.defaultLifecycle.Run(ctx)
}

func (app *Application) ApplyPlugin(ctx context.Context, p Plugin) error {
	if err := app.Inject(p); err != nil {
		return err
	}

	return p.Apply(ctx)
}
