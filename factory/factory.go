package factory

import (
	"context"

	"github.com/bongnv/sen/app"
)

type Factory[T any] interface {
	New(ctx context.Context) (T, error)
}

type FactoryPlugin[T any] struct {
	app.ComponentPlugin
	Factory Factory[T]
}

func (f FactoryPlugin[_]) Apply(ctx context.Context) error {
	t, err := f.Factory.New(ctx)
	if err != nil {
		return err
	}

	f.Component = t
	return f.ComponentPlugin.Apply(ctx)
}

func Component[T any](name string, f Factory[T]) app.Plugin {
	p := &FactoryPlugin[T]{}
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
	FactoryPlugin[T]
}

func (f ServiceFactoryPlugin[T]) Apply(ctx context.Context) error {
	if err := f.FactoryPlugin.Apply(ctx); err != nil {
		return err
	}

	s := f.Component.(T)
	f.App.OnRun(s.Run)
	f.App.OnShutdown(s.Shutdown)
	return nil
}

func Service[T IService](name string, f Factory[T]) app.Plugin {
	p := &ServiceFactoryPlugin[T]{}
	p.Name = name
	p.Factory = f
	return p
}
