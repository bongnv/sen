package zap

import (
	"context"

	"go.uber.org/zap"

	"github.com/bongnv/sen"
)

// Plugin is a sen.Plugin that provides an instance of *zap.Logger.
//
// # Usage
//
//	app.With(&zap.Plugin{
//		Options: zapOptions,
//	})
type Plugin struct {
	Options []zap.Option

	LC       sen.Lifecycle `inject:"lifecycle"`
	Injector sen.Injector  `inject:"injector"`
}

// Initialize initialises zap logger for the application.
// The logger will be regisreted under "logger" tag.
func (p Plugin) Initialize() error {
	logger, err := zap.NewProduction(p.Options...)
	if err != nil {
		return err
	}

	p.LC.AfterRun(func(_ context.Context) error {
		return logger.Sync()
	})

	return p.Injector.Register("logger", logger)
}
