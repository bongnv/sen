package sen

// Plugin represents a plugin in a sen application. It enhances the application
// by proving one or multiple functionalities.
// A plugin can have from zero to many dependencies and they can be injected
// by declaring "inject" tag.
type Plugin interface {
	// Init initialises the plugin and installs the plugin into the application.
	Init() error
}

type componentPlugin struct {
	App       *Application `inject:"app"`
	Name      string
	Component interface{}
}

// Init adds the component to the application as a named dependency.
func (p *componentPlugin) Init() error {
	return p.App.Register(p.Name, p.Component)
}

// Component creates a new component plugin.
// It is a simple plugin to add a component into the application.
// The component should have a name so it will be used as a dependency and
// will be injected to other components when needed.
func Component(name string, c interface{}) Plugin {
	return &componentPlugin{
		Name:      name,
		Component: c,
	}
}

// ModulePlugin composes multiple plugins to act as a plugin.
type ModulePlugin struct {
	App     *Application `inject:"app"`
	Plugins []Plugin
}

// Module creates a ModulePlugin by providing a group of plugins.
func Module(plugins ...Plugin) Plugin {
	return &ModulePlugin{
		Plugins: plugins,
	}
}

// Init installs plugins into the application.
func (m *ModulePlugin) Init() error {
	for _, p := range m.Plugins {
		if err := m.App.Inject(p); err != nil {
			return err
		}

		if err := p.Init(); err != nil {
			return err
		}
	}

	return nil
}
