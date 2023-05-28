package zap

import (
	"context"

	"go.uber.org/zap"

	"github.com/bongnv/sen/pkg/sen"
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

	// redirect std log to logger
	revertStdLog := zap.RedirectStdLog(logger)

	p.LC.PostRun(func(_ context.Context) error {
		revertStdLog()
		// Sync may return errors weirdly. It's best to ignore it for now. See https://github.com/uber-go/zap/issues/991
		_ = logger.Sync()
		return nil
	})

	return p.Injector.Register("logger", logger)
}
