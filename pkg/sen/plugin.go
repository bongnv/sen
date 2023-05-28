package sen

type componentPlugin struct {
	Hub Hub `inject:"hub"`

	name      string
	component any
}

// Initialize adds the component to the application as a named dependency.
func (p *componentPlugin) Initialize() error {
	return p.Hub.Register(p.name, p.component)
}

// Component creates a new component plugin.
// It is a simple plugin to add a component into the application.
// The component will be registered with the given name.
// Init method will be called to initialize the component
// after dependencies are injected.
func Component(name string, c any) Plugin {
	return &componentPlugin{
		name:      name,
		component: c,
	}
}

// Module is a collection of plugins.
func Module(plugins ...Plugin) Plugin {
	return &modulePlugin{
		plugins: plugins,
	}
}

type modulePlugin struct {
	App *Application `inject:"app"`

	plugins []Plugin
}

func (m modulePlugin) Initialize() error {
	return m.App.With(m.plugins...)
}

// OnRun adds multiple hooks to run with the application.
func OnRun(hooks ...Hook) Plugin {
	return &onRunPlugin{
		hooks: hooks,
	}
}

type onRunPlugin struct {
	LC    Lifecycle `inject:"lifecycle"`
	hooks []Hook
}

// Initialize adds the component to the application as a named dependency.
func (p onRunPlugin) Initialize() error {
	for _, h := range p.hooks {
		p.LC.OnRun(h)
	}
	return nil
}

// OnShutdown adds multiple hooks to run with the application.
func OnShutdown(hooks ...Hook) Plugin {
	return &onShutdownPlugin{
		hooks: hooks,
	}
}

type onShutdownPlugin struct {
	LC    Lifecycle `inject:"lifecycle"`
	hooks []Hook
}

// Initialize adds the component to the application as a named dependency.
func (p onShutdownPlugin) Initialize() error {
	for _, h := range p.hooks {
		p.LC.OnShutdown(h)
	}
	return nil
}

// PostRun adds additional logic after all services stop running
// and shutdown logic is executed.
// It's useful for syncing logs, etc.
func PostRun(hooks ...Hook) Plugin {
	return &postRunPlugin{
		hooks: hooks,
	}
}

type postRunPlugin struct {
	LC    Lifecycle `inject:"lifecycle"`
	hooks []Hook
}

// Initialize adds the component to the application as a named dependency.
func (p postRunPlugin) Initialize() error {
	for _, h := range p.hooks {
		p.LC.PostRun(h)
	}
	return nil
}
