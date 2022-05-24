package zap

import (
	"context"

	"go.uber.org/zap"

	"github.com/bongnv/sen"
)

// Plugin creates a new sen plugin that provides zap logger.
func Plugin(options ...zap.Option) sen.Plugin {
	return &zapPlugin{
		options: options,
	}
}

type zapPlugin struct {
	App *sen.Application `inject:"app"`

	options []zap.Option
}

// Init initialises zap logger for the application.
// The logger will be regisreted under "logger" tag.
func (p *zapPlugin) Init() error {
	logger, err := zap.NewProduction(p.options...)
	if err != nil {
		return err
	}

	p.App.OnPostShutdown(func(_ context.Context) error {
		return logger.Sync()
	})

	return p.App.Register("logger", logger)
}
