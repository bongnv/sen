package app

import "context"

// Plugin represents a plugin in a sen application. It enhances the application
// by proving one or multiple functionalities.
// A plugin can have from zero to many dependencies and they can be injected
// by declaring "inject" tag.
type Plugin interface {
	// Apply initialises the plugin and installs the plugin into the application.
	Apply(ctx context.Context) error
}

// ComponentPlugin is a simple plugin to add a component into the application.
// The component should have a name so it will be used as a dependency and
// will be injected to other components when needed.
type ComponentPlugin struct {
	App       *Application `inject:"app"`
	Name      string
	Component interface{}
}

// Apply adds the component to the application as a named dependency.
func (p *ComponentPlugin) Apply(ctx context.Context) error {
	return p.App.Register(p.Name, p.Component)
}

// Component creates a new ComponentPlugin.
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
		if err := m.App.ApplyPlugin(ctx, p); err != nil {
			return err
		}
	}

	return nil
}
