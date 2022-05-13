package sen_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bongnv/sen"
)

type mockFactory struct {
	Data int `inject:"data"`
	err  error
}

func (f mockFactory) New() (interface{}, error) {
	return f.Data, f.err
}

func TestFactory(t *testing.T) {
	t.Run("should return an error when couldn't initialise a new component", func(t *testing.T) {
		p := sen.Factory[interface{}]("mock-component", &mockFactory{
			err: errors.New("random error"),
		})

		app := sen.New()
		err := app.Apply(
			sen.Component("data", 10),
			p,
		)
		require.EqualError(t, err, "random error")
	})

	t.Run("should register successfully when there is no error", func(t *testing.T) {
		p := sen.Factory[interface{}]("mock-component", &mockFactory{})

		app := sen.New()
		err := app.Apply(
			sen.Component("data", 10),
			p,
		)
		require.NoError(t, err)

		component, errGet := sen.GetComponent(app, "mock-component")
		require.NoError(t, errGet)
		require.Equal(t, 10, component)
	})

	t.Run("should propergate error when couldn't populate depdencies", func(t *testing.T) {
		p := sen.Factory[interface{}]("mock-component", &mockFactory{})

		app := sen.New()
		err := app.Apply(
			p,
		)
		require.EqualError(t, err, "injector: data is not registered")
	})
}
