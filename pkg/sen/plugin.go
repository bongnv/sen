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
// The simple plugin adds a component into the application
// under the given name.
func Component(name string, c any) Plugin {
	return &componentPlugin{
		name:      name,
		component: c,
	}
}

// Bundle is a collection of plugins.
func Bundle(plugins ...Plugin) Plugin {
	return &bundlePlugin{
		plugins: plugins,
	}
}

type bundlePlugin struct {
	App *Application `inject:"app"`

	plugins []Plugin
}

func (m bundlePlugin) Initialize() error {
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

// Initialize adds the hook to the application lifecycle.
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

// Initialize adds the hook to the application lifecycle.
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

// Initialize adds the hook to the application lifecycle.
func (p postRunPlugin) Initialize() error {
	for _, h := range p.hooks {
		p.LC.PostRun(h)
	}
	return nil
}
