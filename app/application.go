package app

import (
	"context"
)

// using function or interface?
// using interface + component allows other to interact with the hook
// we don't need as it can be a service
type Hook func(context.Context) error

// there are 2 ways it can be implemented:
// Using factory: nicer, isolated
// Injecting immediately, how to handle error?
// we need to inject immediately so application can inject further components while doing it
// how do we handle error? returning error always is annoying

type Plugin interface {
	Apply(ctx context.Context) error
}

type ComponentPlugin struct {
	App       *Application `inject:"app"`
	Name      string
	Component interface{}
}

func (p *ComponentPlugin) Apply(ctx context.Context) error {
	return p.App.Register(p.Name, p.Component)
}

func Component(name string, c interface{}) Plugin {
	return &ComponentPlugin{
		Name:      name,
		Component: c,
	}
}

type ModulePlugin struct {
	App     *Application `inject:"app"`
	Plugins []Plugin
}

// Module groups multiple plugins to act as a plugin.
func Module(plugins ...Plugin) Plugin {
	return &ModulePlugin{
		Plugins: plugins,
	}
}

func (m *ModulePlugin) Apply(ctx context.Context) error {
	for _, p := range m.Plugins {
		if err := m.App.applyPlugin(ctx, p); err != nil {
			return err
		}
	}

	return nil
}

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

	if err := app.applyPlugin(ctx, app.plugins); err != nil {
		return err
	}

	return app.defaultLifecycle.Run(ctx)
}

func (app *Application) applyPlugin(ctx context.Context, p Plugin) error {
	if err := app.Inject(p); err != nil {
		return err
	}

	return p.Apply(ctx)
}
