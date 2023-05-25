package plugin

import "github.com/bongnv/sen/app"

// Factory is an interface to wrap New method.
// It represents a factory to create components.
type Factory[T any] interface {
	New() (T, error)
}

// Provider creates a new plugin from a factory to provide a new component.
// The plugin will call New method in the given factory to create a new component.
// The component will then be registered into the application.
//
// The factory is useful when we need some dependencies to initialise a component.
func Provider[T any](name string, f Factory[T]) app.Plugin {
	return &providerPlugin[T]{
		Factory: f,
		Name:    name,
	}
}

type providerPlugin[T any] struct {
	App     *app.Application `inject:"app"`
	Factory Factory[T]
	Name    string
}

// Init initialises creates the component and registers it to the application.
func (p *providerPlugin[_]) Initialize() error {
	if err := p.App.Inject(p.Factory); err != nil {
		return err
	}

	component, err := p.Factory.New()
	if err != nil {
		return err
	}

	return p.App.Register(p.Name, component)
}
