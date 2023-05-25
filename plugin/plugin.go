package plugin

import "github.com/bongnv/sen/app"

type componentPlugin struct {
	App *app.Application `inject:"*"`

	name      string
	component any
}

// Init adds the component to the application as a named dependency.
func (p componentPlugin) Initialize() error {
	return p.App.Register(p.name, p.component)
}

// Component creates a new component plugin.
// It is a simple plugin to add a component into the application.
// The component will be registered with the given name.
// Init method will be called to initialize the component
// after dependencies are injected.
func Component(name string, c any) app.Plugin {
	return &componentPlugin{
		name:      name,
		component: c,
	}
}

// Module is a collection of plugins.
func Module(plugins ...app.Plugin) app.Plugin {
	return &modulePlugin{
		plugins: plugins,
	}
}

type modulePlugin struct {
	App *app.Application `inject:"*"`

	plugins []app.Plugin
}

func (m modulePlugin) Initialize() error {
	return m.App.With(m.plugins...)
}
