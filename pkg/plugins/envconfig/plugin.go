package envconfig

import (
	"github.com/caarlos0/env/v8"

	"github.com/bongnv/sen/pkg/sen"
)

// Config is a sen.Plugin to load a config from environment variables
// and registers it to Hub under the given name.
//
// # Usage
//
//	_ = app.With(envconfig.Config("echo.config", &echo.Config{}))
func Config(name string, cfg any) sen.Plugin {
	return &provider{
		name: name,
		cfg:  cfg,
	}
}

type provider struct {
	Hub sen.Hub `inject:"hub"`

	name string
	cfg  any
}

// Initialize loads the config from environment variables and registers it to the Hub.
func (p *provider) Initialize() error {
	if err := env.Parse(p.cfg); err != nil {
		return err
	}

	return p.Hub.Register(p.name, p.cfg)
}
