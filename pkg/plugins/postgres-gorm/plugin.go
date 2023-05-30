package postgresgorm

import (
	"github.com/bongnv/sen/pkg/plugins/envconfig"
	"github.com/bongnv/sen/pkg/sen"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Bundle is a sen.Plugin that provides both Config and gorm.DB for convenience.
//
// # Usage
//
//	app.With(postgresgorm.Bundle())
func Bundle() sen.Plugin {
	return sen.Bundle(
		envconfig.Config("postgresgorm.config", &Config{}),
		&Plugin{},
	)
}

// Config includes configuration to initialize a gorm.DB to a PostgreSQL server.
type Config struct {
	DSN string `env:"POSTGRESQL_DSN,required"`
}

// Plugin is a plugin that provides an instance of gorm.DB. Check gorm.io
// to see how to use gorm to work with databases.
// The plugin requires Config is registered in advance.
//
// # Usage
//
//	app.With(&postgresgorm.Plugin{})
type Plugin struct {
	Hub    sen.Hub `inject:"hub"`
	Config *Config `inject:"postgresgorm.config"`
}

// Initialize creates a new instance of *gorm.DB with the given configuration.
// And then registers it into the application as `gorm`.
func (p *Plugin) Initialize() error {
	db, err := gorm.Open(postgres.Open(p.Config.DSN), &gorm.Config{})
	if err != nil {
		return err
	}

	return p.Hub.Register("gorm", db)
}
