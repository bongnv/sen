package sen_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bongnv/sen"
)

type mockPlugin struct {
	Data int `inject:"data"`
}

func (mockPlugin) Init() error {
	return nil
}

func TestComponent(t *testing.T) {
	t.Run("should register a component into the application", func(t *testing.T) {
		component := &mockComponent{}
		app := sen.New()
		err := app.Apply(
			sen.Component("data", 10),
			sen.Component("need-data", component),
		)
		require.NoError(t, err)
		require.Equal(t, component.Data, 10)
	})

	t.Run("should propergate error if a component cannot be registered", func(t *testing.T) {
		component := &mockComponent{}
		app := sen.New()
		err := app.Apply(sen.Component("need-data", component))
		require.EqualError(t, err, "injector: data is not registered")
	})
}

func TestModule(t *testing.T) {
	t.Run("should apply all plugins into the application", func(t *testing.T) {
		component := &mockComponent{}
		m := sen.Module(
			sen.Component("data", 10),
			sen.Component("need-data", component),
		)
		app := sen.New()
		require.NoError(t, app.Apply(m))
		require.Equal(t, component.Data, 10)
	})

	t.Run("should propagate error if one plugin returns error", func(t *testing.T) {
		m := sen.Module(
			sen.Component("error-plugin", &mockComponent{}),
			sen.Component("ok-plugin", 10),
		)
		app := sen.New()
		err := app.Apply(m)
		require.EqualError(t, err, "injector: data is not registered")
	})

	t.Run("should propagate error if one plugin doesn't have enough dependencies", func(t *testing.T) {
		m := sen.Module(
			&mockPlugin{},
		)

		app := sen.New()
		err := app.Apply(m)
		require.EqualError(t, err, "injector: data is not registered")
	})
}
