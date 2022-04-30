package sen

import (
	"context"

	"github.com/bongnv/sen/app"
	"github.com/bongnv/sen/factory"
	"github.com/bongnv/sen/shutdown"
)

type Foo struct{}

func (*Foo) Run(_ context.Context) error {
	return nil
}

func (*Foo) Shutdown(_ context.Context) error {
	return nil
}

type FooFactory struct{}

func (f *FooFactory) New(ctx context.Context) (*Foo, error) {
	return &Foo{}, nil
}

func Run() error {
	ctx := context.Background()
	app := app.New(
		shutdown.New(),
		app.Component("some-foo", &Foo{}),
		factory.Component[*Foo]("component-only", &FooFactory{}),
		factory.Service[*Foo]("some", &FooFactory{}),
	)

	if err := app.Run(ctx); err != nil {
		return err
	}

	return nil
}
