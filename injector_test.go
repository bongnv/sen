package sen_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bongnv/sen"
)

type mockComponent struct {
	Data int `inject:"data"`
}

func TestInjector(t *testing.T) {
	t.Run("should return error if there is no registered dependency", func(t *testing.T) {
		injector := sen.NewInjector()
		err := injector.Inject(&mockComponent{})
		require.EqualError(t, err, "injector: data is not registered")
	})
}

func TestInjector_Retrieve(t *testing.T) {
	t.Run("should return the component if it's registered", func(t *testing.T) {
		injector := sen.NewInjector()

		err := injector.Register("data", 10)
		require.NoError(t, err)

		data, err := injector.Retrieve("data")
		require.NoError(t, err)
		require.Equal(t, 10, data)
	})

	t.Run("should return an error if the component isn't registered", func(t *testing.T) {
		injector := sen.NewInjector()
		data, err := injector.Retrieve("data")
		require.Equal(t, err, sen.ErrComponentNotRegistered)
		require.Nil(t, data)
	})
}
