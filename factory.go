package sen

// IFactory is an interface to wrap New method.
// It represents a factory to create components.
type IFactory[T any] interface {
	New() (T, error)
}

// Factory creates a new plugin from a factory.
// The plugin will call New method in the given factory to create a new component.
// The component will then be registered into the application.
//
// The factory is useful when we need some dependencies to initialise a component.
func Factory[T any](name string, f IFactory[T]) Plugin {
	return &factoryPlugin[T]{
		Factory: f,
		Name:    name,
	}
}

type factoryPlugin[T any] struct {
	App     *Application `inject:"app"`
	Factory IFactory[T]
	Name    string
}

// Init initialises creates the component and registers it to the application.
func (p *factoryPlugin[_]) Init() error {
	if err := p.App.Inject(p.Factory); err != nil {
		return err
	}

	component, err := p.Factory.New()
	if err != nil {
		return err
	}

	return p.App.Register(p.Name, component)
}
