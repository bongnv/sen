package factory

import (
	"context"

	"github.com/bongnv/sen/app"
)

type Factory[T any] interface {
	New(ctx context.Context) (T, error)
}

type Plugin[T any] struct {
	App     *app.Application `inject:"app"`
	Factory Factory[T]
	Name    string
}

func (f Plugin[_]) Apply(ctx context.Context) error {
	t, err := f.Factory.New(ctx)
	if err != nil {
		return err
	}

	return f.App.ApplyPlugin(ctx, app.Component(f.Name, t))
}

func Component[T any](name string, f Factory[T]) app.Plugin {
	p := &Plugin[T]{}
	p.Name = name
	p.Factory = f
	return p
}

type IService interface {
	// Run is called to start the service. The function shouldn't be returned if the service is still running.
	Run(ctx context.Context) error
	// Shutdown is called to graceful shut down the service.
	Shutdown(ctx context.Context) error
}

type ServiceFactoryPlugin[T IService] struct {
	App     *app.Application `inject:"app"`
	Factory Factory[T]
	Name    string
}

func (f ServiceFactoryPlugin[T]) Apply(ctx context.Context) error {
	s, err := f.Factory.New(ctx)
	if err != nil {
		return err
	}

	f.App.OnRun(s.Run)
	f.App.OnShutdown(s.Shutdown)
	return f.App.ApplyPlugin(ctx, app.Component(f.Name, s))
}

func Service[T IService](name string, f Factory[T]) app.Plugin {
	return &ServiceFactoryPlugin[T]{
		Name:    name,
		Factory: f,
	}
}
