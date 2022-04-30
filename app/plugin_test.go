package app_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bongnv/sen/app"
)

func TestComponent(t *testing.T) {
	t.Run("should register a component into the application", func(t *testing.T) {
		component := &mockComponent{}
		app := app.New(
			app.Component("data", 10),
			app.Component("need-data", component),
		)

		require.NoError(t, app.Run(context.Background()))
		require.Equal(t, component.Data, 10)
	})

	t.Run("should propergate error if a component cannot be registered", func(t *testing.T) {
		component := &mockComponent{}
		app := app.New(
			app.Component("need-data", component),
		)

		require.EqualError(t, app.Run(context.Background()), "injector: data is not registered")
	})
}
